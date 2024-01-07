package controller

import (
	"context"
	"encoding/json"
	"library-app/local/dto"
	"library-app/local/service"
	"net/http"
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

// Register is an HTTP handler function to handle borrow requests.
func (c *BorrowController) BorrowBook(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body into BorrowDTO
	var borrowDTO dto.BorrowDTO
	err := json.NewDecoder(r.Body).Decode(&borrowDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register the borrow using the BorrowService
	err = c.borrowService.CreateNewBorrow(context.Background(), borrowDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book successfully borrowed"))
}
