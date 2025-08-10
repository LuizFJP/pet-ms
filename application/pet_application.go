package application

import (
	"github.com/LuizFJP/pet-ms/domain/entity"
	"github.com/LuizFJP/pet-ms/domain/repository"
)

type petApplication struct {
	pr repository.PetRepository
}

var _ PetApplicationInterface = &petApplication{}

type PetApplicationInterface interface {
	SavePet(pet *entity.Pet) (*entity.Pet, map[string]string)
	GetPet(uuid string) (*entity.Pet, map[string]string)
	UpdatePet(pet *entity.Pet) (*entity.Pet, map[string]string)
	DeletePet(uuid string) (map[string]string, map[string]string)
}

func (p *petApplication) SavePet(pet *entity.Pet) (*entity.Pet, map[string]string) {
	return nil, nil
}

func (p *petApplication) GetPet(uuid string) (*entity.Pet, map[string]string) {
	return nil, nil
}

func (p *petApplication) UpdatePet(pet *entity.Pet) (*entity.Pet, map[string]string) {
	return nil, nil
}

func (p *petApplication) DeletePet(uuid string) (map[string]string, map[string]string) {
	return nil, nil
}
