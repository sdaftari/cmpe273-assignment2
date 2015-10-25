package Controllers

import (  
    "encoding/json"
    "fmt"
    "net/http"
    "io"
    "io/ioutil"
    "strings"
    "errors"
    
    model "Assignment2/models"

    "github.com/julienschmidt/httprouter"
    "github.com/jmoiron/jsonq"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

//Static part of Google API for address
const (
    LocationUrl = "http://maps.google.com/maps/api/geocode/json?address="
)

type (  
    // UserController represents the controller for operating on the User resource
    UserController struct{
        session *mgo.Session
    }
)

func NewUserController(s *mgo.Session) *UserController {  
    //returnd he object of UserController
    return &UserController{s}
}

func (uc UserController) DeleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    //Get the locationId
    id := p.ByName("locationId")

    if !bson.IsObjectIdHex(id) {
        rw.WriteHeader(404)
        return
    }

    objectId := bson.ObjectIdHex(id)

    // Remove user
    if err := uc.session.DB("usersdb").C("userAddresses").RemoveId(objectId); err != nil {
        rw.WriteHeader(404)
        return
    }

    // Write status
    rw.WriteHeader(200)
}

func (uc UserController) UpdateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    userUpdateRequest := model.UserUpdateRequest{}
    userResponse := model.UserResponse{}

    //Get the locationId
    id := p.ByName("locationId")

    if !bson.IsObjectIdHex(id) {
        rw.WriteHeader(404)
        return
    }

    objectId := bson.ObjectIdHex(id)

    //Get the information to be updated
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
    }

    err = json.Unmarshal(body, &userUpdateRequest)
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
    } 

    err1 := uc.session.DB("usersdb").C("userAddresses").Update(bson.M{"_id":objectId }, 
        bson.M{"$set": bson.M{"address": userUpdateRequest.Address, "city": userUpdateRequest.City, "state": userUpdateRequest.State, "zip": userUpdateRequest.Zip}})

    if(err1 != nil) {
        fmt.Println("Fatal error ", err.Error())
    }

    //Get the updated object
    if err := uc.session.DB("usersdb").C("userAddresses").FindId(objectId).One(&userResponse); err != nil {
        rw.WriteHeader(404)
        return
    }

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(userResponse)

    // Write content-type, statuscode, payload
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", uj)
}


func (uc UserController) GetLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    userResponse := model.UserResponse{}

    id := p.ByName("locationId")

    if !bson.IsObjectIdHex(id) {
        rw.WriteHeader(404)
        return
    }

    objectId := bson.ObjectIdHex(id)

    // Fetch user
    if err := uc.session.DB("usersdb").C("userAddresses").FindId(objectId).One(&userResponse); err != nil {
        rw.WriteHeader(404)
        return
    }

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(userResponse)

    // Write content-type, statuscode, payload
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", uj)
}

func (uc UserController) CreateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    user := model.UserRequest{}
    userResponse := model.UserResponse{}

    //Get data from request
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
    }

    err = json.Unmarshal(body, &user)
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
    }    

    //Get Google API data in json
    response, err := GetLocation(user)

    //Get lotitude and longitude
    resp := make(map[string]interface{})

    locationBody, err := ioutil.ReadAll(response)
    err = json.Unmarshal(locationBody, &resp)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    //Create a json query to retrieve latitude and longitude
    jq := jsonq.NewQuery(resp)
    status, err := jq.String("status")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if status != "OK" {
      err = errors.New(status)
      return
    }

    lat, err := jq.Float("results", "0", "geometry", "location", "lat")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    lng, err := jq.Float("results", "0", "geometry", "location", "lng")
    if err != nil {
      return
    }

    //Createing Response object
    userResponse.Id = bson.NewObjectId()
    userResponse.Name = user.Name
    userResponse.Address = user.Address
    userResponse.City = user.City
    userResponse.State = user.State
    userResponse.Zip = user.Zip
    userResponse.Coordinates.Lat = lat
    userResponse.Coordinates.Lng = lng

    //Insert into mongodb
    uc.session.DB("usersdb").C("userAddresses").Insert(userResponse)

    uj, _ := json.Marshal(userResponse)

    // Write content-type, statuscode, payload
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201)
    fmt.Fprintf(rw, "%s", uj)
}

func GetLocation(user model.UserRequest) (io.ReadCloser, error) {
    var userAddress string

    //Get Address
    address := strings.Split(user.Address, " ")
    if len(address) > 1 {
        for i := 0; i < len(address); i++ {
            userAddress = userAddress + address[i] + "+"
        }
    } else {
        userAddress = userAddress + user.Address
    }
    
    //Get City
    city := strings.Split(user.City, " ")
    if len(city) > 1 {
        for i := 0; i < len(city); i++ {
            userAddress = userAddress + city[i] + "+"
        }
    } else {
        userAddress = userAddress + user.City
    }

    //Get State
    state := strings.Split(user.State, " ")
    if len(state) > 1 {
        for i := 0; i < len(state); i++ {
            userAddress = userAddress + state[i] + "+"
        }
    } else {
        userAddress = userAddress + user.State
    }

    userAddress = userAddress + user.Zip

    //Create address url for google api
    userAddress = LocationUrl + userAddress + "&sensor=false"

    //Get response from google api
    response, err := http.Get(userAddress)
    if err != nil {
        return nil, err
    }

    return response.Body, nil;
}