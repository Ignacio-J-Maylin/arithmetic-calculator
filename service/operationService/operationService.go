package operationService

import (
	"database/sql"
	"errors"
	"io"
	"math"
	"net/http"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/repository/operationRepository"
)

func Addition(a, b float64) float64 {
	return a + b
}

func Subtraction(a, b float64) float64 {
	return a - b
}

func Multiplication(a, b float64) float64 {
	return a * b
}

func Division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("negative number")
	}
	return math.Sqrt(a), nil
}

func RandomString() (string, error) {
	response, err := http.Get("https://www.random.org/strings/?num=1&len=10&digits=on&upperalpha=on&loweralpha=on&unique=on&format=plain&rnd=new")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	randomString, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(randomString), nil
}

func GetOperation(db *sql.DB, operationType string) (*models.Operation, error) {
	return operationRepository.GetOperationFromDB(db, operationType)
}
