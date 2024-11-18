package recordService

import (
	"database/sql"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/repository/recordRepository"
)

func CreateRecord(db *sql.DB, operationID int64, userID int64, amount float64, userBalance float64, operationResponse string) error {
	return recordRepository.CreateRecord(db, operationID, userID, amount, userBalance, operationResponse)
}

func GetFilteredRecords(db *sql.DB, filter models.RecordFilter) ([]models.Record, int, error) {
	return recordRepository.GetRecords(db, filter)
}

func SoftDeleteRecord(db *sql.DB, recordID int64, userID int64) error {
	return recordRepository.SoftDeleteRecord(db, recordID, userID)
}
