package dto

import (
	"time"
)

type BorrowDTO struct {
	SSN    string
	Title  string
	Author string
	ISBN   string
	From   time.Time
	To     time.Time
}
