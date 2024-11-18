package userService

import (
	"database/sql"
	"errors"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/repository/userRepository"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, username, password string) error {
	_, err := userRepository.GetUserByUsername(db, username)
	if err == nil {
		return errors.New("user already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return userRepository.CreateUser(db, username, password)
}

func AuthenticateUser(db *sql.DB, username, password string) (int64, bool, error) {
	user, err := userRepository.GetUserByUsername(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, false, nil
	}

	return user.ID, true, nil
}

func AddCreditsToUser(db *sql.DB, userID int64, credits float64) error {
	if credits <= 0 {
		return errors.New("credits must be greater than zero")
	}

	return userRepository.AddCredits(db, userID, credits)
}

func RemoveCreditsFromUser(db *sql.DB, userID int64, credits float64) error {
	if credits <= 0 {
		return errors.New("credits to remove must be greater than zero")
	}

	return userRepository.RemoveCredits(db, userID, credits)
}

func GetUserCredits(db *sql.DB, userID int64) (float64, error) {
	return userRepository.GetCredits(db, userID)
}

func GetAllOperations(db *sql.DB) ([]models.Operation, error) {
	return userRepository.GetAllOperations(db)
}
