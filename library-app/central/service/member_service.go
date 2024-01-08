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

	if isInvalidDTO(registrationDTO) {
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
	log.Printf("Registered new member: %s %s - SSN: %s\n", newMember.Name, newMember.Surname, newMember.SSN)
	return nil
}

func (s *MemberService) GetMemberBySSN(ctx context.Context, ssn string) (model.Member, error) {
	return s.memberRepo.GetMemberBySSN(ctx, ssn)
}

func (s *MemberService) UpdateBorrowCount(ctx context.Context, updateDTO dto.UpdateDTO) error {
	member, err := s.memberRepo.GetMemberBySSN(ctx, updateDTO.SSN)
	if err != nil {
		return err
	}
	member.BorrowCnt = updateDTO.NewCount
	return s.memberRepo.UpdateMember(ctx, member.Id, member)
}

func isInvalidDTO(registrationDTO dto.RegistrationDTO) bool {
	nameMissing := len(registrationDTO.Name) == 0
	surnameMissing := len(registrationDTO.Surname) == 0
	addressMissing := len(registrationDTO.Address) == 0
	ssnMissing := len(registrationDTO.SSN) == 0
	return nameMissing || surnameMissing || addressMissing || ssnMissing
}
