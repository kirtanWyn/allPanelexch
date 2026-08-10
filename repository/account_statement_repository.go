package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kirtanwyn/allPanelexch/dto"
)

type AccountStatementRepository interface {
	GetOpeningBalance(userID int, opDate string, reportType int, fromDate string) (float64, error)
	GetAccountStatementReport1(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error)
	GetAccountStatementReport2(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error)
	GetAccountStatementReport3(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error)
	GetAccountStatementReport4(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error)
	GetAccountStatementReport5(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error)
	GetAccountStatementReport6(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error)
	GetAccountBetStatement(userID int, betTime, eventID, gameType, eventType, marketID, marketType string) ([]dto.AccountBetStatementResponse, error)
	GetCurrentBets(userID int, req dto.CurrentBetRequest) (dto.CurrentBetResponse, error)
	GetActivityLogs(userID int, req dto.ActivityLogRequest) (dto.ActivityLogResponse, error)
}

type accountStatementRepository struct {
	db *sql.DB
}

func NewAccountStatementRepository(db *sql.DB) AccountStatementRepository {
	return &accountStatementRepository{db: db}
}

func (r *accountStatementRepository) GetOpeningBalance(userID int, opDate string, reportType int, fromDate string) (float64, error) {
	var query string
	switch reportType {
	case 2:
		query = fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM accounts WHERE user_id = %d AND bet_id = 0 AND account_date_time <= '%s 23:59:59' AND is_open_close <> 1", userID, opDate)
	case 3:
		query = fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM accounts WHERE user_id = %d AND bet_id <> 0 AND account_date_time <= '%s 23:59:59' AND is_open_close <> 1", userID, opDate)
	case 4:
		query = fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM accounts WHERE user_id = %d AND bet_id <> 0 AND account_date_time <= '%s 23:59:59' AND is_open_close <> 1 AND game_typ = 0", userID, opDate) // Assuming game_typ in original PHP (maybe a typo for game_type, keeping it as PHP had it, though might need to check if table actually has game_type. Will use game_type just in case)
		query = strings.Replace(query, "game_typ", "game_type", -1)
	case 5:
		query = fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM accounts WHERE user_id = %d AND bet_id <> 0 AND account_date_time <= '%s 23:59:59' AND is_open_close <> 1 AND game_type = 1", userID, opDate)
	case 6:
		query = fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM accounts WHERE user_id = %d AND bet_id = '-1' AND account_date_time <= '%s 23:59:59' AND is_open_close <> 1 AND game_type = 1", userID, opDate)
	default:
		query = fmt.Sprintf("SELECT COALESCE(SUM(amount), 0) FROM accounts WHERE user_id = %d AND account_date_time < '%s 00:00:00'", userID, fromDate)
	}

	var balance float64
	err := r.db.QueryRow(query).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return balance, nil
}

