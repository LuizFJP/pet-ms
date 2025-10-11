package entity

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type PetType int

const (
	Dog PetType = iota
	Cat
)

type Pet struct {
	NIdentification uint      `gorm:"AUTO_INCREMENT"`
	Uuid            uuid.UUID `gorm:"primaryKey" json:"uuid"`
	UuidGuardian    uuid.UUID `json:"uuid_guardian"`
	Name            string    `json:"name"`
	BirthYear       int       `json:"birth_year"`
	Breed           string    `json:"breed"`
	Specie          PetType   `json:"specie"`
}

func (p *Pet) Validate(action string) map[string]string {
	errorMessages := make(map[string]string)

	switch strings.ToLower(action) {
	case "":
	default:
		p.validateDefault(errorMessages)
	}

	return errorMessages
}

func (p *Pet) validateDefault(errorMessages map[string]string) {
	if p.Name == "" {
		errorMessages["pet name is required"] = "pet name is empty"
	}

	if p.Breed == "" {
		errorMessages["pet breed is required"] = "pet breed is empty"
	}

	if p.BirthYear > time.Now().Year() {
		errorMessages["birth_year"] = "year out of range"
	}
}
