package repository

import (
	"database/sql"
	"fmt"
)

type CasinoResultRepository interface {
	GetCasinoResults(casinoType string, gameDate string, search string, length int, start int) ([]map[string]interface{}, int, error)
	GetCasinoResultCards(casinoType string, eventID string) (string, string, string, error)
}

type casinoResultRepository struct {
	db *sql.DB
}

func NewCasinoResultRepository(db *sql.DB) CasinoResultRepository {
	return &casinoResultRepository{db: db}
}

func (r *casinoResultRepository) GetCasinoResults(casinoType string, gameDate string, search string, length int, start int) ([]map[string]interface{}, int, error) {
	searchCond := ""
	var args []interface{}
	args = append(args, casinoType, gameDate)

	if search != "" {
		searchCond = " AND (event_id LIKE ? OR game_type LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	countQuery := `SELECT COUNT(*) FROM twenty_teenpatti_result WHERE game_type=? AND DATE(result_time) = ?` + searchCond
	var totalData int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `SELECT event_id, game_type, result_status, cards FROM twenty_teenpatti_result WHERE game_type=? AND DATE(result_time) = ?` + searchCond + ` ORDER BY result_time DESC`

	if length != -1 {
		dataQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var eventID, gameType, resultStatus, cards string
		if err := rows.Scan(&eventID, &gameType, &resultStatus, &cards); err != nil {
			return nil, 0, err
		}
		results = append(results, map[string]interface{}{
			"event_id":      eventID,
			"game_type":     gameType,
			"result_status": resultStatus,
			"cards":         cards,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, totalData, nil
}

func (r *casinoResultRepository) GetCasinoResultCards(casinoType string, eventID string) (string, string, string, error) {
	query := `SELECT cards, result_status, desc_remakrs FROM twenty_teenpatti_result WHERE game_type=? AND event_id=?`
	var cards, resultStatus string
	var descRemarks sql.NullString
	err := r.db.QueryRow(query, casinoType, eventID).Scan(&cards, &resultStatus, &descRemarks)
	if err != nil {
		return "", "", "", err
	}
	return cards, resultStatus, descRemarks.String, nil
}
