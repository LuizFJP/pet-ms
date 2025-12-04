package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPetValidate_EmptyAction_DoesNotRunDefaultValidation(t *testing.T) {
	pet := &Pet{
		Uuid:         uuid.New(),
		UuidGuardian: uuid.New(),
		// Campos inválidos de propósito
		Name:      "",
		Breed:     "",
		BirthYear: time.Now().Year() + 1,
		Specie:    Dog,
	}

	// action == "" cai no case "" e NÃO chama validateDefault
	errs := pet.Validate("")

	if len(errs) != 0 {
		t.Fatalf("expected no errors when action is empty, got: %v", errs)
	}
}

func TestPetValidate_WithAction_PopulatesErrorsOnInvalidPet(t *testing.T) {
	pet := &Pet{
		Uuid:         uuid.New(),
		UuidGuardian: uuid.New(),
		// todos inválidos pra cobrir todos os ifs:
		Name:      "",
		Breed:     "",
		BirthYear: time.Now().Year() + 1, // futuro
		Specie:    Cat,
	}

	// action != "" entra no default do switch e chama validateDefault
	errs := pet.Validate("create")

	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d: %v", len(errs), errs)
	}

	if errs["pet name is required"] != "pet name is empty" {
		t.Errorf("expected name error, got %q", errs["pet name is required"])
	}
	if errs["pet breed is required"] != "pet breed is empty" {
		t.Errorf("expected breed error, got %q", errs["pet breed is required"])
	}
	if errs["birth_year"] != "year out of range" {
		t.Errorf("expected birth_year error, got %q", errs["birth_year"])
	}
}

func TestPetValidate_WithValidPetAndAction_NoErrors(t *testing.T) {
	pet := &Pet{
		NIdentification: 1,
		Uuid:            uuid.New(),
		UuidGuardian:    uuid.New(),
		Name:            "Rex",
		BirthYear:       time.Now().Year(), // ano válido
		Breed:           "SRD",
		Specie:          Dog,
	}

	errs := pet.Validate("update")

	if len(errs) != 0 {
		t.Fatalf("expected no errors for valid pet, got: %v", errs)
	}
}
