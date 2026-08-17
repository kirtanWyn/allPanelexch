package model

// Account represents the accounts table used for user balance tracking.
// Each row stores a ledger entry (deposit, bet hold, win credit, etc.).
// The net balance is SUM(amount) WHERE user_id=? AND status=1.
type Account struct {
	ID       int     `json:"id" db:"id"`
	UserID   int     `json:"user_id" db:"user_id"`
	Amount   float64 `json:"amount" db:"amount"`
	Status   int     `json:"status" db:"status"`
	EventID  int     `json:"event_id" db:"event_id"`
	MarketID string  `json:"market_id" db:"market_id"`
	BetID    int     `json:"bet_id" db:"bet_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
