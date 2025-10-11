package persistence

import (
	"fmt"
	"github.com/LuizFJP/pet-ms/domain/entity"
	"github.com/LuizFJP/pet-ms/domain/repository"
	"github.com/jinzhu/gorm"
)

type PetRepo struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) *PetRepo {
	return &PetRepo{db}
}

var _ repository.PetRepository = &PetRepo{}

func (p *PetRepo) SavePet(pet *entity.Pet) (*entity.Pet, map[string]string) {
	dbErr := map[string]string{}
	err := p.db.Debug().Create(pet).Error
	if err != nil {
		dbErr["db_error"] = err.Error()
		return nil, dbErr
	}
	return pet, nil
}

func (p *PetRepo) GetPet(uuid string) (*entity.Pet, map[string]string) {
	pet := &entity.Pet{}
	dbErr := map[string]string{}
	err := p.db.Debug().Where("uuid = ?", uuid).First(pet).Error
	if err != nil {
		dbErr["db_error"] = err.Error()
		return nil, dbErr
	}
	return pet, nil
}
func (p *PetRepo) UpdatePet(pet *entity.Pet) (*entity.Pet, map[string]string) {
	dbErr := map[string]string{}

	tx := p.db.Debug().
		Model(&entity.Pet{}).
		Where("uuid = ?", pet.Uuid).
		Updates(map[string]interface{}{
			"n_identification": pet.NIdentification,
			"uuid_guardian":    pet.UuidGuardian,
			"name":             pet.Name,
			"birth_year":       pet.BirthYear,
			"breed":            pet.Breed,
			"specie":           pet.Specie,
		})

	if tx.Error != nil {
		dbErr["db_error"] = tx.Error.Error()
		return nil, dbErr
	}
	if tx.RowsAffected == 0 {
		dbErr["not_found"] = "pet not found"
		return nil, dbErr
	}

	updated := &entity.Pet{}
	if err := p.db.Where("uuid = ?", pet.Uuid).First(updated).Error; err != nil {
		dbErr["db_error"] = err.Error()
		return nil, dbErr
	}
	return updated, nil
}
func (p *PetRepo) DeletePet(uuidGuardian string) (map[string]string, map[string]string) {
	dbErr := map[string]string{}
	tx := p.db.Debug().
		Where("uuid_guardian = ?", uuidGuardian).
		Delete(&entity.Pet{})

	if tx.Error != nil {
		dbErr["db_error"] = tx.Error.Error()
		return nil, dbErr
	}
	if tx.RowsAffected == 0 {
		return nil, map[string]string{"not_found": "nenhum pet encontrado para esse guardião"}
	}

	return map[string]string{
		"message": fmt.Sprintf("%d pet(s) deletados!", tx.RowsAffected),
	}, nil
}
