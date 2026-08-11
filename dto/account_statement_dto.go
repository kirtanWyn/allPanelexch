package dto

type AccountStatementRequest struct {
	FromDate   string `form:"from_date" json:"from_date" binding:"required"`
	ToDate     string `form:"to_date" json:"to_date" binding:"required"`
	ReportType int    `form:"report_type" json:"report_type" binding:"required"`
}

type AccountStatementResponse struct {
	AccountDateTime string  `json:"account_date_time"`
	TotalPnl        float64 `json:"total_pnl"`
	Remarks         string  `json:"remakrs"` // Kept spelling as in original PHP for frontend compatibility
	Pop             int     `json:"pop"`
	EventID         string  `json:"event_id"`
	GameType        string  `json:"game_type"`
	EventType       string  `json:"event_type"`
	MarketID        string  `json:"market_id"`
	MarketType      string  `json:"market_type"`
}

type AccountBetStatementRequest struct {
	BetTime    string `form:"bet_time" json:"bet_time" binding:"required"`
	EventID    string `form:"event_id" json:"event_id"`
	GameType   string `form:"game_type" json:"game_type" binding:"required"`
	EventType  string `form:"event_type" json:"event_type"`
	MarketID   string `form:"market_id" json:"market_id"`
	MarketType string `form:"market_type" json:"market_type"`
}

type AccountBetStatementResponse struct {
	SrNo         int     `json:"sr_no"`
	MarketName   string  `json:"market_name"`
	BetType      string  `json:"bet_type"`
	BetOdds      float64 `json:"bet_odds"`
	BetStack     float64 `json:"bet_stack"`
	ResultStatus string  `json:"result_status"`
	BetResult    float64 `json:"bet_result"`
	BetTime      string  `json:"bet_time"`
	BetID        string  `json:"bet_id"`
	TrClass      string  `json:"tr_class"`
}

type CurrentBetRequest struct {
	SEcho          int    `form:"sEcho" json:"sEcho"`
	IDisplayStart  int    `form:"iDisplayStart" json:"iDisplayStart"`
	IDisplayLength int    `form:"iDisplayLength" json:"iDisplayLength"`
	SSearch        string `form:"sSearch" json:"sSearch"`
	ReportType     string `form:"report_type" json:"report_type"` // "sports" or "casino"
	BetType        string `form:"BetType" json:"BetType"`         // "back" or "lay"
}

type CurrentBetData struct {
	EventTypeLabel string  `json:"event_type_label"`
	EventName      string  `json:"event_name"`
	MarketName     string  `json:"market_name"`
	Nation         string  `json:"nation"`
	Datetime       string  `json:"datetime"`
	UserRate       float64 `json:"user_rate"`
	Amount         float64 `json:"amount"`
	BetType        string  `json:"bet_type"`
	Action         string  `json:"action"`
}

type CurrentBetResponse struct {
	SEcho                int              `json:"sEcho"`
	ITotalRecords        int              `json:"iTotalRecords"`
	ITotalDisplayRecords int              `json:"iTotalDisplayRecords"`
	AaData               []CurrentBetData `json:"aaData"`
	TotalAmount          float64          `json:"total_amount"`
	TotalBets            int              `json:"total_bets"`
	Ttt                  string           `json:"ttt"`
}

type ActivityLogRequest struct {
	SEcho          int    `form:"sEcho" json:"sEcho"`
	IDisplayStart  int    `form:"iDisplayStart" json:"iDisplayStart"`
	IDisplayLength int    `form:"iDisplayLength" json:"iDisplayLength"`
	SSearch        string `form:"sSearch" json:"sSearch"`
	FromDate       string `form:"from_date" json:"from_date"`
	ToDate         string `form:"to_date" json:"to_date"`
	ReportType     string `form:"report_type" json:"report_type"` // "endlogin" or "password"
}

type ActivityLogData struct {
	User    string `json:"user"`
	Date    string `json:"date"`
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type ActivityLogResponse struct {
	SEcho           int               `json:"sEcho"`
	RecordsTotal    int               `json:"recordsTotal"`
	RecordsFiltered int               `json:"recordsFiltered"`
	Data            []ActivityLogData `json:"data"`
}

type UpdateButtonValueRequest struct {
	AllButtonValue string `form:"all_button_value" json:"all_button_value" binding:"required"`
	Type           string `form:"type" json:"type"`
}

type UpdateButtonValueResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RefreshBalanceResponse struct {
	Status   string `json:"status"`
	Balance  string `json:"balance"`
	Exposure string `json:"exposure"`
	Winning  string `json:"winning"`
}
