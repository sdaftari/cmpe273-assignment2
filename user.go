package models

import "gopkg.in/mgo.v2/bson"

type (  
    // User represents the structure of our resource
    UserRequest struct {
        Name string `json:"name" bson:"name"`
		Address string `json:"address" bson:"address"`
		City string `json:"city" bson:"city"`
		State string `json:"state" bson:"state"`
		Zip string `json:"zip" bson:"zip"`
    }
)

type (  
    // User represents the structure of our resource
    UserResponse struct {
    	Id bson.ObjectId `json:"id" bson:"_id"`
        Name string `json:"name" bson:"name"`
		Address string `json:"address" bson:"address"`
		City string `json:"city" bson:"city"`
		State string `json:"state" bson:"state"`
		Zip string `json:"zip" bson:"zip"`
		Coordinates AddressCoordinates
    }
)

type (  
    // User represents the structure of our resource
    AddressCoordinates struct {
        Lat float64 `json:"lat" bson:"lat"`
		Lng float64 `json:"lng" bson:"lng"`
    }
)

type (  
    // User represents the structure of our resource
    UserUpdateRequest struct {
		Address string `json:"address" bson:"address"`
		City string `json:"city" bson:"city"`
		State string `json:"state" bson:"state"`
		Zip string `json:"zip" bson:"zip"`
    }
)
