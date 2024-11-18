package userRepository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password, status FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al encriptar la contraseña: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %v", err)
	}

	var userID int64
	result, err := tx.Exec("INSERT INTO users (username, password, status) VALUES (?, ?, ?)", username, string(hashedPassword), "active")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error al insertar el usuario: %v", err)
	}

	userID, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error al obtener el ID del usuario: %v", err)
	}

	_, err = tx.Exec("INSERT INTO balances (user_id, credits) VALUES (?, ?)", userID, 0)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error al insertar el balance: %v", err)
	}

	return tx.Commit()
}

func AddCredits(db *sql.DB, userID int64, creditsToAdd float64) error {
	result, err := db.Exec("UPDATE balances SET credits = credits + ? WHERE user_id = ?", creditsToAdd, userID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		_, err = db.Exec("INSERT INTO balances (user_id, credits) VALUES (?, ?)", userID, creditsToAdd)
		if err != nil {
			return err
		}
	}

	return nil
}
func RemoveCredits(db *sql.DB, userID int64, creditsToRemove float64) error {
	result, err := db.Exec("UPDATE balances SET credits = credits - ? WHERE user_id = ? AND credits >= ?", creditsToRemove, userID, creditsToRemove)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("insufficient credits")
	}

	return nil
}

func GetCredits(db *sql.DB, userID int64) (float64, error) {
	var credits float64
	err := db.QueryRow("SELECT credits FROM balances WHERE user_id = ?", userID).Scan(&credits)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return credits, nil
}

func GetAllOperations(db *sql.DB) ([]models.Operation, error) {
	query := "SELECT id, type FROM operations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []models.Operation
	for rows.Next() {
		var operation models.Operation
		if err := rows.Scan(&operation.ID, &operation.Type); err != nil {
			return nil, err
		}
		operations = append(operations, operation)
	}

	return operations, nil
}
