package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    //"encoding/json"
    "gopkg.in/mgo.v2"

    controllers "Assignment2/controllers"
 )

func main() {
	fmt.Println("Server is listening on 8080!")

	userController := controllers.NewUserController(getSession())

    mux := httprouter.New()

    //Create new location
    mux.POST("/locations", userController.CreateLocation)
    
    //Get the location
    mux.GET("/locations/:locationId", userController.GetLocation)

    //Update the location
    mux.PUT("/locations/:locationId", userController.UpdateLocation)

    //Delete the location
    mux.DELETE("/locations/:locationId", userController.DeleteLocation)
    
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }
    server.ListenAndServe()
}

func getSession() *mgo.Session {  
    // Connect to our local mongo
    s, err := mgo.Dial("mongodb://test:asdf1234#@ds041154.mongolab.com:41154/usersdb")

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    return s
}