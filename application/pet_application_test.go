package application

import (
	"github.com/LuizFJP/pet-ms/domain/entity"
	"github.com/LuizFJP/pet-ms/domain/repository"
	"reflect"
	"testing"
)

// Compile-time check: our mock satisfies the repository interface
var _ repository.PetRepository = (*mockPetRepository)(nil)

// mockPetRepository is a lightweight stub with pluggable behavior
// and call tracking for assertions
type mockPetRepository struct {
	saveFunc   func(p *entity.Pet) (*entity.Pet, map[string]string)
	getFunc    func(id string) (*entity.Pet, map[string]string)
	updateFunc func(p *entity.Pet) (*entity.Pet, map[string]string)
	deleteFunc func(id string) (map[string]string, map[string]string)

	saveCalledWith   *entity.Pet
	getCalledWith    string
	updateCalledWith *entity.Pet
	deleteCalledWith string
}

func (m *mockPetRepository) SavePet(p *entity.Pet) (*entity.Pet, map[string]string) {
	m.saveCalledWith = p
	if m.saveFunc != nil {
		return m.saveFunc(p)
	}
	return p, nil
}

func (m *mockPetRepository) GetPet(id string) (*entity.Pet, map[string]string) {
	m.getCalledWith = id
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return &entity.Pet{}, nil
}

func (m *mockPetRepository) UpdatePet(p *entity.Pet) (*entity.Pet, map[string]string) {
	m.updateCalledWith = p
	if m.updateFunc != nil {
		return m.updateFunc(p)
	}
	return p, nil
}

func (m *mockPetRepository) DeletePet(id string) (map[string]string, map[string]string) {
	m.deleteCalledWith = id
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return map[string]string{"status": "deleted"}, nil
}

func TestNewPetApplication_ReturnsConcreteAndWrapsRepo(t *testing.T) {
	mock := &mockPetRepository{}
	app := NewPetApplication(mock)
	if app == nil {
		t.Fatalf("NewPetApplication returned nil")
	}

	// Verify it is the concrete implementation we expect (same package allows this)
	pa, ok := app.(*petApplication)
	if !ok {
		t.Fatalf("expected *petApplication, got %T", app)
	}
	if pa.pr != mock {
		t.Fatalf("petApplication.pr should be the provided repository instance")
	}
}

func TestSavePet_DelegatesToRepository(t *testing.T) {
	in := &entity.Pet{ /* fill fields if you have them */ }
	wantPet := &entity.Pet{}
	wantErrs := map[string]string{"ok": "true"}

	mock := &mockPetRepository{
		saveFunc: func(p *entity.Pet) (*entity.Pet, map[string]string) {
			if p != in {
				t.Fatalf("repo.SavePet received wrong pointer")
			}
			return wantPet, wantErrs
		},
	}

	app := NewPetApplication(mock)
	gotPet, gotErrs := app.SavePet(in)

	if mock.saveCalledWith != in {
		t.Fatalf("SavePet should forward the same pointer to repo")
	}
	if gotPet != wantPet {
		t.Fatalf("SavePet should return repo's pet. got=%p want=%p", gotPet, wantPet)
	}
	if !reflect.DeepEqual(gotErrs, wantErrs) {
		t.Fatalf("SavePet should return repo's errors. got=%v want=%v", gotErrs, wantErrs)
	}
}

func TestGetPet_DelegatesToRepository(t *testing.T) {
	wantID := "123"
	wantPet := &entity.Pet{}
	wantErrs := map[string]string{"ok": "true"}

	mock := &mockPetRepository{
		getFunc: func(id string) (*entity.Pet, map[string]string) {
			if id != wantID {
				t.Fatalf("repo.GetPet received wrong id: got=%s want=%s", id, wantID)
			}
			return wantPet, wantErrs
		},
	}

	app := NewPetApplication(mock)
	gotPet, gotErrs := app.GetPet(wantID)

	if mock.getCalledWith != wantID {
		t.Fatalf("GetPet should pass the id to repo. got=%s want=%s", mock.getCalledWith, wantID)
	}
	if gotPet != wantPet {
		t.Fatalf("GetPet should return repo's pet. got=%p want=%p", gotPet, wantPet)
	}
	if !reflect.DeepEqual(gotErrs, wantErrs) {
		t.Fatalf("GetPet should return repo's errors. got=%v want=%v", gotErrs, wantErrs)
	}
}

func TestUpdatePet_DelegatesToRepository(t *testing.T) {
	in := &entity.Pet{}
	wantPet := &entity.Pet{}
	wantErrs := map[string]string{"updated": "true"}

	mock := &mockPetRepository{
		updateFunc: func(p *entity.Pet) (*entity.Pet, map[string]string) {
			if p != in {
				t.Fatalf("repo.UpdatePet received wrong pointer")
			}
			return wantPet, wantErrs
		},
	}

	app := NewPetApplication(mock)
	gotPet, gotErrs := app.UpdatePet(in)

	if mock.updateCalledWith != in {
		t.Fatalf("UpdatePet should forward the same pointer to repo")
	}
	if gotPet != wantPet {
		t.Fatalf("UpdatePet should return repo's pet. got=%p want=%p", gotPet, wantPet)
	}
	if !reflect.DeepEqual(gotErrs, wantErrs) {
		t.Fatalf("UpdatePet should return repo's errors. got=%v want=%v", gotErrs, wantErrs)
	}
}

func TestDeletePet_DelegatesToRepository(t *testing.T) {
	wantID := "abc-uuid"
	wantResp := map[string]string{"status": "ok"}
	wantErrs := map[string]string{"error": ""}

	mock := &mockPetRepository{
		deleteFunc: func(id string) (map[string]string, map[string]string) {
			if id != wantID {
				t.Fatalf("repo.DeletePet received wrong id: got=%s want=%s", id, wantID)
			}
			return wantResp, wantErrs
		},
	}

	app := NewPetApplication(mock)
	gotResp, gotErrs := app.DeletePet(wantID)

	if mock.deleteCalledWith != wantID {
		t.Fatalf("DeletePet should pass the id to repo. got=%s want=%s", mock.deleteCalledWith, wantID)
	}
	if !reflect.DeepEqual(gotResp, wantResp) {
		t.Fatalf("DeletePet should return repo's response. got=%v want=%v", gotResp, wantResp)
	}
	if !reflect.DeepEqual(gotErrs, wantErrs) {
		t.Fatalf("DeletePet should return repo's errors. got=%v want=%v", gotErrs, wantErrs)
	}
}
