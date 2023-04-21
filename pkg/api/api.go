package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kpotier/carburant/pkg/gas"
)

// Refresh is a variable that determines how often the gas price data is fetched
// from the French government API.
var Refresh = 30 * time.Minute

// API structure is the container for the fetched gas price data from the French
// government. It holds information such as the current price, the date and time
// of the latest update, and any other relevant metadata. This structure is used
// by the application to access and manipulate the gas price data.
type API struct {
	log *log.Logger

	done   chan bool
	ticker *time.Ticker

	mux  sync.RWMutex
	data []gas.Data

	muxDB  sync.RWMutex
	pathDB string
	db     map[string][]gas.DataGas
}

// The Start function initializes an instance of the API structure and launches
// a goroutine to periodically refresh the gas price data from the French
// government API. This function takes in a logger to log any relevant
// information or errors and a path to the database to store the fetched data.
// Once called, the Start function will continue to run in the background,
// periodically refreshing the gas price data and storing it in the specified
// database. The initialized API instance is returned by this function and can
// be used by the application to access and manipulate the gas price data. Stop
// may be called to kill the launched goroutine.
func Start(l *log.Logger, path string) (*API, error) {
	a := &API{
		log:    l,
		pathDB: path,
		mux:    sync.RWMutex{},
	}

	// Open database
	err := a.dbGet()
	if err != nil {
		return nil, fmt.Errorf("dbGet: %w", err)
	}

	// Fetch data
	err = a.fetch()
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	// Start ticker and refresher
	a.done = make(chan bool)
	a.ticker = time.NewTicker(Refresh)
	go func() {
		for {
			select {
			case <-a.ticker.C:
				err := a.fetch()
				if err != nil {
					a.log.Fatal("fetch:", err)
					return
				}
			case <-a.done:
				return
			}
		}
	}()

	return a, nil
}

// Stop kills the goroutine instanced in Start.
func (a *API) Stop() {
	if a.ticker != nil {
		a.done <- true
		a.ticker.Stop()
	}
}

// The Handler function returns an http.Handler that can be used with the http
// package to handle incoming HTTP requests. The Handler function allows the
// application to expose the fetched gas price data over a network via HTTP and
// can be used to build a simple RESTful API for the gas price data."
func (a *API) Handler() http.Handler {
	return a
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.log.Printf("api: %s: %s %s", r.RemoteAddr, r.Method, r.URL)
	mux := http.NewServeMux()
	mux.HandleFunc("/stations", a.getStations)
	mux.HandleFunc("/history", a.getHistory)
	mux.HandleFunc("/favorites", a.getFavorites)
	mux.ServeHTTP(w, r)
}

func (a *API) fetch() error {
	data, err := gas.Fetch()
	if err != nil {
		return err
	}

	// Put the latest data into memory
	a.mux.Lock()
	a.data = data
	a.mux.Unlock()

	// Save the prices into db
	err = a.dbUpdate()
	return err
}
