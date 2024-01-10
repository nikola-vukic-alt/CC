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
	"os"
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
func (s *BorrowService) CreateNewBorrow(ctx context.Context, borrowDTO dto.BorrowDTO) (error, int, dto.DetailedBorrowDTO) {
	if isInvalidDTO(borrowDTO) {
		return errors.New("All the fields are required."), http.StatusBadRequest, dto.DetailedBorrowDTO{}
	}

	client := &http.Client{}

	member, err, statusCode := getMemberBySSN(borrowDTO.SSN, client)
	if err != nil {
		return errors.New("Member not registererd"), http.StatusBadRequest, dto.DetailedBorrowDTO{}
	}
	if member.BorrowCnt > 2 {
		return errors.New("Member is already borrowing three books."), http.StatusBadRequest, dto.DetailedBorrowDTO{}
	}
	_, err, statusCode = s.borrowRepo.GetMembersBorrow(ctx, member.Id, borrowDTO.Title)
	if err != nil && statusCode != http.StatusNotFound {
		return err, statusCode, dto.DetailedBorrowDTO{}
	}
	if statusCode == http.StatusOK {
		return errors.New("You have already borrowed this book."), http.StatusBadRequest, dto.DetailedBorrowDTO{}
	}
	err, statusCode = updateBorrowCount(member, client, true)
	if err != nil {
		return err, statusCode, dto.DetailedBorrowDTO{}
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
		return err, http.StatusInternalServerError, dto.DetailedBorrowDTO{}
	}
	log.Printf("Member: %s %s, book title: %s - borrow count: %d\n",
		member.Name,
		member.Surname,
		newBorrow.Title,
		member.BorrowCnt+1)

	return nil, http.StatusOK, dto.DetailedBorrowDTO{
		SSN:       member.SSN,
		BorrowCnt: member.BorrowCnt + 1,
		Title:     newBorrow.Title,
		Author:    newBorrow.Author,
		ISBN:      newBorrow.ISBN,
		From:      time.Now(),
	}
}

// return type: (error, statusCode)
func (s *BorrowService) ReturnBorrow(ctx context.Context, returnDTO dto.ReturnDTO) (error, int, dto.DetailedReturnDTO) {
	if isInvalidReturn(returnDTO) {
		return errors.New("All the fields are required\n"), http.StatusBadRequest, dto.DetailedReturnDTO{}
	}
	client := &http.Client{}
	member, err, _ := getMemberBySSN(returnDTO.SSN, client)
	if err != nil {
		return errors.New("Member not found\n"), http.StatusBadRequest, dto.DetailedReturnDTO{}
	}
	borrow, err, _ := s.borrowRepo.GetMembersBorrow(ctx, member.Id, returnDTO.Title)
	if err != nil {
		return errors.New("Borrow not found\n"), http.StatusBadRequest, dto.DetailedReturnDTO{}
	}
	if borrow.From.Before(borrow.To) {
		return errors.New("You have already returned this borrow\n"), http.StatusBadRequest, dto.DetailedReturnDTO{}
	}
	borrow.To = time.Now()
	err = s.borrowRepo.UpdateBorrow(ctx, borrow.Id, borrow)
	if err != nil {
		return err, http.StatusInternalServerError, dto.DetailedReturnDTO{}
	}
	err, statusCode := updateBorrowCount(member, client, false)
	if err != nil {
		return err, statusCode, dto.DetailedReturnDTO{}
	}
	log.Printf("Member: %s %s, book title: %s - borrow count: %d\n",
		member.Name,
		member.Surname,
		borrow.Title,
		member.BorrowCnt-1)
	return nil, http.StatusOK, dto.DetailedReturnDTO{
		SSN:       member.SSN,
		BorrowCnt: member.BorrowCnt - 1,
		Title:     borrow.Title,
		Author:    borrow.Author,
		ISBN:      borrow.ISBN,
		From:      borrow.From,
		To:        time.Now(),
	}
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
	host := os.Getenv("CENTRAL_LIBRARY")
	log.Printf("Central host: %s\n", host)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:8080/get?ssn=%s", ssn, host), nil)
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

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:8080/update-borrow-count", os.Getenv("CENTRAL_LIBRARY")), bytes.NewBuffer(requestBody))
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
