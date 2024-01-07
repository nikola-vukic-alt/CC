package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-app/local/dto"
	"library-app/local/model"
	"library-app/local/repository"
	"log"
	"net/http"

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

	Member := getMemberBySSN(borrowDTO.SSN, client)
	if Member.BorrowCnt > 2 {
		return errors.New("Member is already borrowing three books.")
	}

	newBorrow := model.Borrow{
		UserId: Member.Id,
		Title:  borrowDTO.Title,
		Author: borrowDTO.Author,
		ISBN:   borrowDTO.ISBN,
		From:   borrowDTO.From,
		To:     borrowDTO.To,
	}

	err := s.borrowRepo.SaveBorrow(ctx, newBorrow)
	if err != nil {
		log.Printf("Error registering borrow: %v\n", err)
		return err
	}

	return nil
}

func isInvalidDTO(borrowDTO dto.BorrowDTO) bool {
	titleMissing := len(borrowDTO.Title) == 0
	ssnMissing := len(borrowDTO.SSN) == 0
	fromMissing := borrowDTO.From.IsZero()
	toMissing := borrowDTO.To.IsZero()
	return titleMissing || ssnMissing || fromMissing || toMissing
}

func getMemberBySSN(ssn string, client *http.Client) Member {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/get?ssn=%s", ssn), nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return Member{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return Member{}
	}

	var m Member
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		fmt.Println("Error decoding response body:", err)
		return Member{}
	}

	return m
}
