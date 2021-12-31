package models

import "time"

type User struct {
	ID           int       `json:"id" bson:"id"`
	Name         string    `json:"name" bson:"name"`
	DOB          time.Time `json:"dob" bson:"dob"`
	Address      Address   `json:"address" bson:"address"`
	IsVerified   bool      `json:"isVerified" bson:"isVerified"`
	Remarks      string    `json:"remarks" bson:"remarks"`
	Subscription string    `json:"subscription" bson:"subscription"`
}

type Address struct {
	Block  string `json:"block" bson:"block"`
	Street string `json:"street" bson:"street"`
	City   string `json:"city" bson:"city"`
}
