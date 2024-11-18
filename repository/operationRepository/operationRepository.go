package operationRepository

import (
	"database/sql"
	"errors"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
)

func GetOperationFromDB(db *sql.DB, operationType string) (*models.Operation, error) {
	var operation models.Operation

	err := db.QueryRow("SELECT id, type, cost, status FROM operations WHERE type = ? AND status = 'active'", operationType).
		Scan(&operation.ID, &operation.Type, &operation.Cost, &operation.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("operation type not found or inactive")
		}
		return nil, err
	}

	return &operation, nil
}
