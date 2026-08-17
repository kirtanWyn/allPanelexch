package model

// EventMarketId represents the event_market_id table — maps event/market/type combinations.
type EventMarketId struct {
	ID         int    `json:"id" db:"id"`
	EventID    int    `json:"event_id" db:"event_id"`
	MarketID   string `json:"market_id" db:"market_id"`
	RunnerID   string `json:"runner_id" db:"runner_id"`
	MarketName string `json:"market_name" db:"market_name"`
	MarketType string `json:"market_type" db:"market_type"`
}

// EventMarketLimit represents the event_market_limit table — min/max limits per event.
type EventMarketLimit struct {
	ID           int     `json:"id" db:"id"`
	EventID      int     `json:"event_id" db:"event_id"`
	OddsMarketID string  `json:"oddsmarketId" db:"oddsmarketId"`
	MatchMin     int     `json:"match_min" db:"match_min"`
	MatchMax     int     `json:"match_max" db:"match_max"`
	BookmakerMin int     `json:"bookmaker_min" db:"bookmaker_min"`
	BookmakerMax int     `json:"bookmaker_max" db:"bookmaker_max"`
	BookmakerLive int    `json:"bookmaker_live" db:"bookmaker_live"`
	Status       int     `json:"status" db:"status"`
	EventName    string  `json:"event_name" db:"event_name"`
}

// BetDelayMaster represents the bet_delay_master table — configurable delays per market type.
type BetDelayMaster struct {
	ID             int     `json:"id" db:"id"`
	MarketTypeID   string  `json:"market_type_id" db:"market_type_id"`
	DelayValue     float64 `json:"delay_value" db:"delay_value"`
}

// EventDelayMaster represents the event_delay_master table — per-event delay overrides.
type EventDelayMaster struct {
	ID      int    `json:"id" db:"id"`
	EventID string `json:"event_id" db:"event_id"`
}

// BetMarketSuspendMaster represents the bet_market_suspend_master table — suspend rules.
type BetMarketSuspendMaster struct {
	ID         int    `json:"id" db:"id"`
	SportID    string `json:"sport_id" db:"sport_id"`
	EventID    string `json:"event_id" db:"event_id"`
	MarketType string `json:"market_type" db:"market_type"`
}

// TossEndTime represents the toss_end_time table — toss betting cutoff.
type TossEndTime struct {
	ID      int    `json:"id" db:"id"`
	EventID string `json:"event_id" db:"event_id"`
	EndDate string `json:"end_date" db:"end_date"`
}
