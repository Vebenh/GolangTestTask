package api

import (
	"GoTestTask/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"sync"

	"gorm.io/gorm"

	parser "GoTestTask/pkg/xmlparser"
)

const (
	ParseURL = "https://www.treasury.gov/ofac/downloads/sdn.xml"
)

type App struct {
	DB *gorm.DB
}

type Response struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
	Code   int    `json:"code,omitempty"`
}

func (app *App) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	entries := make(chan db.SdnEntry)
	errors := make(chan error, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Update has started... ")
	viper.Set("state", db.StateUpdating)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := parser.ParseXML(ctx, ParseURL, entries); err != nil {
			errors <- err
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := db.WriteToDB(ctx, app.DB, entries); err != nil {
			errors <- err
			cancel()
		}
	}()

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			fmt.Println("Update failed:", ctx.Err())
			response := Response{
				Result: false,
				Info:   "service unavailable",
				Code:   http.StatusServiceUnavailable,
			}
			sendResponse(response, w)
			return
		}
	}

	response := Response{
		Result: true,
		Info:   "",
		Code:   http.StatusOK,
	}
	sendResponse(response, w)

	viper.Set("state", db.StateOk)
	fmt.Println("Update is completed... ")
}

func (app *App) StateHandler(w http.ResponseWriter, r *http.Request) {
	switch viper.Get("state") {
	case db.StateEmpty:
		response := Response{
			Result: false,
			Info:   "empty",
		}
		sendResponse(response, w)
	case db.StateUpdating:
		response := Response{
			Result: false,
			Info:   "updating",
		}
		sendResponse(response, w)
	case db.StateOk:
		response := Response{
			Result: true,
			Info:   "ok",
		}
		sendResponse(response, w)
	}
}

func (app *App) GetNamesHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Println("Name:", name)
	searchType := r.URL.Query().Get("type")
	fmt.Println("Type:", searchType)

	if searchType != "strong" {
		searchType = "weak"
	}

	persons, err := db.GetPerson(app.DB, name, searchType)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	sendResponse(persons, w)
}

func (app *App) CleanTablesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start cleaning... ")
	app.DB.Delete(&db.Program{}, true)
	app.DB.Delete(&db.Aka{}, true)
	app.DB.Delete(&db.Address{}, true)
	app.DB.Delete(&db.PublishInformation{}, true)
	app.DB.Delete(&db.SdnEntry{}, true)
	viper.Set("state", db.StateEmpty)
	fmt.Println("Cleaning completed...")
}

func sendResponse[T any](resp T, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
