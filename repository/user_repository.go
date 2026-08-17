package repository

import (
	"database/sql"

	"github.com/kirtanwyn/allPanelexch/model"
)

// UserRepository defines the contract for user-related queries.
type UserRepository interface {
	GetUserStatus(userID int) (*model.UserMasterParentStatus, error)
	GetUserFull(userID int) (*model.UserMaster, error)
	GetParentStatuses(parentIDs []int) (map[int]*model.UserMasterParentStatus, error)
	GetAccountBalance(userID int) (float64, error)
	GetSportAccess(userID int) (soccerAccess, tennisAccess int, err error)
	GetUserStakeLimits(userID int) (*model.UserMaster, error)
	GetUserExposureLimits(userID int) (*model.UserMaster, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetUserStatus fetches the core status fields for a user.
func (r *userRepository) GetUserStatus(userID int) (*model.UserMasterParentStatus, error) {
	var status model.UserMasterParentStatus
	query := `SELECT Status, bet_status, fancy_bet_status FROM user_master WHERE Id = ?`
	err := r.db.QueryRow(query, userID).Scan(&status.Status, &status.BetStatus, &status.FancyBetStatus)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetUserFull fetches all user_master fields needed for bet placement.
func (r *userRepository) GetUserFull(userID int) (*model.UserMaster, error) {
	var u model.UserMaster
	query := `SELECT 
		Id, Name, Status, bet_status, fancy_bet_status,
		parentDL, parentMDL, parentSuperMDL, parentKingAdmin,
		soccer_access, tennis_access,
		min_stake, max_stake,
		min_cricket_stake, max_cricket_stake,
		min_soccer_stake, max_soccer_stake,
		min_tennis_stake, max_tennis_stake,
		min_fancy_stake, max_fancy_stake,
		net_exposure_limit,
		minimum_odds, maximum_odds,
		bet_email_notify
	FROM user_master WHERE Id = ?`
	err := r.db.QueryRow(query, userID).Scan(
		&u.ID, &u.Name, &u.Status, &u.BetStatus, &u.FancyBetStatus,
		&u.ParentDL, &u.ParentMDL, &u.ParentSuperMDL, &u.ParentKingAdmin,
		&u.SoccerAccess, &u.TennisAccess,
		&u.MinStake, &u.MaxStake,
		&u.MinCricketStake, &u.MaxCricketStake,
		&u.MinSoccerStake, &u.MaxSoccerStake,
		&u.MinTennisStake, &u.MaxTennisStake,
		&u.MinFancyStake, &u.MaxFancyStake,
		&u.NetExposureLimit,
		&u.MinimumOdds, &u.MaximumOdds,
		&u.BetEmailNotify,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetParentStatuses fetches Status/BetStatus/FancyBetStatus for all parent IDs.
func (r *userRepository) GetParentStatuses(parentIDs []int) (map[int]*model.UserMasterParentStatus, error) {
	result := make(map[int]*model.UserMasterParentStatus)
	for _, pid := range parentIDs {
		if pid <= 0 {
			continue
		}
		var s model.UserMasterParentStatus
		query := `SELECT Status, bet_status, fancy_bet_status FROM user_master WHERE Id = ?`
		err := r.db.QueryRow(query, pid).Scan(&s.Status, &s.BetStatus, &s.FancyBetStatus)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, err
		}
		result[pid] = &s
	}
	return result, nil
}

// GetAccountBalance returns the sum of accounts where status=1.
func (r *userRepository) GetAccountBalance(userID int) (float64, error) {
	var totalBalance sql.NullFloat64
	query := `SELECT SUM(amount) as total_balance FROM accounts WHERE user_id = ? AND status = 1`
	err := r.db.QueryRow(query, userID).Scan(&totalBalance)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return totalBalance.Float64, nil
}

// GetSportAccess returns soccer_access and tennis_access flags.
func (r *userRepository) GetSportAccess(userID int) (soccerAccess, tennisAccess int, err error) {
	query := `SELECT soccer_access, tennis_access FROM user_master WHERE Id = ?`
	err = r.db.QueryRow(query, userID).Scan(&soccerAccess, &tennisAccess)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	return soccerAccess, tennisAccess, err
}

// GetUserStakeLimits returns stake limit fields for a user.
func (r *userRepository) GetUserStakeLimits(userID int) (*model.UserMaster, error) {
	var u model.UserMaster
	query := `SELECT 
		min_stake, max_stake,
		min_cricket_stake, max_cricket_stake,
		min_soccer_stake, max_soccer_stake,
		min_tennis_stake, max_tennis_stake,
		min_fancy_stake, max_fancy_stake
	FROM user_master WHERE Id = ?`
	err := r.db.QueryRow(query, userID).Scan(
		&u.MinStake, &u.MaxStake,
		&u.MinCricketStake, &u.MaxCricketStake,
		&u.MinSoccerStake, &u.MaxSoccerStake,
		&u.MinTennisStake, &u.MaxTennisStake,
		&u.MinFancyStake, &u.MaxFancyStake,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserExposureLimits returns net_exposure_limit and odds limits.
func (r *userRepository) GetUserExposureLimits(userID int) (*model.UserMaster, error) {
	var u model.UserMaster
	query := `SELECT net_exposure_limit, minimum_odds, maximum_odds FROM user_master WHERE Id = ?`
	err := r.db.QueryRow(query, userID).Scan(&u.NetExposureLimit, &u.MinimumOdds, &u.MaximumOdds)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}