package model

type SiteUnderMaintenance struct {
	SiteStatus int `json:"site_status" db:"site_status"`
}

type UserLoginMaster struct {
	UserID              int    `json:"user_id" db:"UserID"`
	EmailID             string `json:"email_id" db:"Email_ID"`
	UserType            int    `json:"user_type" db:"UserType"`
	Password            string  `json:"-" db:"Password"`
	UserPasswordSalt    *string `json:"-" db:"user_password_salt"`
	UserPasswordSaltKey *string `json:"-" db:"user_password_salt_key"`
	SecretKey           *string `json:"-" db:"SecretKey"`
	FirstPasswordChanged int    `json:"first_password_changed" db:"first_password_changed"`
	ParentDL            int    `json:"parent_dl" db:"parentDL"`
	ParentMDL           int    `json:"parent_mdl" db:"parentMDL"`
	ParentSuperMDL      int    `json:"parent_super_mdl" db:"parentSuperMDL"`
	ParentKingAdmin     int    `json:"parent_king_admin" db:"parentKingAdmin"`
	LoginString         string `json:"login_string" db:"loginString"`
	ApiAuthToken        string `json:"api_auth_token" db:"api_auth_token"`
}

type UserMaster struct {
	ID                     int    `json:"id" db:"Id"`
	Name                   string `json:"name" db:"Name"`
	Status                 int    `json:"status" db:"Status"`
	BetStatus              int    `json:"bet_status" db:"bet_status"`
	FancyBetStatus         int    `json:"fancy_bet_status" db:"fancy_bet_status"`
	ParentDL               int    `json:"parent_dl" db:"parentDL"`
	ParentMDL              int    `json:"parent_mdl" db:"parentMDL"`
	ParentSuperMDL         int    `json:"parent_super_mdl" db:"parentSuperMDL"`
	ParentKingAdmin        int    `json:"parent_king_admin" db:"parentKingAdmin"`
	SoccerAccess           int    `json:"soccer_access" db:"soccer_access"`
	TennisAccess           int    `json:"tennis_access" db:"tennis_access"`
	MinStake               int    `json:"min_stake" db:"min_stake"`
	MaxStake               int    `json:"max_stake" db:"max_stake"`
	MinCricketStake        int    `json:"min_cricket_stake" db:"min_cricket_stake"`
	MaxCricketStake        int    `json:"max_cricket_stake" db:"max_cricket_stake"`
	MinSoccerStake         int    `json:"min_soccer_stake" db:"min_soccer_stake"`
	MaxSoccerStake         int    `json:"max_soccer_stake" db:"max_soccer_stake"`
	MinTennisStake         int    `json:"min_tennis_stake" db:"min_tennis_stake"`
	MaxTennisStake         int    `json:"max_tennis_stake" db:"max_tennis_stake"`
	MinFancyStake          int    `json:"min_fancy_stake" db:"min_fancy_stake"`
	MaxFancyStake          int    `json:"max_fancy_stake" db:"max_fancy_stake"`
	NetExposureLimit       int    `json:"net_exposure_limit" db:"net_exposure_limit"`
	MinimumOdds            int    `json:"minimum_odds" db:"minimum_odds"`
	MaximumOdds            int    `json:"maximum_odds" db:"maximum_odds"`
	BetEmailNotify         int    `json:"bet_email_notify" db:"bet_email_notify"`
	UserVerificationStatus string `json:"user_verification_status" db:"user_verification_status"`
	UserVerificationType   *string `json:"user_verification_type" db:"user_verification_type"`
	IsUserVerified         string `json:"is_user_verified" db:"is_user_verified"`
}

// UserMasterParentStatus holds the status fields when querying parent hierarchy
type UserMasterParentStatus struct {
	Status       int `json:"status" db:"Status"`
	BetStatus    int `json:"bet_status" db:"bet_status"`
	FancyBetStatus int `json:"fancy_bet_status" db:"fancy_bet_status"`
}

type LoginIPAddress struct {
	UserID        int    `json:"user_id" db:"user_id"`
	IPAddress     string `json:"ip_address" db:"ip_address"`
	LoginDateTime string `json:"login_date_time" db:"login_date_time"`
	UserAgent     string `json:"user_agent" db:"user_agent"`
}

type ActivityLog struct {
	UserID    int    `json:"user_id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	IPAddress string `json:"ip_address" db:"ip_address"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	DateTime  string `json:"date_time" db:"date_time"`
	LogType   string `json:"log_type" db:"log_type"`
}
