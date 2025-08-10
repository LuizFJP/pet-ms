package repository

import "github.com/LuizFJP/pet-ms/domain/entity"

type PetRepository interface {
	SavePet(pet *entity.Pet) (*entity.Pet, map[string]string)
	GetPet(uuid string) (*entity.Pet, map[string]string)
	UpdatePet(pet *entity.Pet) (*entity.Pet, map[string]string)
	DeletePet(uuid string) (map[string]string, map[string]string)
}
