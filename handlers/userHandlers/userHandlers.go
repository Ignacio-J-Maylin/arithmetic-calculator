package userHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/models"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/service/authHelpers"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/service/operationService"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/service/recordService"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/service/userService"
)

type OperationRequest struct {
	OperationType string  `json:"operation_type"`
	A             float64 `json:"a"`
	B             float64 `json:"b,omitempty"`
}

func HandleCredits(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authHelpers.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			credits, err := userService.GetUserCredits(db, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]float64{"credits": credits})

		case http.MethodPut:
			var requestBody struct {
				Credits float64           `json:"credits"`
				Action  models.ActionType `json:"action"`
			}
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if requestBody.Action == models.AddAction {
				err := userService.AddCreditsToUser(db, userID, requestBody.Credits)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else if requestBody.Action == models.RemoveAction {
				err := userService.RemoveCreditsFromUser(db, userID, requestBody.Credits)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid action, use 'add' or 'remove'", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Credits updated successfully"})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func PerformOperation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authHelpers.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var req OperationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		operation, err := operationService.GetOperation(db, req.OperationType)
		if err != nil {
			http.Error(w, "Failed to retrieve operation", http.StatusInternalServerError)
			return
		}

		credits, err := userService.GetUserCredits(db, userID)
		if err != nil {
			http.Error(w, "Failed to retrieve user credits", http.StatusInternalServerError)
			return
		}
		if credits < operation.Cost {
			http.Error(w, "Insufficient credits", http.StatusPaymentRequired)
			return
		}

		var result interface{}
		switch req.OperationType {
		case "addition":
			result = operationService.Addition(req.A, req.B)
		case "subtraction":
			result = operationService.Subtraction(req.A, req.B)
		case "multiplication":
			result = operationService.Multiplication(req.A, req.B)
		case "division":
			result, err = operationService.Division(req.A, req.B)
		case "square_root":
			result, err = operationService.SquareRoot(req.A)
		case "random_string":
			result, err = operationService.RandomString()
		default:
			http.Error(w, "Invalid operation type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resultString string
		switch v := result.(type) {
		case string:
			resultString = v
		case float64:
			resultString = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			http.Error(w, "Unsupported result type", http.StatusInternalServerError)
			return
		}

		if err := userService.RemoveCreditsFromUser(db, userID, operation.Cost); err != nil {
			http.Error(w, "Failed to deduct credits", http.StatusInternalServerError)
			return
		}

		if err := recordService.CreateRecord(db, operation.ID, userID, operation.Cost, credits-operation.Cost, resultString); err != nil {
			http.Error(w, "Failed to record operation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	}
}

func GetRecordsHistory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authHelpers.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		operationName := r.URL.Query().Get("operation_name")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		orderBy := r.URL.Query().Get("order_by")
		orderDir := r.URL.Query().Get("order_dir")

		var filter models.RecordFilter
		filter.UserID = &userID

		if operationName != "" {
			filter.OperationName = &operationName
		}
		if startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err == nil {
				filter.StartDate = &startDate
			}
		}
		if endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil {
				filter.EndDate = &endDate
			}
		}
		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err == nil {
				filter.Limit = limit
			}
		} else {
			filter.Limit = 10
		}
		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err == nil {
				filter.Offset = offset
			}
		}

		filter.OrderBy = orderBy
		filter.OrderDir = orderDir

		records, totalRecords, err := recordService.GetFilteredRecords(db, filter)
		if err != nil {
			log.Printf("Error retrieving records: %v", err)
			http.Error(w, "Failed to retrieve records", http.StatusInternalServerError)
			return
		}

		totalPages := int(math.Ceil(float64(totalRecords) / float64(filter.Limit)))
		currentPage := (filter.Offset / filter.Limit) + 1

		response := models.PaginatedResponse{
			TotalRecords:   totalRecords,
			CurrentPage:    currentPage,
			TotalPages:     totalPages,
			RecordsPerPage: filter.Limit,
			Records:        records,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteRecordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authHelpers.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		recordIDStr := r.URL.Query().Get("record_id")
		if recordIDStr == "" {
			http.Error(w, "Record ID is required", http.StatusBadRequest)
			return
		}

		recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid record ID", http.StatusBadRequest)
			return
		}

		err = recordService.SoftDeleteRecord(db, recordID, userID)
		if err != nil {
			if err == models.ErrRecordNotFound {
				http.Error(w, "Record not found or unauthorized", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to delete record", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Record soft-deleted successfully"})
	}
}

func GetOperations(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operations, err := userService.GetAllOperations(db)
		if err != nil {
			log.Printf("Error retrieving operations: %v", err)
			http.Error(w, "Failed to retrieve operations", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"operations": operations,
		})
	}
}
