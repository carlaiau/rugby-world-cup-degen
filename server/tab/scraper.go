package tab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/carlaiau/degen/types"
)

const APIEndpoint = "https://content.tab.co.nz/content-service/api/v1/q/time-band-event-list?marketGroupTypesIncluded=MONEYLINE%2CROLLING_SPREAD%2CROLLING_TOTAL%2CSTATIC_SPREAD%2CSTATIC_TOTAL&allowedEventSorts=MTCH&includeChildMarkets=true&prioritisePrimaryMarkets=true&maxTotalItems=60&maxEventsPerCompetition=7&maxCompetitionsPerSportPerBand=3&maxEventsForNextToGo=5&startTimeOffsetForNextToGo=600&maxMarkets=14&excludeEventsWithNoMarkets=true&includeCommentary=false&includeMedia=false&drilldownTagIds=18525"

const EventEndpoint = "https://content.tab.co.nz/content-service/api/v1/q/event-list?maxEvents=12&orderEventsBy=startTime&maxMarkets=10&orderMarketsBy=displayOrder&marketSortsIncluded=HL%2CWH&marketGroupTypesIncluded=MONEYLINE%2CROLLING_SPREAD%2CROLLING_TOTAL%2CSTATIC_SPREAD%2CSTATIC_TOTAL&excludeEventsWithNoMarkets=false&eventSortsIncluded=MTCH&includeChildMarkets=true&prioritisePrimaryMarkets=true&includeCommentary=false&includeMedia=false&drilldownTagIds=18525"

func SortOutcomes(market *Market) {
	sort.Slice(market.Outcomes, func(i, j int) bool {
		return market.Outcomes[i].DisplayOrder < market.Outcomes[j].DisplayOrder
	})
}

func Scraper(currentEvents map[string]types.DegenEvent) error {

	resp, err := http.Get(EventEndpoint)
	if err != nil {
		fmt.Printf("TAB: Error fetching data from API: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("TAB: Error reading response body: %v\n", err)
		return err
	}

	fmt.Println(string(body))

	var res TimeBandEventResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Printf("TAB: Error unmarshalling response: %v\n", err)
		return err
	}

	events := res.Data.Events

	for _, event := range events {
		if !strings.Contains(event.Meeting.Name, "World Cup") {
			continue
		}
		for _, market := range event.Markets {
			SortOutcomes(&market)
			switch market.GroupCode {

			case "MATCH_RESULT_2_WAY_INC_ET", "ROLLING_HANDICAP_2_WAY_MIDDLE_LINE", "ROLLING_TOTAL_POINTS_O/U_MIDDLE_LINE":
				name := "LINE"
				if market.GroupCode == "ROLLING_TOTAL_POINTS_O/U_MIDDLE_LINE" {
					name = "TOTAL"
				}
				if market.GroupCode == "MATCH_RESULT_2_WAY_INC_ET" {
					name = "H2H"
				}
				var handicapValue float64
				if name != "H2H" {
					if v, ok := market.HandicapValue.(float64); ok {
						handicapValue = v
					}
				}
				newMarket := types.DegenMarket{
					Bookmaker:  "NZ TAB",
					Name:       name,
					LastUpdate: time.Now(),
					Line:       handicapValue,
					Outcomes:   []types.DegenOutcome{},
				}

				for _, outcome := range market.Outcomes {
					newMarket.Outcomes = append(newMarket.Outcomes, types.DegenOutcome{
						DisplayOrder: outcome.DisplayOrder,
						Price:        outcome.Prices[0].Decimal,
					})
				}

				// Here is where we check to see if a the event already exists and if so then we
				eventKey := fmt.Sprintf("%s-%s", event.Teams[0].Name, event.Teams[1].Name)
				marketKey := fmt.Sprintf("nz_tab-%s", newMarket.Name)

				if _, ok := currentEvents[eventKey]; ok {
					currentEvents[eventKey].Markets[marketKey] = newMarket
				} else {
					newEvent := types.DegenEvent{
						StartTime: event.StartTime,
						Name:      event.Name,
						Markets:   map[string]types.DegenMarket{},
						Home:      event.Teams[0].Name,
						Away:      event.Teams[1].Name,
					}
					newEvent.Markets[marketKey] = newMarket
					currentEvents[eventKey] = newEvent
				}
			}
		}

	}

	return nil
}
