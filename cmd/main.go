package main

import (
	"fmt"
	"log"
	"net/http"

	"GoTestTask/pkg/api"
	"GoTestTask/pkg/db"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

const (
	baseURL            = "localhost:8080"
	updatePostfix      = "/update"
	statePostfix       = "/state"
	getNamePostfix     = "/get_names"
	cleanTablesPostfix = "/clean"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.Set("state", db.StateEmpty)
}

func main() {
	fmt.Println("Starting application...")
	fmt.Println("Getting config data...")
	configData := getConfigData()

	fmt.Println("Connecting to database...")
	database, err := db.InitGorm(configData)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	fmt.Println("Connecting succeed...")

	App := api.App{database}
	r := chi.NewRouter()
	r.Get(updatePostfix, App.UpdateHandler)
	r.Get(statePostfix, App.StateHandler)
	r.Get(getNamePostfix, App.GetNamesHandler)
	r.Get(cleanTablesPostfix, App.CleanTablesHandler)

	err = http.ListenAndServe(baseURL, r)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigData() db.DbConfig {
	return db.DbConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Dbname:   viper.GetString("database.dbname")}
}
