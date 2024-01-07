package controller

import (
	"context"
	"encoding/json"
	"library-app/central/dto"
	"library-app/central/service" // Update this import path based on your project structure
	"net/http"
)

// MemberController handles HTTP requests related to members.
type MemberController struct {
	memberService *service.MemberService
}

// NewMemberController creates a new MemberController instance.
func NewMemberController(memberService *service.MemberService) *MemberController {
	return &MemberController{
		memberService: memberService,
	}
}

// Register is an HTTP handler function to handle member registration requests.
func (c *MemberController) Register(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body into RegistrationDTO
	var registrationDTO dto.RegistrationDTO
	err := json.NewDecoder(r.Body).Decode(&registrationDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register the member using the MemberService
	err = c.memberService.RegisterMember(context.Background(), registrationDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Member successfully registered"))
}
