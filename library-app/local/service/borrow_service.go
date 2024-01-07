package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-app/local/dto"
	"library-app/local/model"
	"library-app/local/repository"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	Id        primitive.ObjectID
	Name      string
	Surname   string
	Address   string
	SSN       string
	BorrowCnt int
}

type UpdateDTO struct {
	SSN      string
	NewCount int
}

type BorrowService struct {
	borrowRepo *repository.BorrowRepository
}

func NewBorrowService(borrowRepo *repository.BorrowRepository) *BorrowService {
	return &BorrowService{
		borrowRepo: borrowRepo,
	}
}

func (s *BorrowService) CreateNewBorrow(ctx context.Context, borrowDTO dto.BorrowDTO) error {
	if isInvalidDTO(borrowDTO) {
		return errors.New("All the fields are required.")
	}

	client := &http.Client{}

	member, err := getMemberBySSN(borrowDTO.SSN, client)
	if err != nil {
		return err
	}
	if member.BorrowCnt > 2 {
		return errors.New("Member is already borrowing three books.")
	}

	err = updateBorrowCount(member, client, true)
	if err != nil {
		return err
	}

	newBorrow := model.Borrow{
		UserId: member.Id,
		Title:  borrowDTO.Title,
		Author: borrowDTO.Author,
		ISBN:   borrowDTO.ISBN,
		From:   time.Now(),
	}

	err = s.borrowRepo.SaveBorrow(ctx, newBorrow)
	if err != nil {
		log.Printf("Error registering borrow: %v\n", err)
		return err
	}

	return nil
}

func (s *BorrowService) ReturnBorrow(ctx context.Context, returnDTO dto.ReturnDTO) error {
	if isInvalidReturn(returnDTO) {
		return errors.New("All the fields are required.")
	}
	client := &http.Client{}
	member, err := getMemberBySSN(returnDTO.SSN, client)
	if err != nil {
		return err
	}
	borrow, err := s.borrowRepo.GetMembersBorrow(ctx, member.Id, returnDTO.Title)
	if err != nil {
		return err
	}
	borrow.To = time.Now()
	err = s.borrowRepo.UpdateBorrow(ctx, borrow.Id, borrow)
	if err != nil {
		return err
	}
	err = updateBorrowCount(member, client, false)
	if err != nil {
		return err
	}
	return nil
}

func isInvalidReturn(returnDTO dto.ReturnDTO) bool {
	ssnMissing := len(returnDTO.SSN) == 0
	titleMissing := len(returnDTO.Title) == 0
	return titleMissing || ssnMissing
}

func isInvalidDTO(borrowDTO dto.BorrowDTO) bool {
	ssnMissing := len(borrowDTO.SSN) == 0
	titleMissing := len(borrowDTO.Title) == 0
	authorMissing := len(borrowDTO.Author) == 0
	isbnMissing := len(borrowDTO.ISBN) == 0
	return titleMissing || ssnMissing || authorMissing || isbnMissing
}

func getMemberBySSN(ssn string, client *http.Client) (Member, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/get?ssn=%s", ssn), nil)
	if err != nil {
		return Member{}, fmt.Errorf("Error creating HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return Member{}, fmt.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Member{}, fmt.Errorf("Unexpected status code: %v", resp.StatusCode)
	}

	var m Member
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return Member{}, fmt.Errorf("Error decoding response body: %v", err)
	}

	return m, nil
}

func updateBorrowCount(member Member, client *http.Client, shouldIncrease bool) error {
	var updateDTO UpdateDTO
	updateDTO.SSN = member.SSN
	if shouldIncrease {
		updateDTO.NewCount = member.BorrowCnt + 1
	} else {
		updateDTO.NewCount = member.BorrowCnt - 1
	}
	requestBody, err := json.Marshal(updateDTO)
	if err != nil {
		return fmt.Errorf("Error encoding UpdateDTO into JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", "http://localhost:8080/update-borrow-count", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("New borrow count is %d.\n", updateDTO.NewCount)
	return nil
}
