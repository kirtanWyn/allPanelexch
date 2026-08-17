package model

// BlockEvent represents the block_event table — blocks user from specific sport+event combos.
type BlockEvent struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"user_id" db:"UserId"`
	SportType  string `json:"sport_type" db:"sport_type"`
	CasinoName string `json:"casino_name" db:"casino_name"`
}

// BlockEventId represents the block_event_id table — blocks event/market by upline.
type BlockEventId struct {
	ID        int    `json:"id" db:"id"`
	EventType int    `json:"event_type" db:"event_type"`
	EventID   int    `json:"event_id" db:"event_id"`
	MarketID  string `json:"market_id" db:"market_id"`
	BlockBy   int    `json:"block_by" db:"block_by"`
}

// BlockSport represents the block_sport table — blocks entire sport type by upline.
type BlockSport struct {
	ID        int    `json:"id" db:"id"`
	EventType int    `json:"event_type" db:"event_type"`
	BlockBy   int    `json:"block_by" db:"block_by"`
}

// BetBlockDetails represents the bet_block_details table — blocks specific user on event.
type BetBlockDetails struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	EventID   int    `json:"event_id" db:"event_id"`
	BlockType int    `json:"block_type" db:"block_type"`
}
