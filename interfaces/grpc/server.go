package grpc

import (
	"context"
	"fmt"
	"github.com/LuizFJP/pet-ms/application"
	"github.com/LuizFJP/pet-ms/domain/entity"
	pb "github.com/LuizFJP/pet-ms/proto"
	"github.com/google/uuid"
	"strconv"
)

type PetServer struct {
	pa application.PetApplicationInterface
	pb.UnimplementedPetServiceServer
}

func NewPetServer(pa application.PetApplicationInterface) *PetServer {
	return &PetServer{pa: pa}
}

func (s *PetServer) Create(ctx context.Context, input *pb.CreatePetRequest) (*pb.CreatePetResponse, error) {
	petEntity := &entity.Pet{
		Name:         input.Name,
		Uuid:         uuid.New(),
		UuidGuardian: uuid.MustParse(input.UuidGuardian),
		BirthYear:    int(input.BirthYear),
		Breed:        input.Breed,
		Specie:       entity.PetType(input.Specie),
	}
	petEntity.Validate("default")

	res, errData := s.pa.SavePet(petEntity)
	if errData != nil {
		return nil, fmt.Errorf("something went wrong: %v", errData["message"])
	}

	petResponse := &pb.CreatePetResponse{
		NIdentification: int64(res.NIdentification),
		Uuid:            res.Uuid.String(),
		UuidGuardian:    res.UuidGuardian.String(),
		Name:            res.Name,
		BirthYear:       uint64(res.BirthYear),
		Breed:           res.Breed,
		Specie:          strconv.FormatInt(int64(res.Specie), 10),
	}

	return petResponse, nil
}

func (s *PetServer) Update(ctx context.Context, input *pb.UpdatePetRequest) (*pb.UpdatePetResponse, error) {
	petEntity := &entity.Pet{
		Uuid:      uuid.MustParse(input.Uuid),
		Name:      input.Name,
		BirthYear: int(input.BirthYear),
		Breed:     input.Breed,
		Specie:    entity.PetType(input.Specie),
	}
	petEntity.Validate("default")

	res, errData := s.pa.UpdatePet(petEntity)
	if errData != nil {
		return nil, fmt.Errorf("something went wrong: %v", errData["message"])
	}

	petResponse := &pb.UpdatePetResponse{
		NIdentification: int64(res.NIdentification),
		Uuid:            res.Uuid.String(),
		UuidGuardian:    res.UuidGuardian.String(),
		Name:            res.Name,
		BirthYear:       uint64(res.BirthYear),
		Breed:           res.Breed,
		Specie:          strconv.FormatInt(int64(res.Specie), 10),
	}
	return petResponse, nil
}

func (s *PetServer) Get(ctx context.Context, input *pb.GetPetRequest) (*pb.GetPetResponse, error) {
	res, errData := s.pa.GetPet(input.Uuid)
	if errData != nil {
		return nil, fmt.Errorf("something went wrong: %v", errData["message"])
	}
	petResponse := &pb.GetPetResponse{
		NIdentification: int64(res.NIdentification),
		Uuid:            res.Uuid.String(),
		UuidGuardian:    res.UuidGuardian.String(),
		Name:            res.Name,
		BirthYear:       uint64(res.BirthYear),
		Breed:           res.Breed,
		Specie:          strconv.FormatInt(int64(res.Specie), 10),
	}

	return petResponse, nil
}

func (s *PetServer) Delete(ctx context.Context, input *pb.DeletePetRequest) (*pb.DeletePetResponse, error) {
	res, errData := s.pa.DeletePet(input.UuidGuardian)
	if errData != nil {
		return nil, fmt.Errorf("something went wrong: %v", errData["message"])
	}

	deleteResponse := &pb.DeletePetResponse{
		Message: res["message"],
	}

	return deleteResponse, nil
}