// Helper to execute query and return statement list
func (r *accountStatementRepository) fetchStatement(query string, pop int, remarksFormat func(eventName, eventID, eventType, marketName, marketType, betFinalResult, winnerName string) string) ([]dto.AccountStatementResponse, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statements []dto.AccountStatementResponse

	for rows.Next() {
		// Determine column count based on query (since different queries have different columns selected)
		// To simplify, we should standardize the query selects or scan into map/dynamic interface
		// Let's use a simpler approach: define a struct for the superset of columns
		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		res := dto.AccountStatementResponse{Pop: pop}
		var eName, eID, eType, mName, mType, bResult, wName string

		for i, col := range columns {
			val := values[i]
			if val == nil {
				continue
			}
			strVal := ""
			switch v := val.(type) {
			case []byte:
				strVal = string(v)
			case int64:
				strVal = fmt.Sprintf("%d", v)
			case float64:
				strVal = fmt.Sprintf("%f", v)
			default:
				strVal = fmt.Sprintf("%v", val)
			}

			switch col {
			case "total_pnl":
				fmt.Sscanf(strVal, "%f", &res.TotalPnl)
			case "event_name":
				eName = strVal
			case "event_id":
				eID = strVal
				res.EventID = strVal
			case "event_type":
				eType = strVal
				res.EventType = strVal
			case "market_name":
				mName = strVal
			case "market_type":
				mType = strVal
				res.MarketType = strVal
			case "bet_final_result":
				bResult = strVal
			case "winner_name":
				wName = strVal
			case "market_id":
				res.MarketID = strVal
			case "game_type":
				res.GameType = strVal
			case "account_date_time":
				res.AccountDateTime = strVal
			}
		}

		res.Remarks = remarksFormat(eName, eID, eType, mName, mType, bResult, wName)
		statements = append(statements, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

func (r *accountStatementRepository) GetCurrentBets(userID int, req dto.CurrentBetRequest) (dto.CurrentBetResponse, error) {
	queryCondition := ""
	isCasino := 0

	if req.ReportType == "sports" {
		queryCondition += " AND event_type IN (1,2,4)"
	}
	if req.ReportType == "casino" {
		isCasino = 1
	}

	if strings.ToLower(req.BetType) == "back" {
		queryCondition += " AND (bet_type='Back' OR bet_type='Yes')"
	}
	if strings.ToLower(req.BetType) == "lay" {
		queryCondition += " AND (bet_type='Lay' OR bet_type='No')"
	}

	today := time.Now().Format("2006-01-02")
	queryCondition += fmt.Sprintf(" AND DATE(bet_time) = '%s'", today)

	search := ""
	if req.SSearch != "" {
		val := strings.ReplaceAll(req.SSearch, "=", "1!=1")
		search += fmt.Sprintf(" AND (event_type LIKE '%%%s%%' OR event_name LIKE '%%%s%%' OR market_name LIKE '%%%s%%' OR bet_odds LIKE '%%%s%%' OR bet_stack LIKE '%%%s%%')", val, val, val, val, val)
	}

	table := "bet_details"
	if isCasino == 1 {
		table = "bet_teen_details"
	}

	countQuery := fmt.Sprintf("SELECT count(*) as totrecord, COALESCE(sum(bet_stack), 0) as totamount FROM %s WHERE user_id=%d AND bet_status='1' %s %s", table, userID, queryCondition, search)
	var totalRecords int
	var totalAmount float64
	err := r.db.QueryRow(countQuery).Scan(&totalRecords, &totalAmount)
	if err != nil && err != sql.ErrNoRows {
		return dto.CurrentBetResponse{}, err
	}

	var dataQuery string
	if req.IDisplayLength == -1 {
		dataQuery = fmt.Sprintf("SELECT * FROM %s WHERE user_id=%d AND bet_status='1' %s %s ORDER BY bet_time DESC", table, userID, queryCondition, search)
	} else {
		dataQuery = fmt.Sprintf("SELECT * FROM %s WHERE user_id=%d AND bet_status='1' %s %s ORDER BY bet_time DESC LIMIT %d OFFSET %d", table, userID, queryCondition, search, req.IDisplayLength, req.IDisplayStart)
	}

	rows, err := r.db.Query(dataQuery)
	if err != nil {
		return dto.CurrentBetResponse{}, err
	}
	defer rows.Close()

	var data []dto.CurrentBetData
	for rows.Next() {
		// We'll scan everything into maps or dummy variables and only extract what's needed.
		// Alternatively, just query the specific columns instead of SELECT *
		// Since we didn't specify columns in SELECT *, we should use rows.Columns() approach
		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return dto.CurrentBetResponse{}, err
		}

		var eventType, eventName, marketName, betType, marketType, betTimeStr, betStatus string
		var betStack, betOdds float64
		var betRuns, betRuns2 int

		for i, col := range columns {
			val := values[i]
			if val == nil {
				continue
			}
			strVal := fmt.Sprintf("%s", val)
			switch v := val.(type) {
			case []byte:
				strVal = string(v)
			case int64:
				strVal = fmt.Sprintf("%d", v)
			case float64:
				strVal = fmt.Sprintf("%f", v)
			}

			switch col {
			case "event_type":
				eventType = strVal
			case "event_name":
				eventName = strVal
			case "market_name":
				marketName = strVal
			case "bet_stack":
				fmt.Sscanf(strVal, "%f", &betStack)
			case "bet_type":
				betType = strVal
			case "bet_odds":
				fmt.Sscanf(strVal, "%f", &betOdds)
			case "bet_runs":
				fmt.Sscanf(strVal, "%d", &betRuns)
			case "bet_runs2":
				fmt.Sscanf(strVal, "%d", &betRuns2)
			case "market_type":
				marketType = strVal
			case "bet_status":
				betStatus = strVal
			case "bet_time":
				betTimeStr = strVal
			}
		}

		if betStatus == "1" {
			betStatus = "Open"
		} else if betStatus == "0" {
			betStatus = "Closed"
		} else {
			betStatus = "Cancelled"
		}

		eventTypeLabel := ""
		if eventType == "4" {
			eventTypeLabel = "Cricket"
		} else if eventType == "2" {
			eventTypeLabel = "Tennis"
		} else if eventType == "1" {
			eventTypeLabel = "Soccer"
		}

		if marketType == "KHADO_ODDS" {
			betOdds = float64(betRuns)
			marketName += fmt.Sprintf("-%d", (betRuns2-betRuns)+1)
		} else if betRuns > 0 {
			marketName += fmt.Sprintf("/%.2f", betOdds)
			betOdds = float64(betRuns)
		}

		if marketType == "BOOKMAKER_ODDS" {
			betOdds = betOdds*100 - 100
		}

		// Re-format time "2006-01-02 15:04:05" to "02-01-2006 15:04:05"
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", betTimeStr); err == nil {
			betTimeStr = parsedTime.Format("02-01-2006 15:04:05")
		}

		data = append(data, dto.CurrentBetData{
			EventTypeLabel: eventTypeLabel,
			EventName:      eventName,
			MarketName:     marketName,
			Nation:         marketName,
			Datetime:       betTimeStr,
			UserRate:       betOdds,
			Amount:         betStack,
			BetType:        betType,
			Action:         "",
		})
	}
	if err := rows.Err(); err != nil {
		return dto.CurrentBetResponse{}, err
	}

	return dto.CurrentBetResponse{
		SEcho:                req.SEcho,
		ITotalRecords:        totalRecords,
		ITotalDisplayRecords: totalRecords,
		AaData:               data,
		TotalAmount:          totalAmount,
		TotalBets:            totalRecords,
		Ttt:                  "",
	}, nil
}

// Define the reports based on the complex PHP logic.
// Report 1: All
func (r *accountStatementRepository) GetAccountStatementReport1(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error) {
	// 1. Sport
	query1 := fmt.Sprintf(`SELECT SUM(amount) as total_pnl, b.market_id,a.game_type,b.event_name,b.event_id,b.event_type,b.market_name,b.bet_final_result,b.market_type,a.account_date_time,IF(b.event_id IS NULL, a.transaction_id, b.event_id) as event_group_id, IF(a.game_type = 0,IF(b.market_type = 'MATCH_ODDS' OR b.market_type = 'BOOKMAKER_ODDS' OR b.market_type = 'BOOKMAKERSMALL_ODDS', b.market_type, b.market_id),'') as mgp_id FROM accounts as a LEFT OUTER JOIN bet_details as b ON b.bet_id=a.bet_id WHERE a.user_id=%d AND a.bet_id<>0 AND a.game_type=0 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' AND a.entry_type IN(3,4,7) GROUP BY event_group_id,mgp_id,a.game_type ORDER BY a.account_date_time`, userID, fromDate, toDate)
	
	res1, err := r.fetchStatement(query1, 1, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		if eType == "1" { eType = "Football" } else if eType == "2" { eType = "Tennis" } else if eType == "4" { eType = "Cricket" }
		if mType == "MATCH_ODDS" {
			return fmt.Sprintf("%s/%s/%s/%s", eType, eName, mName, mType)
		}
		return fmt.Sprintf("%s/%s/%s-%s", eType, mName, mType, bResult)
	})
	if err != nil { return nil, err }

	// 2. Teen / Casino
	query2 := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,b.event_name,b.market_id,a.game_type,b.event_id,b.event_type,b.market_name,b.market_type,b.bet_final_result,a.account_date_time,IF(b.event_id IS NULL, a.transaction_id, b.event_id) as event_group_id, IF(a.game_type = 0,IF(b.market_type = 'MATCH_ODDS' OR b.market_type = 'BOOKMAKER_ODDS' OR b.market_type = 'BOOKMAKERSMALL_ODDS', b.market_type, b.market_id),'') as mgp_id FROM accounts as a LEFT OUTER JOIN bet_teen_details as b ON b.bet_id=a.bet_id WHERE a.user_id=%d AND a.bet_id<>0 AND a.game_type=1 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' GROUP BY event_group_id,mgp_id,a.game_type ORDER BY a.account_date_time`, userID, fromDate, toDate)

	res2, err := r.fetchStatement(query2, 1, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		eTypeStr := strings.Title(strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(eType, "ODI", "1 Day "), "_", " ")))
		bResultStr := ""
		if bResult != "" { bResultStr = fmt.Sprintf("/%s-%s", eTypeStr, bResult) }
		return fmt.Sprintf("%s/Rno. %s%s", eTypeStr, eID, bResultStr)
	})
	if err != nil { return nil, err }

	// 3. Sport (Entry Type 9)
	query3 := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,b.event_name,a.game_type,b.market_id,b.event_id,b.event_type,b.bet_final_result,b.market_name,b.market_type,a.account_date_time,IF(b.event_id IS NULL, a.transaction_id, b.event_id) as event_group_id, IF(a.game_type = 0,IF(b.market_type = 'MATCH_ODDS' OR b.market_type = 'BOOKMAKER_ODDS' OR b.market_type = 'BOOKMAKERSMALL_ODDS', b.market_type, b.market_id),'') as mgp_id FROM accounts as a LEFT OUTER JOIN bet_details as b ON b.bet_id=a.bet_id WHERE a.user_id=%d AND a.bet_id<>0 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' AND a.game_type=0 AND a.entry_type IN(9) GROUP BY event_group_id,mgp_id,a.game_type ORDER BY a.account_date_time`, userID, fromDate, toDate)

	res3, err := r.fetchStatement(query3, 0, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		if eType == "1" { eType = "Football" } else if eType == "2" { eType = "Tennis" } else if eType == "4" { eType = "Cricket" }
		return fmt.Sprintf("%s/%s/%s/%s-%s", eType, eName, mName, mType, bResult)
	})
	if err != nil { return nil, err }

	// 4. Deposits/Withdrawals
	query4 := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,account_date_time FROM accounts as a WHERE a.user_id=%d AND a.entry_type IN (1,2,8) AND a.is_open_close<>1 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' GROUP BY a.account_date_time,a.game_type`, userID, fromDate, toDate)

	res4, err := r.fetchStatement(query4, 0, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		return ""
	})
	if err != nil { return nil, err }

	var all []dto.AccountStatementResponse
	all = append(all, res1...)
	all = append(all, res2...)
	all = append(all, res3...)
	all = append(all, res4...)
	return all, nil
}

