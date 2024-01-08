package dto

import "time"

type DetailedBorrowDTO struct {
	SSN       string
	BorrowCnt int
	Title     string
	Author    string
	ISBN      string
	From      time.Time
}

type DetailedReturnDTO struct {
	SSN       string
	BorrowCnt int
	Title     string
	Author    string
	ISBN      string
	From      time.Time
	To        time.Time
}
