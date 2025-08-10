package persistence

import (
	"fmt"
	"github.com/LuizFJP/pet-ms/domain/entity"
	"github.com/LuizFJP/pet-ms/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
)

type Repositories struct {
	Pet repository.PetRepository
	db  *gorm.DB
}

func NewPetRepo(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		Pet: NewPetRepository(db),
		db:  db,
	}, nil
}

func (s *Repositories) Close() error {
	return s.db.Close()
}

func (s *Repositories) Automigrate() error {
	return s.db.AutoMigrate(&entity.Pet{}).Error
}
