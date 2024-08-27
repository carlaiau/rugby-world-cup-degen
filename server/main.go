package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/carlaiau/degen/types"

	the_odds_api "github.com/carlaiau/degen/the-odds-api"
)

var (
	handlerHit    bool
	lastRun       time.Time
	ticker        *time.Ticker
	currentEvents map[string]types.DegenEvent
)

type proxyIp struct {
	Ip   string
	Port int
}

func getData() {

	/*
		err := tab.Scraper(currentEvents)
		if err != nil {
			fmt.Printf("Error scraping TAB: %v", err)
		}
	*/

	the_odds_api.Scrapper(currentEvents)
}

func convertMapsIntoJsonArrays(currentEvents map[string]types.DegenEvent) []types.DegenEventForClient {
	// Restructure the data and sort it
	events := []types.DegenEventForClient{}

	for _, event := range currentEvents {

		markets := make([]types.DegenMarket, 0, len(event.Markets))
		for _, v := range event.Markets {
			markets = append(markets, v)
		}
		// Sort the values by StartTime
		sort.Slice(markets, func(j, k int) bool {
			return markets[j].Name < markets[k].Name
		})

		events = append(events, types.DegenEventForClient{
			StartTime: event.StartTime,
			Name:      event.Name,
			Markets:   markets,
			Home:      event.Home,
			Away:      event.Away,
		})

	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return events
}
func handler(w http.ResponseWriter, r *http.Request) {
	handlerHit = true

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Vary", "Origin")
	json.NewEncoder(w).Encode(convertMapsIntoJsonArrays(currentEvents))

}

func main() {

	handlerHit = true
	currentEvents = make(map[string]types.DegenEvent)
	ticker = time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if handlerHit && time.Since(lastRun) >= 5*time.Minute {
					getData()
					lastRun = time.Now()
					handlerHit = false
				}
			}
		}
	}()

	http.HandleFunc("/", handler)

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	// Note the Port value used here.
	http.ListenAndServe(":"+httpPort, nil)

	fmt.Println("Server started on port " + httpPort)

}
