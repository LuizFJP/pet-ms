package grpc

import (
	"context"
	"github.com/LuizFJP/pet-ms/application"
	"github.com/LuizFJP/pet-ms/domain/entity"
	pb "github.com/LuizFJP/pet-ms/proto"
	"github.com/google/uuid"
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

	s.pa.SavePet(petEntity)
	return nil, nil
}

func (s *PetServer) Update(ctx context.Context, input *pb.UpdatePetRequest) (*pb.UpdatePetResponse, error) {
	return nil, nil
}

func (s *PetServer) Get(ctx context.Context, input *pb.GetPetRequest) (*pb.GetPetResponse, error) {
	return nil, nil
}

func (s *PetServer) Delete(ctx context.Context, input *pb.DeletePetRequest) (*pb.DeletePetResponse, error) {
	return nil, nil
}
