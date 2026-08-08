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
	Status                 string `json:"status" db:"Status"` // In PHP compared as "1"
	UserVerificationStatus string `json:"user_verification_status" db:"user_verification_status"`
	UserVerificationType   *string `json:"user_verification_type" db:"user_verification_type"`
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
