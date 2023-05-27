package model

import (
	"image"
	"time"
)

// This file contains all structs representing all models for the HealthCheck-API

type Model interface {
	ID() string
	FID() string
}

var CredentialParameters = []string{"username", "password"}

type UserCredentials struct {
	Username string `firestore:"username" json:"username" mapstructure:"username"`
	Password string `firestore:"password" json:"password" `
	User     `json:",omitempty" mapstructure:",squash,omitempty"`
}

var UserParameters = []string{"username", "password", "name", "telephone"}

type User struct {
	Id        string `firestore:"_id" json:"_id" mapstructure:"_id,omitempty"`
	Name      string `firestore:"name" json:"name" `
	Telephone string `firestore:"telephone" json:"telephone" `
	Admin     bool   `firestore:"admin" json:"admin"`
}

var ProjectParameters = []string{"name"}

type Project struct {
	Id   string `firestore:"_id" json:"_id" mapstructure:"_id,omitempty"`
	Name string `firestore:"name" json:"name" mapstructure:"name"`
}

func (p *Project) ID() string {
	return p.Id
}

func (p *Project) FID() string {
	return ""
}

var LocationParameters = []string{"name", "address"}

type Location struct {
	Id        string          `firestore:"_id" json:"_id" mapstructure:"_id,omitempty"`
	Name      string          `firestore:"name" json:"name" mapstructure:"name"`
	Address   LocationAddress `firestore:"address" json:"address" mapstructure:"address"`
	ProjectId *string         `firestore:"projectId,omitempty" json:"projectId,omitempty" mapstructure:"projectId,omitempty"`
}

func (l *Location) ID() string {
	return l.Id
}

func (l *Location) FID() string {
	if l.ProjectId == nil {
		return ""
	}
	return *l.ProjectId
}

var AddressParameters = []string{"street", "zipcode", "city", "country"}

type LocationAddress struct {
	Street  string `firestore:"street" json:"street" mapstructure:"street"`
	Zipcode string `firestore:"zipcode" json:"zipcode" mapstructure:"zipcode"`
	City    string `firestore:"city" json:"city" mapstructure:"city"`
	Country string `firestore:"country" json:"country" mapstructure:"country"`
}

var DeviceParameters = []string{"locationId", "name"}

type Device struct {
	Id         string `firestore:"_id" json:"_id" mapstructure:"_id,omitempty"`
	LocationId string `firestore:"locationId" json:"locationId" mapstructure:"locationId"`
	Name       string `firestore:"name" json:"name" mapstructure:"name"`
	Serial     *int   `firestore:"serial,omitempty" json:"serial,omitempty" mapstructure:"serial,omitempty"`
	LastPing   *Ping  `firestore:"lastPing,omitempty" json:"lastPing,omitempty" mapstructure:"lastPing,omitempty"`
}

func (d *Device) ID() string {
	return d.Id
}

func (d *Device) FID() string {
	return d.LocationId
}

type LocalImage struct {
	Image *image.Image
}

var PingParameters = []string{"id", "status", "timestamp", "version"}

type Ping struct {
	Id        string    `firestore:"-" json:"id" mapstructure:"id"`
	Status    string    `firestore:"status" json:"status" mapstructure:"status"`
	Timestamp time.Time `firestore:"timestamp" json:"timestamp" mapstructure:"-"`
	Version   string    `firestore:"version" json:"version" mapstructure:"version"`
}

func (p *Ping) ID() string {
	return p.Id
}

func (p *Ping) FID() string {
	return ""
}

var DeviceStatus = [...]string{"ONLINE", "WARNING", "OFFLINE"}

var InventoryParameters = []string{"name"}

type Inventory struct {
	Id          string    `firestore:"_id" json:"_id" mapstructure:"_id,omitempty"`
	Name        string    `firestore:"name" json:"name" mapstructure:"name"`
	DeviceId    *string   `firestore:"deviceId,omitempty" json:"deviceId" mapstructure:"deviceId"`
	Notes       *string   `firestore:"notes,omitempty" json:"notes,omitempty" mapstructure:"notes,omitempty"`
	BuyDate     time.Time `firestore:"buyDate,omitempty" json:"buyDate" mapstructure:"-"`
	MontageDate time.Time `firestore:"montageDate,omitempty" json:"montageDate" mapstructure:"-"`
}

func (i *Inventory) ID() string {
	return i.Id
}

func (i *Inventory) FID() string {
	if i.DeviceId == nil {
		return ""
	}
	return *i.DeviceId
}
