package service

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/kirtanwyn/allPanelexch/dto"
	"github.com/kirtanwyn/allPanelexch/repository"
)

type AccountStatementService interface {
	GetAccountStatement(userID int, req dto.AccountStatementRequest) ([]dto.AccountStatementResponse, error)
	GetAccountBetStatement(userID int, req dto.AccountBetStatementRequest) ([]dto.AccountBetStatementResponse, error)
	GetCurrentBets(userID int, req dto.CurrentBetRequest) (dto.CurrentBetResponse, error)
	GetActivityLogs(userID int, req dto.ActivityLogRequest) (dto.ActivityLogResponse, error)
	UpdateButtonValue(userID int, req dto.UpdateButtonValueRequest) (dto.UpdateButtonValueResponse, error)
	RefreshBalance(userID int) (dto.RefreshBalanceResponse, error)
}

type accountStatementService struct {
	repo repository.AccountStatementRepository
}

func NewAccountStatementService(repo repository.AccountStatementRepository) AccountStatementService {
	return &accountStatementService{repo: repo}
}

func (s *accountStatementService) GetAccountStatement(userID int, req dto.AccountStatementRequest) ([]dto.AccountStatementResponse, error) {
	minDate := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	if req.FromDate == "" {
		return nil, fmt.Errorf("Please select start date.")
	}
	if req.ToDate == "" {
		return nil, fmt.Errorf("Please select end date.")
	}
	if minDate >= req.FromDate {
		return nil, fmt.Errorf("You can not view before 90 days profit/loss report.")
	}

	fromDateParsed, _ := time.Parse("2006-01-02", req.FromDate)
	opDate := fromDateParsed.AddDate(0, 0, -1).Format("2006-01-02")

	// 1. Get Opening Balance
	openingBalance, err := s.repo.GetOpeningBalance(userID, opDate, req.ReportType, req.FromDate)
	if err != nil {
		return nil, err
	}

	openingRecord := dto.AccountStatementResponse{
		AccountDateTime: fmt.Sprintf("%s 00:00:00", req.FromDate), // Format later if needed to d-m-Y H:i:s
		TotalPnl:        openingBalance,
		Remarks:         "Opening Balance",
		Pop:             0,
	}
	
	// Ensure format matches PHP: d-m-Y H:i:s
	if t, err := time.Parse("2006-01-02 15:04:05", openingRecord.AccountDateTime); err == nil {
		openingRecord.AccountDateTime = t.Format("02-01-2006 15:04:05")
	}

	var statements []dto.AccountStatementResponse
	statements = append(statements, openingRecord)

	// 2. Fetch specific report data
	var reportData []dto.AccountStatementResponse
	switch req.ReportType {
	case 1:
		reportData, err = s.repo.GetAccountStatementReport1(userID, req.FromDate, req.ToDate)
	case 2:
		reportData, err = s.repo.GetAccountStatementReport2(userID, req.FromDate, req.ToDate)
	case 3:
		reportData, err = s.repo.GetAccountStatementReport3(userID, req.FromDate, req.ToDate)
	case 4:
		reportData, err = s.repo.GetAccountStatementReport4(userID, req.FromDate, req.ToDate)
	case 5:
		reportData, err = s.repo.GetAccountStatementReport5(userID, req.FromDate, req.ToDate)
	case 6:
		reportData, err = s.repo.GetAccountStatementReport6(userID, req.FromDate, req.ToDate)
	default:
		reportData, err = s.repo.GetAccountStatementReport1(userID, req.FromDate, req.ToDate)
	}

	if err != nil {
		return nil, err
	}

	// Format dates in report data
	for i := range reportData {
		if t, err := time.Parse("2006-01-02 15:04:05", reportData[i].AccountDateTime); err == nil {
			reportData[i].AccountDateTime = t.Format("02-01-2006 15:04:05")
		}
	}

	statements = append(statements, reportData...)

	// 3. Sort by date
	sort.SliceStable(statements, func(i, j int) bool {
		ti, _ := time.Parse("02-01-2006 15:04:05", statements[i].AccountDateTime)
		tj, _ := time.Parse("02-01-2006 15:04:05", statements[j].AccountDateTime)
		return ti.Before(tj)
	})

	return statements, nil
}

