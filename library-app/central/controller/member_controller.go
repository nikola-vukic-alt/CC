package controller

import (
	"context"
	"encoding/json"
	"library-app/central/dto"
	"library-app/central/service"
	"log"
	"net/http"
)

type MemberController struct {
	memberService *service.MemberService
}

func NewMemberController(memberService *service.MemberService) *MemberController {
	return &MemberController{
		memberService: memberService,
	}
}

func (c *MemberController) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Registration request received by the central library.")
	var registrationDTO dto.RegistrationDTO
	err := json.NewDecoder(r.Body).Decode(&registrationDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.memberService.RegisterMember(context.Background(), registrationDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Member successfully registered"))
}

func (c *MemberController) GetMemberBySSN(w http.ResponseWriter, r *http.Request) {
	ssn := r.URL.Query().Get("ssn")
	if ssn == "" {
		http.Error(w, "SSN parameter is required", http.StatusBadRequest)
		return
	}

	member, err := c.memberService.GetMemberBySSN(context.Background(), ssn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (c *MemberController) UpdateBorrowCount(w http.ResponseWriter, r *http.Request) {
	var updateDTO dto.UpdateDTO

	err := json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.memberService.UpdateBorrowCount(context.Background(), updateDTO)
	if err != nil {
		log.Printf("Error updating borrow count: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Borrow count updated successfully"))
}
