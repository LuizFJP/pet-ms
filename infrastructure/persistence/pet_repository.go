package persistence

import (
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
	err := p.db.Debug().Create(&pet).Error
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
	err := p.db.Debug().Save(&pet).Error
	if err != nil {
		dbErr["db_error"] = err.Error()
		return nil, dbErr
	}
	return pet, nil
}
func (p *PetRepo) DeletePet(uuid string) (map[string]string, map[string]string) {
	pet := &entity.Pet{}
	dbErr := map[string]string{}
	err := p.db.Debug().Where("uuid = ?", uuid).Delete(pet).Error
	if err != nil {
		dbErr["db_error"] = err.Error()
		return nil, dbErr
	}
	return map[string]string{"message": "pet deletado!"}, nil
}
