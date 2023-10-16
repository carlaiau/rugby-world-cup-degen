package the_odds_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/carlaiau/degen/types"
	"github.com/joho/godotenv"
)

const (
	API_URL   = "https://api.the-odds-api.com/v4"
	SPORT_KEY = "rugbyunion_world_cup"
)

func Scrapper(currentEvents map[string]types.DegenEvent) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	API_KEY := os.Getenv("THE_ODDS_API_KEY")

	url := fmt.Sprintf("%s/sports/%s/odds/?regions=us,us2,uk,au,eu&oddsFormat=decimal&apiKey=%s&markets=totals,spreads", API_URL, SPORT_KEY, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data from API: %v", err)
		return err
	}
	defer resp.Body.Close()

	res := []Event{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
		return err
	}

	for _, event := range res {
		for _, bookmaker := range event.Bookmakers {
			for _, market := range bookmaker.Markets {
				name := ""
				switch market.Key {
				case "totals":
					name = "TOTAL"
				case "h2h":
					name = "H2H"
				case "spreads":
					name = "LINE"
				default:
					continue
				}

				newMarket := types.DegenMarket{
					Bookmaker:  bookmaker.Title,
					Name:       name,
					LastUpdate: bookmaker.LastUpdate,
					Line:       0,
					Outcomes:   []types.DegenOutcome{},
				}
				for i, outcome := range market.Outcomes {
					if i == 0 && name != "H2H" {
						newMarket.Line = outcome.Point
					}
					displayOrder := 0
					if outcome.Name != event.HomeTeam {
						displayOrder = 1
					}
					newMarket.Outcomes = append(newMarket.Outcomes, types.DegenOutcome{
						DisplayOrder: displayOrder,
						Price:        outcome.Price,
					})
				}

				// Down here deterime if this is a new event or the event already exists
				eventKey := fmt.Sprintf("%s-%s", event.HomeTeam, event.AwayTeam)
				marketKey := fmt.Sprintf("%s-%s", bookmaker.Key, name)

				if _, ok := currentEvents[eventKey]; ok {
					currentEvents[eventKey].Markets[marketKey] = newMarket
				} else {
					newEvent := types.DegenEvent{
						StartTime: event.CommenceTime,
						Name:      event.HomeTeam + " vs " + event.AwayTeam,
						Markets:   map[string]types.DegenMarket{},
						Home:      event.HomeTeam,
						Away:      event.AwayTeam,
					}
					newEvent.Markets[marketKey] = newMarket
					currentEvents[eventKey] = newEvent
				}

			}
		}

	}
	return nil
}