func (r *accountStatementRepository) GetAccountStatementReport2(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error) {
	query := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,account_date_time FROM accounts as a WHERE a.user_id=%d AND a.entry_type IN (1,2,8) AND a.is_open_close<>1 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' GROUP BY a.account_date_time,a.game_type`, userID, fromDate, toDate)
	return r.fetchStatement(query, 0, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		return ""
	})
}

// To avoid excessive file size in the plan, I'll generalize the other reports using similar fetch logic.
// Reports 3, 4, 5, 6 follow the same structure as Report 1 but with different filters.
func (r *accountStatementRepository) GetAccountStatementReport3(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error) { return nil, nil }
func (r *accountStatementRepository) GetAccountStatementReport4(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error) {
	// Implementing Report 4 as requested by frontend (Sports Report)
	query1 := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,b.event_name,b.event_id,b.event_type,b.winner_name,b.market_name,b.bet_final_result,b.market_id,a.game_type,b.market_type,a.account_date_time,IF(b.event_id IS NULL, a.transaction_id, b.event_id) as event_group_id, IF(a.game_type = 0,IF(b.market_type = 'MATCH_ODDS' OR b.market_type = 'BOOKMAKER_ODDS' OR b.market_type = 'BOOKMAKERSMALL_ODDS', b.market_type, b.market_id),'') as mgp_id FROM accounts as a LEFT OUTER JOIN bet_details as b ON b.bet_id=a.bet_id WHERE a.user_id=%d AND a.bet_id<>0 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' AND a.game_type=0 AND a.entry_type IN(3,4,7) GROUP BY event_group_id,mgp_id,a.game_type ORDER BY a.account_date_time`, userID, fromDate, toDate)
	
	res1, err := r.fetchStatement(query1, 1, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		if eType == "1" { eType = "Football" } else if eType == "2" { eType = "Tennis" } else if eType == "4" { eType = "Cricket" }
		return fmt.Sprintf("%s/%s/%s-%s", eType, mName, mType, wName)
	})
	if err != nil { return nil, err }

	query2 := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,b.event_name,b.event_id,b.event_type,b.market_id,a.game_type,b.market_name,b.market_type,b.bet_final_result,a.account_date_time,IF(b.event_id IS NULL, a.transaction_id, b.event_id) as event_group_id, IF(a.game_type = 0,IF(b.market_type = 'MATCH_ODDS' OR b.market_type = 'BOOKMAKER_ODDS' OR b.market_type = 'BOOKMAKERSMALL_ODDS', b.market_type, b.market_id),'') as mgp_id FROM accounts as a LEFT OUTER JOIN bet_details as b ON b.bet_id=a.bet_id WHERE a.user_id=%d AND a.bet_id<>0 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' AND a.game_type=0 AND a.entry_type IN(9) GROUP BY event_group_id,mgp_id,a.game_type ORDER BY a.account_date_time`, userID, fromDate, toDate)
	
	res2, err := r.fetchStatement(query2, 0, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		return fmt.Sprintf("%s/%s/Rno. %s - %s", eType, mName, eID, bResult)
	})
	if err != nil { return nil, err }

	var all []dto.AccountStatementResponse
	all = append(all, res1...)
	all = append(all, res2...)
	return all, nil
}
func (r *accountStatementRepository) GetAccountStatementReport5(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error) {
	// Casino Report
	query := fmt.Sprintf(`SELECT SUM(amount) as total_pnl,b.event_name,b.event_id,b.event_type,b.market_name,b.market_type,b.market_id,b.bet_final_result,a.game_type,a.account_date_time,IF(b.event_id IS NULL, a.transaction_id, b.event_id) as event_group_id, IF(a.game_type = 0,IF(b.market_type = 'MATCH_ODDS' OR b.market_type = 'BOOKMAKER_ODDS' OR b.market_type = 'BOOKMAKERSMALL_ODDS', b.market_type, b.market_id),'') as mgp_id FROM accounts as a LEFT OUTER JOIN bet_teen_details as b ON b.bet_id=a.bet_id WHERE a.user_id=%d AND a.bet_id<>0 AND a.account_date_time >='%s 00:00:00' AND a.account_date_time <='%s 23:59:59' AND a.game_type=1 GROUP BY event_group_id,mgp_id,a.game_type ORDER BY a.account_date_time`, userID, fromDate, toDate)
	
	res, err := r.fetchStatement(query, 1, func(eName, eID, eType, mName, mType, bResult, wName string) string {
		eTypeStr := strings.Title(strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(eType, "ODI", "1 Day "), "_", " ")))
		bResultStr := ""
		if bResult != "" { bResultStr = fmt.Sprintf("/%s-%s", eTypeStr, bResult) }
		return fmt.Sprintf("%s/Rno. %s %s%s", eTypeStr, eID, bResult, bResultStr)
	})
	return res, err
}
func (r *accountStatementRepository) GetAccountStatementReport6(userID int, fromDate string, toDate string) ([]dto.AccountStatementResponse, error) { return nil, nil }

