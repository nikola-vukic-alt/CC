package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"library-app/local/dto"
	"library-app/local/service"
	"log"
	"net/http"
	"os"
)

// BorrowController handles HTTP requests related to borrows.
type BorrowController struct {
	borrowService *service.BorrowService
}

// NewBorrowController creates a new BorrowController instance.
func NewBorrowController(borrowService *service.BorrowService) *BorrowController {
	return &BorrowController{
		borrowService: borrowService,
	}
}

func (c *BorrowController) RegisterMember(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register request received by the %s library.\n", os.Getenv("LOCAL_NAME"))
	var registrationDTO dto.RegistrationDTO
	err := json.NewDecoder(r.Body).Decode(&registrationDTO)
	if err != nil {
		http.Error(w, "Invalid request body\n", http.StatusBadRequest)
		return
	}

	jsonBody, err := json.Marshal(registrationDTO)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	endpoint := fmt.Sprintf("http://%s:8080/register", os.Getenv("CENTRAL_LIBRARY"))

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(bodyBytes)
}

func (c *BorrowController) BorrowBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Borrow request received by the %s library.\n", os.Getenv("LOCAL_NAME"))
	log.Println("Trying to create a new borrow 1...")
	var borrowDTO dto.BorrowDTO
	err := json.NewDecoder(r.Body).Decode(&borrowDTO)
	if err != nil {
		log.Printf("Error: SSN - %s\n", borrowDTO.SSN)
		http.Error(w, "Invalid request body\n", http.StatusBadRequest)
		return
	}
	log.Println("Trying to create a new borrow 2...")
	err, statusCode, newBorrow := c.borrowService.CreateNewBorrow(context.Background(), borrowDTO)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	responseJSON, err := json.Marshal(newBorrow)
	if err != nil {
		http.Error(w, "Error encoding response body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func (c *BorrowController) ReturnBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Return request received by the %s library.\n", os.Getenv("LOCAL_NAME"))
	var returnDTO dto.ReturnDTO
	err := json.NewDecoder(r.Body).Decode(&returnDTO)
	if err != nil {
		http.Error(w, "Invalid request body\n", http.StatusBadRequest)
		return
	}

	err, statusCode, newReturn := c.borrowService.ReturnBorrow(context.Background(), returnDTO)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	responseJSON, err := json.Marshal(newReturn)
	if err != nil {
		http.Error(w, "Error encoding response body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
