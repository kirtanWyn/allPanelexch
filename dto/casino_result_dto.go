package dto

type CasinoResultRequest struct {
	SEcho          int    `form:"sEcho" json:"sEcho"`
	IDisplayStart  int    `form:"iDisplayStart" json:"iDisplayStart"`
	IDisplayLength int    `form:"iDisplayLength" json:"iDisplayLength"`
	SSearch        string `form:"sSearch" json:"sSearch"`
	GameDate       string `form:"game_date" json:"game_date"`
	CasinoType     string `form:"casino_type" json:"casino_type"`
}

type CasinoResultItem struct {
	Round  string `json:"round"`
	Winner string `json:"winner"`
}

type CasinoResultResponse struct {
	SEcho           int                `json:"sEcho"`
	RecordsTotal    int                `json:"recordsTotal"`
	RecordsFiltered int                `json:"recordsFiltered"`
	Data            []CasinoResultItem `json:"data"`
}

type GetResultCardsRequest struct {
	EventID    string `form:"event_id" json:"event_id"`
	CasinoType string `form:"casino_type" json:"casino_type"`
}
