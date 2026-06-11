package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Card is used by pop to map your cards database table to your go code.
type Card struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name       string    `db:"name" json:"name"`
	SFID       string    `db:"sfid" json:"sfid"`
	Set        string    `db:"set" json:"set"`
	Quantity   int       `db:"quantity" json:"quantity"`
	BinderID   uuid.UUID `db:"binder_id" json:"-"`
	Binder     *Binder   `belongs_to:"binder" json:"binder,omitempty"`
	BinderName string    `db:"binder_name" json:"binder_name"`
	Foil       bool      `db:"foil" json:"foil"`
	Price      float64   `db:"price" json:"price"`
	Remarks    string    `db:"remarks" json:"remarks"`
	Hidden     bool      `db:"hidden" json:"hidden"`
}

// String is not required by pop and may be deleted
func (c Card) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Cards is not required by pop and may be deleted
type Cards []Card

// String is not required by pop and may be deleted
func (c Cards) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Card) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Card) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Card) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