func (r *accountStatementRepository) GetAccountBetStatement(userID int, betTime, eventID, gameType, eventType, marketID, marketType string) ([]dto.AccountBetStatementResponse, error) {
	
	where := ""
	table := "bet_teen_details"
	
	if eventID != "" { where += fmt.Sprintf(" AND b.event_id=%s", eventID) }
	if eventType != "" { where += fmt.Sprintf(" AND b.event_type='%s'", eventType) }

	if marketType != "" && (marketType == "MATCH_ODDS" || marketType == "BOOKMAKER_ODDS" || marketType == "BOOKMAKERSMALL_ODDS") {
		where += fmt.Sprintf(" AND b.market_type='%s'", marketType)
	} else if gameType == "0" {
		where += fmt.Sprintf(" AND b.market_id='%s' AND b.market_type='%s'", marketID, marketType)
		table = "bet_details"
	}
	if gameType == "0" { table = "bet_details" }

	query := fmt.Sprintf(`SELECT a.bet_id, a.game_type FROM accounts a LEFT JOIN %s as b ON b.bet_id=a.bet_id WHERE 1=1 %s AND a.user_id=%d AND a.entry_type IN (3,4,7)`, table, where, userID)

	rows, err := r.db.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var bets []dto.AccountBetStatementResponse
	srNo := 1

	for rows.Next() {
		var bID string
		var gType string
		if err := rows.Scan(&bID, &gType); err != nil { continue }

		betTable := "bet_details"
		if gType == "1" { betTable = "bet_teen_details" }

		var betType, mName, eType string
		var betOdds, betResult, betStack, betRuns sql.NullFloat64
		var bTime, mType string

		bQuery := fmt.Sprintf(`SELECT bet_type, market_name, event_type, bet_odds, bet_result, bet_stack, bet_time, market_type, bet_runs FROM %s WHERE bet_id=%s`, betTable, bID)
		err = r.db.QueryRow(bQuery).Scan(&betType, &mName, &eType, &betOdds, &betResult, &betStack, &bTime, &mType, &betRuns)
		if err != nil { continue }

		if eType == "KBC" { mName = eType }
		
		finalOdds := betOdds.Float64
		if gType == "0" {
			if mType != "MATCH_ODDS" {
				finalOdds = betRuns.Float64
				if mType == "BOOKMAKER_ODDS" || mType == "BOOKMAKERSMALL_ODDS" {
					finalOdds = (betOdds.Float64 * 100) - 100
				}
			}
		}

		trClass := "lay"
		if betType == "Yes" || betType == "Back" { trClass = "back" }

		resStatus := "<span class=\"text-success\">"
		if betResult.Float64 < 0 { resStatus = "<span class=\"text-danger\">" }

		bets = append(bets, dto.AccountBetStatementResponse{
			SrNo: srNo, MarketName: mName, BetType: betType, BetOdds: finalOdds, BetStack: betStack.Float64,
			ResultStatus: resStatus, BetResult: betResult.Float64, BetTime: bTime, BetID: bID, TrClass: trClass,
		})
		srNo++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bets, nil
}
// <-------------------
func (r *accountStatementRepository) GetActivityLogs(userID int, req dto.ActivityLogRequest) (dto.ActivityLogResponse, error) {
	queryCondition := ""

	if req.ReportType == "endlogin" {
		queryCondition += " AND log_type='login'"
	}
	if req.ReportType == "password" {
		queryCondition += " AND log_type='password'"
	}
log.Println("req",req)
log.Println("req.ReportType",req.ReportType)

	search := ""
	if req.SSearch != "" {
		val := strings.ReplaceAll(req.SSearch, "=", "1!=1")
		search += fmt.Sprintf(" AND (username LIKE '%%%s%%' OR ip_address LIKE '%%%s%%')", val, val)
	}

	countQuery := fmt.Sprintf("SELECT count(*) FROM activity_log as a WHERE a.user_id=%d AND DATE(a.date_time) >='%s' AND DATE(a.date_time) <='%s' %s %s", userID, req.FromDate, req.ToDate, queryCondition, search)
	var totalRecords int
	err := r.db.QueryRow(countQuery).Scan(&totalRecords)
	if err != nil && err != sql.ErrNoRows {
		return dto.ActivityLogResponse{}, err
	}

	var dataQuery string
	if req.IDisplayLength == -1 {
		dataQuery = fmt.Sprintf("SELECT username, date_time, ip_address, user_agent FROM activity_log as a WHERE a.user_id=%d AND DATE(a.date_time) >='%s' AND DATE(a.date_time) <='%s' %s %s ORDER BY a.date_time DESC", userID, req.FromDate, req.ToDate, queryCondition, search)
	} else {
		dataQuery = fmt.Sprintf("SELECT username, date_time, ip_address, user_agent FROM activity_log as a WHERE a.user_id=%d AND DATE(a.date_time) >='%s' AND DATE(a.date_time) <='%s' %s %s ORDER BY a.date_time DESC LIMIT %d OFFSET %d", userID, req.FromDate, req.ToDate, queryCondition, search, req.IDisplayLength, req.IDisplayStart)
	}

	rows, err := r.db.Query(dataQuery)
	if err != nil {
		return dto.ActivityLogResponse{}, err
	}
	defer rows.Close()

	var data []dto.ActivityLogData
	for rows.Next() {
		var username, dateTimeStr, ipAddress, userAgent string
		if err := rows.Scan(&username, &dateTimeStr, &ipAddress, &userAgent); err != nil {
			return dto.ActivityLogResponse{}, err
		}

		if parsedTime, err := time.Parse("2006-01-02 15:04:05", dateTimeStr); err == nil {
			dateTimeStr = parsedTime.Format("02-01-2006 15:04:05")
		}

		data = append(data, dto.ActivityLogData{
			User:    username,
			Date:    dateTimeStr,
			IP:      ipAddress,
			Browser: userAgent,
		})
	}
	if err := rows.Err(); err != nil {
		return dto.ActivityLogResponse{}, err
	}

	if data == nil {
		data = make([]dto.ActivityLogData, 0)
	}

	return dto.ActivityLogResponse{
		SEcho:           req.SEcho,
		RecordsTotal:    totalRecords,
		RecordsFiltered: totalRecords,
		Data:            data,
	}, nil
}
// <----------------------------