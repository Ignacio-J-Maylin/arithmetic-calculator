package models

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID       int64
	Username string
	Password string
	Status   string
}

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type Balance struct {
	ID      int64
	UserID  int64
	Credits float64
}

type Record struct {
	ID                int64     `json:"id"`
	OperationName     string    `json:"operation_name"`
	UserID            int64     `json:"user_id"`
	Amount            float64   `json:"amount"`
	UserBalance       float64   `json:"user_balance"`
	OperationResponse string    `json:"operation_response"`
	Date              time.Time `json:"date"`
}

type Operation struct {
	ID     int64
	Type   string
	Status string
	Cost   float64
}

const (
	OperationAddition       = "addition"
	OperationSubtraction    = "subtraction"
	OperationMultiplication = "multiplication"
	OperationDivision       = "division"
	OperationSquareRoot     = "square_root"
	OperationRandomString   = "random_string"
)

type ActionType string

const (
	AddAction    ActionType = "add"
	RemoveAction ActionType = "remove"
)

type PaginatedResponse struct {
	TotalRecords   int      `json:"total_records"`
	CurrentPage    int      `json:"current_page"`
	TotalPages     int      `json:"total_pages"`
	RecordsPerPage int      `json:"records_per_page"`
	Records        []Record `json:"records"`
}

type RecordFilter struct {
	UserID        *int64     `json:"user_id"`
	OperationName *string    `json:"operation_name"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
	OrderBy       string     `json:"order_by"`
	OrderDir      string     `json:"order_dir"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var ErrRecordNotFound = errors.New("record not found")

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
