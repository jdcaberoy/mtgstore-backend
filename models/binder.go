package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Binder is used by pop to map your binders database table to your go code.
type Binder struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	OwnerID       uuid.UUID `db:"owner_id" json:"owner_id"`
	OwnerUsername string    `db:"owner_username" json:"owner_username"`
	Hidden        bool      `db:"hidden" json:"hidden"`

	Owner *User `belongs_to:"user" json:"user,omitempty" fk_id:"owner_id" `
	Cards Cards `has_many:"cards" json:"cards,omitempty"`
}

// String is not required by pop and may be deleted
func (b Binder) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Binders is not required by pop and may be deleted
type Binders []Binder

// String is not required by pop and may be deleted
func (b Binders) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *Binder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *Binder) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *Binder) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
