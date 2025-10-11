package persistence

import (
	_ "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"

	"github.com/LuizFJP/pet-ms/domain/entity"
)

func TestRepositories_Close(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro criando sqlmock: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	mock.ExpectClose()

	gdb, err := gorm.Open("postgres", sqlDB)
	if err != nil {
		t.Fatalf("erro abrindo gorm com sqlmock: %v", err)
	}

	repos := &Repositories{db: gdb}

	if err := repos.Close(); err != nil {
		t.Fatalf("Close() retornou erro: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Close() não chamou o Close do sql.DB: %v", err)
	}
}

func TestRepositories_Automigrate_SQLite(t *testing.T) {
	gdb, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("erro abrindo sqlite em memória: %v", err)
	}
	gdb.LogMode(false)

	repos := &Repositories{db: gdb}

	_ = entity.Pet{}

	if err := repos.Automigrate(); err != nil {
		t.Fatalf("Automigrate() retornou erro: %v", err)
	}

	if !gdb.HasTable(&entity.Pet{}) {
		t.Fatalf("esperava que a tabela de Pet existisse após Automigrate")
	}
}
