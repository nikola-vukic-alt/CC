package service

import (
	"context"
	"errors"
	"fmt"
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
	nameMissing := len(registrationDTO.Name) == 0
	surnameMissing := len(registrationDTO.Surname) == 0
	addressMissing := len(registrationDTO.Address) == 0
	ssnMissing := len(registrationDTO.SSN) == 0
	if nameMissing || surnameMissing || addressMissing || ssnMissing {
		return errors.New("All the fields are required.")
	}

	existingMember, err := s.memberRepo.GetMemberBySSN(ctx, registrationDTO.SSN)
	if existingMember.SSN == registrationDTO.SSN {
		return fmt.Errorf("Member with SSN: %s already exists.", registrationDTO.SSN)
	}

	newMember := model.Member{
		Name:      registrationDTO.Name,
		Surname:   registrationDTO.Surname,
		Address:   registrationDTO.Address,
		SSN:       registrationDTO.SSN,
		BorrowCnt: 0,
	}

	err = s.memberRepo.SaveMember(ctx, newMember)
	if err != nil {
		log.Printf("Error registering member: %v\n", err)
		return err
	}

	return nil
}
