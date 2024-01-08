package controller

import (
	"context"
	"encoding/json"
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

func (c *BorrowController) BorrowBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Borrow request received by the %s library.\n", os.Getenv("LOCAL_NAME"))
	var borrowDTO dto.BorrowDTO
	err := json.NewDecoder(r.Body).Decode(&borrowDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, isBadRequest := c.borrowService.CreateNewBorrow(context.Background(), borrowDTO)
	if err != nil {
		if isBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book borrowed successfully"))
}

func (c *BorrowController) ReturnBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Return request received by the %s library.\n", os.Getenv("LOCAL_NAME"))
	var returnDTO dto.ReturnDTO
	err := json.NewDecoder(r.Body).Decode(&returnDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, isBadRequest := c.borrowService.ReturnBorrow(context.Background(), returnDTO)
	if err != nil {
		if isBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Borrow successfully returned"))
}
