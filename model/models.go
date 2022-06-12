package model

import (
	"image"
	"time"
)

// This file contains all structs representing all models for the HealthCheck-API

var CredentialParameters = []string{"username", "password"}

type UserCredentials struct {
	Username string `firestore:"username" json:"username" mapstructure:"username"`
	Password string `firestore:"password" json:"password" `
	User     *User  `firestore:"user,omitempty" json:"user,omitempty" `
}

var UserParameters = []string{"_id", "name", "telephone", "admin"}

//only needed for user registration (which has low priority in implementation)
type User struct {
	Id        string `firestore:"_id" json:"_id" mapstructure:"_id"`
	Name      string `firestore:"name" json:"name" `
	Telephone string `firestore:"telephone" json:"telephone" `
	Admin     bool   `firestore:"admin" json:"admin"`
}

var ProjectParameters = []string{"_id", "name"}

type Project struct {
	Id   string `firestore:"_id" json:"_id" mapstructure:"_id"`
	Name string `firestore:"name" json:"name" mapstructure:"name"`
}

var LocationParameters = []string{"_id", "name", "address", "projectId"}

type Location struct {
	Id        string          `firestore:"_id" json:"_id" mapstructure:"_id"`
	Name      string          `firestore:"name" json:"name" mapstructure:"name"`
	Address   LocationAddress `firestore:"address" json:"address" mapstructure:"address"`
	ProjectId string          `firestore:"projectId,omitempty" json:"projectId,omitempty" mapstructure:"projectId,omitempty"`
}

var AddressParameters = []string{"street", "zipcode", "city", "country"}

type LocationAddress struct {
	Street  string `firestore:"street" json:"street" mapstructure:"street"`
	Zipcode string `firestore:"zipcode" json:"zipcode" mapstructure:"zipcode"`
	City    string `firestore:"city" json:"city" mapstructure:"city"`
	Country string `firestore:"country" json:"country" mapstructure:"country"`
}

var DeviceParameters = []string{"_id", "locationId", "name"}

type Device struct {
	Id         string `firestore:"_id" json:"_id" mapstructure:"_id"`
	LocationId string `firestore:"locationId" json:"locationId" mapstructure:"locationId"`
	Name       string `firestore:"name" json:"name" mapstructure:"name"`
	Serial     *int   `firestore:"serial,omitempty" json:"serial,omitempty" mapstructure:"serial,omitempty"`
	LastPing   *Ping  `firestore:"lastPing,omitempty" json:"lastPing,omitempty" mapstructure:"lastPing,omitempty"`
}

type LocalDevice struct {
	LastPing *Ping
	Image    *image.Image
}

var PingParameters = []string{"id", "status", "timestamp", "version"}

type Ping struct {
	Id        string    `firestore:"-" json:"id" mapstructure:"id"`
	Status    string    `firestore:"status" json:"status" mapstructure:"status"`
	Timestamp time.Time `firestore:"timestamp" json:"timestamp" mapstructure:"-"` //change to Date
	Version   string    `firestore:"version" json:"version" mapstructure:"version"`
}

var DeviceStatus = [...]string{"ONLINE", "WARNING", "OFFLINE"}

var InventoryParameters = []string{"_id", "name"}

type Inventory struct {
	Id          string    `firestore:"_id" json:"_id" mapstructure:"_id"`
	Name        string    `firestore:"name" json:"name" mapstructure:"name"`
	DeviceId    string    `firestore:"deviceId" json:"locationId" mapstructure:"deviceId"`
	Notes       string    `firestore:"notes,omitempty" json:"notes,omitempty" mapstructure:"notes,omitempty"`
	BuyDate     time.Time `firestore:"buyDate,omitempty" json:"buyDate" mapstructure:"-"`
	MontageDate time.Time `firestore:"montageDate,omitempty" json:"montageDate" mapstructure:"-"`
}
