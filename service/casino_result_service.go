package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kirtanwyn/allPanelexch/dto"
	"github.com/kirtanwyn/allPanelexch/repository"
)

type CasinoResultService interface {
	GetCasinoResults(req dto.CasinoResultRequest) (dto.CasinoResultResponse, error)
	GetResultCardsHTML(req dto.GetResultCardsRequest) (string, error)
}

type casinoResultService struct {
	repo repository.CasinoResultRepository
}

func NewCasinoResultService(repo repository.CasinoResultRepository) CasinoResultService {
	return &casinoResultService{repo: repo}
}

func (s *casinoResultService) GetCasinoResults(req dto.CasinoResultRequest) (dto.CasinoResultResponse, error) {
	gameDate := req.GameDate
	// For datatable, length is iDisplayLength, start is iDisplayStart
	length := req.IDisplayLength
	start := req.IDisplayStart
	
	results, totalData, err := s.repo.GetCasinoResults(req.CasinoType, gameDate, req.SSearch, length, start)
	if err != nil {
		return dto.CasinoResultResponse{}, err
	}

	var data []dto.CasinoResultItem
	for _, row := range results {
		eventID := row["event_id"].(string)
		gameType := row["game_type"].(string)
		resultStatus := row["result_status"].(string)
		
		resultStatusText := ""
		if req.CasinoType == "poker" || req.CasinoType == "poker20" {
			resultStatusText = "Poker " + resultStatus + " - "
		} else if req.CasinoType == "poker6" || req.CasinoType == "card32" || req.CasinoType == "card32eu" {
			resultStatusText = "Player " + resultStatus + " - "
		}

		link := fmt.Sprintf(`<span onclick="view_casiono_result(%s)">%s</span>`, eventID, eventID)
		
		data = append(data, dto.CasinoResultItem{
			Round:  link,
			Winner: resultStatusText + " " + gameType,
		})
	}

	if data == nil {
		data = []dto.CasinoResultItem{}
	}

	return dto.CasinoResultResponse{
		SEcho:           req.SEcho,
		RecordsTotal:    totalData,
		RecordsFiltered: totalData, // simplified
		Data:            data,
	}, nil
}

func (s *casinoResultService) GetResultCardsHTML(req dto.GetResultCardsRequest) (string, error) {
	cardsStr, resultStatus, descRemarks, err := s.repo.GetCasinoResultCards(req.CasinoType, req.EventID)
	if err != nil {
		return "", err
	}

	var cards []string
	if err := json.Unmarshal([]byte(cardsStr), &cards); err != nil {
		// handle non-json or empty gracefully
	}

	var htmlBuilder strings.Builder
	// Basic implementation for teen games
	if req.CasinoType == "teen41" || req.CasinoType == "teen42" || req.CasinoType == "teen33" || req.CasinoType == "teen32" || req.CasinoType == "teen6" || req.CasinoType == "teen3" || req.CasinoType == "teen" || req.CasinoType == "teen62" {
		rdescArray := strings.Split(descRemarks, "#")
		
		webURL := "" // define if needed
		
		// Fallbacks for card lengths
		cardA1, cardA2, cardA3 := "", "", ""
		cardB1, cardB2, cardB3 := "", "", ""
		
		if len(cards) > 0 { cardA1 = webURL + "https://allpanelexch.tv//storage/front/img/cards_new/" + cards[0] + ".png" }
		if len(cards) > 2 { cardA2 = webURL + "https://allpanelexch.tv//storage/front/img/cards_new/" + cards[2] + ".png" }
		if len(cards) > 4 { cardA3 = webURL + "https://allpanelexch.tv//storage/front/img/cards_new/" + cards[4] + ".png" }
		
		if len(cards) > 1 { cardB1 = webURL + "https://allpanelexch.tv//storage/front/img/cards_new/" + cards[1] + ".png" }
		if len(cards) > 3 { cardB2 = webURL + "https://allpanelexch.tv//storage/front/img/cards_new/" + cards[3] + ".png" }
		if len(cards) > 5 { cardB3 = webURL + "https://allpanelexch.tv//storage/front/img/cards_new/" + cards[5] + ".png" }
		
		htmlBuilder.WriteString(`<div class="d-flex justify-content-between">
			<div class="d-flex flex-column text-center">
				<h4>Player A</h4>
				<div class="result-image d-flex align-items-center  gap-2">`)
				
		if resultStatus == "1" {
			htmlBuilder.WriteString(`<div class="winner-icon mt-3" style="position: unset;"><i class="fas fa-trophy mr-2"></i></div>`)
		}
		
		htmlBuilder.WriteString(fmt.Sprintf(`<img src="%s" class="mr-2"><img src="%s" class="mr-2"><img src="%s" class="mr-2">`, cardA1, cardA2, cardA3))
		htmlBuilder.WriteString(`</div></div>
			<div class="d-flex flex-column text-center">
				<h4>Player B</h4>
				<div class="result-image d-flex align-items-center  gap-2">`)
		
		htmlBuilder.WriteString(fmt.Sprintf(`<img src="%s" class="mr-2"><img src="%s" class="mr-2"><img src="%s" class="mr-2">`, cardB1, cardB2, cardB3))
		
		if resultStatus == "2" {
			htmlBuilder.WriteString(`<div class="winner-icon mt-3" style="position: unset;"><i class="fas fa-trophy mr-2"></i></div>`)
		}
		
		htmlBuilder.WriteString(`</div></div></div>`)

		if req.CasinoType == "teen" || req.CasinoType == "teen62" {
			winner, oddEven, cons := "", "", ""
			if len(rdescArray) > 0 { winner = rdescArray[0] }
			if len(rdescArray) > 2 { oddEven = rdescArray[2] }
			if len(rdescArray) > 3 { cons = rdescArray[3] }
			
			htmlBuilder.WriteString(fmt.Sprintf(`<div class="row mt-2 justify-content-center">
				<div class="col-md-6">
					<div class="casino-result-desc" style="display: -webkit-flex;flex-wrap: wrap;padding: 6px;box-shadow: 0 0 4px -1px;margin-top:10px">
						<div class="casino-result-desc-item" style="display: flex;justify-content: center;width: 100%%; flex-wrap: wrap;">
							<div><span style="color:gray">Winner: </span>%s</div>
							<div><span style="color:gray">Odd/Even: </span>%s</div>
							<div><span style="color:gray">Consecutive: </span>%s</div>
						</div>
					</div>
				</div>
			</div>`, winner, oddEven, cons))
		}
	} else {
		htmlBuilder.WriteString(`<div class="alert alert-warning">Result view for this game type is not fully migrated to the new API yet. Please check old PHP file for reference.</div>`)
	}

	return htmlBuilder.String(), nil
}
