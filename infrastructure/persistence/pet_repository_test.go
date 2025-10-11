package persistence

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LuizFJP/pet-ms/domain/entity"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open("sqlite3", ":memory:")
	require.NoError(t, err, "failed to open sqlite in-memory for tests")

	db.LogMode(false)

	require.NoError(t, db.AutoMigrate(&entity.Pet{}).Error, "failed to automigrate Pet")
	return db
}

func TestPetRepository_SavePet_Success(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()
	repo := NewPetRepository(db)

	p := &entity.Pet{
		Uuid:            uuid.New(),
		NIdentification: 1,
		UuidGuardian:    uuid.New(),
		Name:            "Mingau",
		BirthYear:       2020,
		Breed:           "SRD",
		Specie:          1,
	}

	saved, errMap := repo.SavePet(p)
	require.Nil(t, errMap, "expected no db_error on save")
	require.NotNil(t, saved, "expected saved pet not to be nil")

	var got entity.Pet
	require.NoError(t, db.Where("uuid = ?", p.Uuid).First(&got).Error)
	assert.Equal(t, p.Name, got.Name)
	assert.Equal(t, p.UuidGuardian, got.UuidGuardian)
}

func TestPetRepository_GetPet_FoundAndNotFound(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	repo := NewPetRepository(db)

	known := &entity.Pet{
		Uuid:            uuid.New(),
		NIdentification: 2,
		UuidGuardian:    uuid.New(),
		Name:            "Thor",
		BirthYear:       2018,
		Breed:           "Labrador",
		Specie:          2,
	}

	require.NoError(t, db.Create(known).Error)

	got, errMap := repo.GetPet(known.Uuid.String())
	require.Nil(t, errMap)
	require.NotNil(t, got)
	assert.Equal(t, known.Name, got.Name)

	_, notFound := repo.GetPet(uuid.New().String())
	require.NotNil(t, notFound)
	assert.Contains(t, notFound, "db_error", "gorm returns an error for not found on First()")

}

func TestPetRepository_UpdatePet_Success(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()
	repo := NewPetRepository(db)

	original := &entity.Pet{
		Uuid:            uuid.New(),
		NIdentification: 2,
		UuidGuardian:    uuid.New(),
		Name:            "Luna",
		BirthYear:       2019,
		Breed:           "Beagle",
		Specie:          2,
	}

	require.NoError(t, db.Create(original).Error)

	original.Name = "Luna Updated"
	original.BirthYear = 2021
	original.Breed = "Beagle Tricolor"

	updated, errMap := repo.UpdatePet(original)
	require.Nil(t, errMap, "unexpected error map on update")
	require.NotNil(t, updated)

	assert.Equal(t, "Luna Updated", updated.Name)
	assert.Equal(t, 2021, updated.BirthYear)
	assert.Equal(t, "Beagle Tricolor", updated.Breed)
}

func TestPetRepository_UpdatePet_NotFound(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	repo := NewPetRepository(db)

	ghost := &entity.Pet{
		Uuid:            uuid.New(),
		NIdentification: 2,
		UuidGuardian:    uuid.New(),
		Name:            "Ghost",
		BirthYear:       2015,
		Breed:           "Unknown",
		Specie:          2,
	}

	updated, errMap := repo.UpdatePet(ghost)
	require.Nil(t, updated)
	require.NotNil(t, errMap)
	assert.Contains(t, errMap, "not_found")
	assert.Equal(t, "pet not found", errMap["not_found"])
}

func TestPetRepository_DeletePet_Success(t *testing.T) {
	db := newTestDB(t)
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing db")
		}
	}(db)

	repo := NewPetRepository(db)

	guardian := uuid.New()
	otherGuardian := uuid.New()

	pets := []entity.Pet{
		{
			Uuid:            uuid.New(),
			NIdentification: 1,
			UuidGuardian:    guardian,
			Name:            "Pingo",
			BirthYear:       2017,
			Breed:           "SRD",
			Specie:          2,
		},
		{
			Uuid:            uuid.New(),
			NIdentification: 2,
			UuidGuardian:    guardian,
			Name:            "Nina",
			BirthYear:       2022,
			Breed:           "SRD",
			Specie:          1,
		},
		{
			Uuid:            uuid.New(),
			NIdentification: 3,
			UuidGuardian:    otherGuardian,
			Name:            "Bidu",
			BirthYear:       2023,
			Breed:           "SRD",
			Specie:          2,
		},
	}

	for _, p := range pets {
		require.NoError(t, db.Create(&p).Error)
	}
	msg, errMap := repo.DeletePet(guardian.String())
	require.Nil(t, errMap)
	require.NotNil(t, msg)
	assert.Equal(t, "2 pet(s) deletados!", msg["message"])

	var count int
	require.NoError(t, db.Model(&entity.Pet{}).Count(&count).Error)
	assert.Equal(t, 1, count)
}

func TestPetRepository_DeletePet_NotFound(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	repo := NewPetRepository(db)

	_, errMap := repo.DeletePet(uuid.New().String())
	require.NotNil(t, errMap)
	assert.Contains(t, errMap, "not_found")
	assert.Equal(t, "nenhum pet encontrado para esse guardião", errMap["not_found"])
}

func TestPetRepository_SavePet_DBError(t *testing.T) {
	db := newTestDB(t)

	repo := NewPetRepository(db)

	require.NoError(t, db.Close())

	p := &entity.Pet{
		Uuid:            uuid.New(),
		NIdentification: 3,
		UuidGuardian:    uuid.New(),
		Name:            "Fail",
		BirthYear:       2000,
		Breed:           "None",
		Specie:          3,
	}

	saved, errMap := repo.SavePet(p)
	require.Nil(t, saved)
	require.NotNil(t, errMap)
	assert.Contains(t, errMap, "db_error")
	assert.Contains(t, fmt.Sprintf("%v", errMap["db_error"]), "closed", "expect error mentions closed DB")
}
