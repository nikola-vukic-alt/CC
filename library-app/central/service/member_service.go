package service

import (
	"context"
	"library-app/central/dto"
	"library-app/central/model"
	"library-app/central/repository"
	"log"
)

type MemberService struct {
	memberRepo *repository.MemberRepository
}

func NewMemberService(memberRepo *repository.MemberRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
	}
}

func (s *MemberService) RegisterMember(ctx context.Context, registrationDTO dto.RegistrationDTO) error {
	// You can perform validation on registrationDTO fields if needed

	newMember := model.Member{
		Name:    registrationDTO.Name,
		Surname: registrationDTO.Surname,
		Address: registrationDTO.Address,
		SSN:     registrationDTO.SSN,
	}

	err := s.memberRepo.SaveMember(ctx, newMember)
	if err != nil {
		log.Printf("Error registering member: %v\n", err)
		return err
	}

	return nil
}
