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

// return type: (error, statusCode)
func (s *BorrowService) CreateNewBorrow(ctx context.Context, borrowDTO dto.BorrowDTO) (error, int) {
	if isInvalidDTO(borrowDTO) {
		return errors.New("All the fields are required."), http.StatusBadRequest
	}

	client := &http.Client{}

	member, err, statusCode := getMemberBySSN(borrowDTO.SSN, client)
	if err != nil {
		return errors.New("Member not registererd"), http.StatusBadRequest
	}
	if member.BorrowCnt > 2 {
		return errors.New("Member is already borrowing three books."), http.StatusBadRequest
	}
	_, err, statusCode = s.borrowRepo.GetMembersBorrow(ctx, member.Id, borrowDTO.Title)
	if err != nil && statusCode != http.StatusNotFound {
		return err, statusCode
	}
	if statusCode == http.StatusOK {
		return errors.New("You have already borrowed this book."), http.StatusBadRequest
	}
	err, statusCode = updateBorrowCount(member, client, true)
	if err != nil {
		return err, statusCode
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
		return err, http.StatusInternalServerError
	}
	log.Printf("Member: %s %s, book title: %s - borrow count: %d\n",
		member.Name,
		member.Surname,
		newBorrow.Title,
		member.BorrowCnt+1)
	return nil, http.StatusOK
}

// return type: (error, statusCode)
func (s *BorrowService) ReturnBorrow(ctx context.Context, returnDTO dto.ReturnDTO) (error, int) {
	if isInvalidReturn(returnDTO) {
		return errors.New("All the fields are required\n"), http.StatusBadRequest
	}
	client := &http.Client{}
	member, err, _ := getMemberBySSN(returnDTO.SSN, client)
	if err != nil {
		return errors.New("Member not found\n"), http.StatusBadRequest
	}
	borrow, err, _ := s.borrowRepo.GetMembersBorrow(ctx, member.Id, returnDTO.Title)
	if err != nil {
		return errors.New("Borrow not found\n"), http.StatusBadRequest
	}
	if borrow.From.Before(borrow.To) {
		return errors.New("You have already returned this borrow\n"), http.StatusBadRequest
	}
	borrow.To = time.Now()
	err = s.borrowRepo.UpdateBorrow(ctx, borrow.Id, borrow)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	err, statusCode := updateBorrowCount(member, client, false)
	if err != nil {
		return err, statusCode
	}
	log.Printf("Member: %s %s, book title: %s - borrow count: %d\n",
		member.Name,
		member.Surname,
		borrow.Title,
		member.BorrowCnt-1)
	return nil, http.StatusOK
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

func getMemberBySSN(ssn string, client *http.Client) (Member, error, int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://central_library:8080/get?ssn=%s", ssn), nil)
	if err != nil {
		return Member{}, fmt.Errorf("Error creating HTTP request: %v\n", err), http.StatusInternalServerError
	}

	resp, err := client.Do(req)
	if err != nil {
		return Member{}, fmt.Errorf("Error sending HTTP request: %v\n", err), resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Member{}, fmt.Errorf("Unexpected status code: %v\n", resp.StatusCode), resp.StatusCode
	}

	var m Member
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return Member{}, fmt.Errorf("Error decoding response body: %v\n", err), http.StatusInternalServerError
	}

	return m, nil, http.StatusOK
}

func updateBorrowCount(member Member, client *http.Client, shouldIncrease bool) (error, int) {
	var updateDTO UpdateDTO
	updateDTO.SSN = member.SSN
	if shouldIncrease {
		updateDTO.NewCount = member.BorrowCnt + 1
	} else {
		updateDTO.NewCount = member.BorrowCnt - 1
	}
	requestBody, err := json.Marshal(updateDTO)
	if err != nil {
		return fmt.Errorf("Error encoding UpdateDTO into JSON: %v\n", err), http.StatusInternalServerError
	}

	req, err := http.NewRequest("PUT", "http://central_library:8080/update-borrow-count", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %v\n", err), http.StatusInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending HTTP request: %v\n", err), resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %v\n", resp.StatusCode), resp.StatusCode
	}
	return nil, http.StatusOK
}
