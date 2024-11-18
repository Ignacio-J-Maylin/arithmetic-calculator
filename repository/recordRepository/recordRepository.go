package recordRepository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
)

func GetRecords(db *sql.DB, filter models.RecordFilter) ([]models.Record, int, error) {
	query := `
		SELECT r.id, o.type AS operation_name, r.user_id, r.amount, r.user_balance, r.operation_response, r.date
		FROM records r
		JOIN operations o ON r.operation_id = o.id
		WHERE r.deleted_at IS NULL`
	countQuery := `
		SELECT COUNT(*)
		FROM records r
		JOIN operations o ON r.operation_id = o.id
		WHERE r.deleted_at IS NULL`
	args := []interface{}{}

	if filter.UserID != nil {
		query += " AND r.user_id = ?"
		countQuery += " AND r.user_id = ?"
		args = append(args, *filter.UserID)
	}
	if filter.OperationName != nil {
		query += " AND LOWER(o.type) LIKE LOWER(?)"
		countQuery += " AND LOWER(o.type) LIKE LOWER(?)"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.OperationName))
	}
	if filter.StartDate != nil {
		query += " AND r.date >= ?"
		countQuery += " AND r.date >= ?"
		args = append(args, *filter.StartDate)
	}
	if filter.EndDate != nil {
		query += " AND r.date <= ?"
		countQuery += " AND r.date <= ?"
		args = append(args, *filter.EndDate)
	}

	var totalRecords int
	if err := db.QueryRow(countQuery, args...).Scan(&totalRecords); err != nil {
		return nil, 0, err
	}

	if filter.OrderBy != "" {
		direction := "asc"
		if filter.OrderDir == "desc" {
			direction = "desc"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.OrderBy, direction)
	} else {
		query += " ORDER BY r.date DESC"
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, filter.Offset)

	log.Printf("Executing query: %s", query)
	log.Printf("With arguments: %v", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	records := []models.Record{}
	for rows.Next() {
		var record models.Record
		if err := rows.Scan(&record.ID, &record.OperationName, &record.UserID, &record.Amount, &record.UserBalance, &record.OperationResponse, &record.Date); err != nil {
			return nil, 0, err
		}
		records = append(records, record)
	}

	return records, totalRecords, nil
}

func CreateRecord(db *sql.DB, operationID int64, userID int64, amount float64, userBalance float64, operationResponse string) error {
	_, err := db.Exec(`
		INSERT INTO records (operation_id, user_id, amount, user_balance, operation_response, date) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		operationID, userID, amount, userBalance, operationResponse, time.Now(),
	)
	return err
}

func SoftDeleteRecord(db *sql.DB, recordID int64, userID int64) error {
	query := "UPDATE records SET deleted_at = ? WHERE id = ? AND user_id = ? AND deleted_at IS NULL"
	result, err := db.Exec(query, time.Now(), recordID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}
