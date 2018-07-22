package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// Picture object
type Picture struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Upvotes         uint32    `json:"upvotes" db:"upvotes"`
	Downvotes       uint32    `json:"downvotes" db:"downvotes"`
	Sorting         float32   `json:"sorting" db:"sorting"`
	ConfidenceLevel float32   `json:"confidenceLevel" db:"confidence_level"`
	Filename        string    `json:"filename" db:"filename"`
	Hash            string    `json:"hash" db:"hash"`
	Import          Import    `json:"-" belongs_to:"import"`
	ImportID        uuid.UUID `json:"importId" db:"import_id"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

func (p Picture) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Pictures collection
type Pictures []Picture

func (p Pictures) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *Picture) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (p *Picture) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (p *Picture) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