func (s *accountStatementService) GetAccountBetStatement(userID int, req dto.AccountBetStatementRequest) ([]dto.AccountBetStatementResponse, error) {
	
	// Convert betTime to Y-m-d H:i:s
	t, err := time.Parse(time.RFC3339, req.BetTime) // Or whatever format it comes in
	if err != nil {
		t, err = time.Parse("02-01-2006 15:04:05", req.BetTime)
		if err != nil {
			// fallback
		}
	}
	betTime := t.Format("2006-01-02 15:04:05")
	if req.BetTime == "" {
		betTime = req.BetTime
	}

	bets, err := s.repo.GetAccountBetStatement(userID, betTime, req.EventID, req.GameType, req.EventType, req.MarketID, req.MarketType)
	if err != nil {
		return nil, err
	}
	
	// Format time
	for i := range bets {
		if t, err := time.Parse("2006-01-02 15:04:05", bets[i].BetTime); err == nil {
			bets[i].BetTime = t.Format("02-01-2006 15:04:05")
		}
	}

	return bets, nil
}

func (s *accountStatementService) GetCurrentBets(userID int, req dto.CurrentBetRequest) (dto.CurrentBetResponse, error) {
	return s.repo.GetCurrentBets(userID, req)
}
// <----------------
func (s *accountStatementService) GetActivityLogs(userID int, req dto.ActivityLogRequest) (dto.ActivityLogResponse, error) {
	return s.repo.GetActivityLogs(userID, req)
}

func (s *accountStatementService) UpdateButtonValue(userID int, req dto.UpdateButtonValueRequest) (dto.UpdateButtonValueResponse, error) {
	err := s.repo.UpdateButtonValue(userID, req.AllButtonValue, req.Type)
	if err != nil {
		return dto.UpdateButtonValueResponse{
			Status:  "ok",
			Message: "Something went wrong,please try again later",
		}, nil // returning nil error to match original PHP behavior (which always returns ok with error message)
	}

	return dto.UpdateButtonValueResponse{
		Status:  "ok",
		Message: "Button value Changed",
	}, nil
}

func formatFloat(val float64) string {
	if math.Floor(val) == val {
		return fmt.Sprintf("%.0f", val)
	}
	return fmt.Sprintf("%.2f", val)
}

func (s *accountStatementService) RefreshBalance(userID int) (dto.RefreshBalanceResponse, error) {
	accountBalance, err := s.repo.GetAccountBalance(userID)
	if err != nil {
		return dto.RefreshBalanceResponse{}, err
	}

	unmatchedExposureBalance, err := s.repo.GetUnmatchedExposure(userID)
	if err != nil {
		return dto.RefreshBalanceResponse{}, err
	}

	exposureBalance, err := s.repo.GetTotalNetExposure(userID)
	if err != nil {
		return dto.RefreshBalanceResponse{}, err
	}

	exposureWinningBalance, err := s.repo.GetTotalOnlyWinning(userID)
	if err != nil {
		return dto.RefreshBalanceResponse{}, err
	}

	checkPlusExpo := exposureBalance + unmatchedExposureBalance
	if checkPlusExpo > 0 {
		_ = accountBalance
	} else {
		_ = accountBalance + exposureBalance + unmatchedExposureBalance
	}

	exposureBalance = exposureBalance * (-1)
	unmatchedExposureBalance = unmatchedExposureBalance * (-1)

	exposure := exposureBalance + unmatchedExposureBalance

	return dto.RefreshBalanceResponse{
		Status:   "ok",
		Balance:  formatFloat(accountBalance),
		Exposure: formatFloat(exposure),
		Winning:  fmt.Sprintf("%.2f", exposureWinningBalance),
	}, nil
}
