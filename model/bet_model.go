package model

// BetDetails represents the bet_details table — one row per placed bet.
type BetDetails struct {
	ID                int     `json:"id" db:"id"`
	BetIPAddress      string  `json:"bet_ip_address" db:"bet_ip_address"`
	BetUserAgent      string  `json:"bet_user_agent" db:"bet_user_agent"`
	EventID           int     `json:"event_id" db:"event_id"`
	EventType         int     `json:"event_type" db:"event_type"`
	OddsMarketID      string  `json:"oddsmarketId" db:"oddsmarketId"`
	MarketID          string  `json:"market_id" db:"market_id"`
	UserID            int     `json:"user_id" db:"user_id"`
	EventName         string  `json:"event_name" db:"event_name"`
	MarketName        string  `json:"market_name" db:"market_name"`
	MarketType        string  `json:"market_type" db:"market_type"`
	BetType           string  `json:"bet_type" db:"bet_type"`
	BetRuns           float64 `json:"bet_runs" db:"bet_runs"`
	BetRuns2          float64 `json:"bet_runs2" db:"bet_runs2"`
	BetOdds           float64 `json:"bet_odds" db:"bet_odds"`
	BetStack          float64 `json:"bet_stack" db:"bet_stack"`
	BetResult         int     `json:"bet_result" db:"bet_result"`
	BetMarginUsed     float64 `json:"bet_margin_used" db:"bet_margin_used"`
	BetTime           string  `json:"bet_time" db:"bet_time"`
	BetStatus         int     `json:"bet_status" db:"bet_status"`
	BetWinAmount      float64 `json:"bet_win_amount" db:"bet_win_amount"`
	RunnerID          string  `json:"runner_id" db:"runner_id"`
	RunnerName1       string  `json:"runner_name1" db:"runner_name1"`
	DisplayMarketType string  `json:"display_market_type" db:"display_market_type"`
}

// ExposureDetails represents the exposure_details table — tracks net exposure per user/event/market.
type ExposureDetails struct {
	ID                 int     `json:"id" db:"id"`
	UserID             int     `json:"user_id" db:"user_id"`
	EventID            string  `json:"event_id" db:"event_id"`
	MarketID           string  `json:"market_id" db:"market_id"`
	MarketType         string  `json:"market_type" db:"market_type"`
	ExposureAmount     float64 `json:"exposure_amount" db:"exposure_amount"`
	ExposureDatetime   string  `json:"exposure_datetime" db:"exposure_datetime"`
	EventType          string  `json:"event_type" db:"event_type"`
	MaxWinningAmount   float64 `json:"max_winning_amount" db:"max_winning_amount"`
	MeterMarketID      string  `json:"meter_market_id" db:"meter_market_id"`
	CasinoBackName     string  `json:"casino_back_name" db:"casino_back_name"`
}

// UnmatchedBetDetails represents the unmatched_bet_details table.
type UnmatchedBetDetails struct {
	ID             int     `json:"id" db:"id"`
	UserID         int     `json:"user_id" db:"user_id"`
	EventID        int     `json:"event_id" db:"event_id"`
	MarketID       string  `json:"market_id" db:"market_id"`
	MarketType     string  `json:"market_type" db:"market_type"`
	BetType        string  `json:"bet_type" db:"bet_type"`
	BetMarginUsed  float64 `json:"bet_margin_used" db:"bet_margin_used"`
	BetWinAmount   float64 `json:"bet_win_amount" db:"bet_win_amount"`
	BetStatus      int     `json:"bet_status" db:"bet_status"`
}

// BetSuccessLog represents the bet_success_log table.
type BetSuccessLog struct {
	ID              int    `json:"id" db:"id"`
	UserID          int    `json:"user_id" db:"user_id"`
	LogMarketName   string `json:"log_market_name" db:"log_market_name"`
	APIData         string `json:"api_data" db:"api_data"`
	APIURL          string `json:"api_url" db:"api_url"`
	BetID           int    `json:"bet_id" db:"bet_id"`
	LogDetails      string `json:"log_details" db:"log_details"`
	PageCallTime    string `json:"page_call_time" db:"page_call_time"`
	LogTime         string `json:"log_time" db:"log_time"`
	GameType        string `json:"game_type" db:"game_type"`
}
