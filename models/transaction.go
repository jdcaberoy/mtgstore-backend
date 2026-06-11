package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Transaction is used by pop to map your transactions database table to your go code.
type Transaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	SenderID   uuid.UUID `json:"sender_id" db:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id" db:"receiver_id"`

	Sender   User `json:"sender" belongs_to:"user" fk_id:"sender_id"`
	Receiver User `json:"receiver" belongs_to:"user" fk_id:"receiver_id"`

	RequestorName     string            `json:"requestor_name" db:"requestor_name"`
	SenderName        string            `json:"sender_name" db:"sender_name"`
	ReceiverName      string            `json:"receiver_name" db:"receiver_name"`
	TransactionStatus TransactionStatus `json:"transaction_status" db:"transaction_status"`
}

type TransactionStatus string

const (
	Pending   TransactionStatus = "Pending"   // request sent by initiator
	Approved                    = "Approved"  // request approved by other party
	Declined                    = "Declined"  // request declined by other party
	Cancelled                   = "Cancelled" // request cancelled by initiator
)

// String is not required by pop and may be deleted
func (t Transaction) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Transactions is not required by pop and may be deleted
type Transactions []Transaction

// String is not required by pop and may be deleted
func (t Transactions) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Transaction) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Transaction) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Transaction) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
