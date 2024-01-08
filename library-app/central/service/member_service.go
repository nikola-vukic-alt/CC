package service

import (
	"context"
	"errors"
	"fmt"
	"library-app/central/dto"
	"library-app/central/model"
	"library-app/central/repository"
	"log"
	"net/http"
)

type MemberService struct {
	memberRepo *repository.MemberRepository
}

func NewMemberService(memberRepo *repository.MemberRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
	}
}

func (s *MemberService) RegisterMember(ctx context.Context, registrationDTO dto.RegistrationDTO) (error, int, dto.MemberDTO) {

	if isInvalidDTO(registrationDTO) {
		return errors.New("All the fields are required."), http.StatusBadRequest, dto.MemberDTO{}
	}

	_, err, statusCode := s.memberRepo.GetMemberBySSN(ctx, registrationDTO.SSN)
	if err != nil && statusCode != http.StatusNotFound {
		return err, statusCode, dto.MemberDTO{}
	}
	if statusCode == http.StatusOK {
		return fmt.Errorf("Member with SSN: %s already exists.", registrationDTO.SSN), http.StatusBadRequest, dto.MemberDTO{}
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
		return err, http.StatusInternalServerError, dto.MemberDTO{}
	}
	log.Printf("Registered new member: %s %s - SSN: %s\n", newMember.Name, newMember.Surname, newMember.SSN)
	return nil, http.StatusCreated, dto.MemberDTO{
		Name:      newMember.Name,
		Surname:   newMember.Surname,
		Address:   newMember.Address,
		SSN:       newMember.SSN,
		BorrowCnt: newMember.BorrowCnt,
	}
}

func (s *MemberService) GetMemberBySSN(ctx context.Context, ssn string) (model.Member, error, int) {
	return s.memberRepo.GetMemberBySSN(ctx, ssn)
}

func (s *MemberService) UpdateBorrowCount(ctx context.Context, updateDTO dto.UpdateDTO) (error, int) {
	member, err, statusCode := s.memberRepo.GetMemberBySSN(ctx, updateDTO.SSN)
	if err != nil {
		return errors.New("Member not found"), statusCode
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
