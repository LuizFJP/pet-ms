package grpc

import (
	"context"
	"testing"

	"github.com/LuizFJP/pet-ms/domain/entity"
	pb "github.com/LuizFJP/pet-ms/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type appMock struct {
	savePetFn   func(*entity.Pet) (*entity.Pet, map[string]string)
	updatePetFn func(*entity.Pet) (*entity.Pet, map[string]string)
	getPetFn    func(string) (*entity.Pet, map[string]string)
	deletePetFn func(string) (map[string]string, map[string]string)
}

func (m *appMock) SavePet(p *entity.Pet) (*entity.Pet, map[string]string) {
	if m.savePetFn != nil {
		return m.savePetFn(p)
	}
	return nil, map[string]string{"message": "not implemented"}
}

func (m *appMock) UpdatePet(p *entity.Pet) (*entity.Pet, map[string]string) {
	if m.updatePetFn != nil {
		return m.updatePetFn(p)
	}
	return nil, map[string]string{"message": "not implemented"}
}

func (m *appMock) GetPet(id string) (*entity.Pet, map[string]string) {
	if m.getPetFn != nil {
		return m.getPetFn(id)
	}
	return nil, map[string]string{"message": "not implemented"}
}

func (m *appMock) DeletePet(uuidGuardian string) (map[string]string, map[string]string) {
	if m.deletePetFn != nil {
		return m.deletePetFn(uuidGuardian)
	}
	return nil, map[string]string{"message": "not implemented"}
}

func makePet() *entity.Pet {
	return &entity.Pet{
		NIdentification: 101,
		Uuid:            uuid.New(),
		UuidGuardian:    uuid.New(),
		Name:            "Mingau",
		BirthYear:       2020,
		Breed:           "SRD",
		Specie:          2, // assume enum/int underlying
	}
}

func TestPetServer_Create_Success(t *testing.T) {
	app := &appMock{
		savePetFn: func(p *entity.Pet) (*entity.Pet, map[string]string) {
			// echo back with DB-assigned identification
			ret := *p
			ret.NIdentification = 123
			return &ret, nil
		},
	}
	s := NewPetServer(app)

	req := &pb.CreatePetRequest{
		Name:         "Mingau",
		UuidGuardian: uuid.New().String(),
		BirthYear:    2020,
		Breed:        "SRD",
		Specie:       2, // untyped int constant fits any int field
	}

	resp, err := s.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int64(123), resp.NIdentification)
	assert.NotEmpty(t, resp.Uuid)
	assert.Equal(t, req.UuidGuardian, resp.UuidGuardian)
	assert.Equal(t, "Mingau", resp.Name)
	assert.Equal(t, uint64(2020), resp.BirthYear)
	assert.Equal(t, "SRD", resp.Breed)
	assert.Equal(t, "2", resp.Specie) // server formats int -> string
}

func TestPetServer_Create_Error(t *testing.T) {
	app := &appMock{
		savePetFn: func(p *entity.Pet) (*entity.Pet, map[string]string) {
			return nil, map[string]string{"message": "save failed"}
		},
	}
	s := NewPetServer(app)

	req := &pb.CreatePetRequest{
		Name:         "Any",
		UuidGuardian: uuid.New().String(),
		BirthYear:    2024,
		Breed:        "Any",
		Specie:       1,
	}

	resp, err := s.Create(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "save failed")
}

func TestPetServer_Update_Success(t *testing.T) {
	app := &appMock{
		updatePetFn: func(p *entity.Pet) (*entity.Pet, map[string]string) {
			ret := *p
			ret.NIdentification = 777
			ret.UuidGuardian = uuid.New()
			return &ret, nil
		},
	}
	s := NewPetServer(app)

	petID := uuid.New()
	req := &pb.UpdatePetRequest{
		Uuid:      petID.String(),
		Name:      "Luna Updated",
		BirthYear: 2021,
		Breed:     "Beagle Tricolor",
		Specie:    3,
	}

	resp, err := s.Update(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int64(777), resp.NIdentification)
	assert.Equal(t, petID.String(), resp.Uuid)
	assert.Equal(t, "Luna Updated", resp.Name)
	assert.Equal(t, uint64(2021), resp.BirthYear)
	assert.Equal(t, "Beagle Tricolor", resp.Breed)
	assert.Equal(t, "3", resp.Specie)
}

func TestPetServer_Update_Error(t *testing.T) {
	app := &appMock{
		updatePetFn: func(p *entity.Pet) (*entity.Pet, map[string]string) {
			return nil, map[string]string{"message": "update failed"}
		},
	}
	s := NewPetServer(app)

	req := &pb.UpdatePetRequest{
		Uuid:      uuid.New().String(),
		Name:      "X",
		BirthYear: 2000,
		Breed:     "Y",
		Specie:    1,
	}

	resp, err := s.Update(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "update failed")
}

func TestPetServer_Get_Success(t *testing.T) {
	pet := makePet()
	app := &appMock{
		getPetFn: func(id string) (*entity.Pet, map[string]string) {
			return pet, nil
		},
	}
	s := NewPetServer(app)

	resp, err := s.Get(context.Background(), &pb.GetPetRequest{Uuid: pet.Uuid.String()})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int64(pet.NIdentification), resp.NIdentification)
	assert.Equal(t, pet.Uuid.String(), resp.Uuid)
	assert.Equal(t, pet.UuidGuardian.String(), resp.UuidGuardian)
	assert.Equal(t, pet.Name, resp.Name)
	assert.Equal(t, uint64(pet.BirthYear), resp.BirthYear)
	assert.Equal(t, pet.Breed, resp.Breed)
	assert.Equal(t, "2", resp.Specie)
}

func TestPetServer_Get_Error(t *testing.T) {
	app := &appMock{
		getPetFn: func(id string) (*entity.Pet, map[string]string) {
			return nil, map[string]string{"message": "not found"}
		},
	}
	s := NewPetServer(app)

	resp, err := s.Get(context.Background(), &pb.GetPetRequest{Uuid: uuid.New().String()})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not found")
}

func TestPetServer_Delete_Success(t *testing.T) {
	app := &appMock{
		deletePetFn: func(guardian string) (map[string]string, map[string]string) {
			return map[string]string{"message": "2 pet(s) deletados!"}, nil
		},
	}
	s := NewPetServer(app)

	guardian := uuid.New().String()
	resp, err := s.Delete(context.Background(), &pb.DeletePetRequest{UuidGuardian: guardian})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "2 pet(s) deletados!", resp.Message)
}

func TestPetServer_Delete_Error(t *testing.T) {
	app := &appMock{
		deletePetFn: func(guardian string) (map[string]string, map[string]string) {
			return nil, map[string]string{"message": "no pets"}
		},
	}
	s := NewPetServer(app)

	resp, err := s.Delete(context.Background(), &pb.DeletePetRequest{UuidGuardian: uuid.New().String()})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "no pets")
}
