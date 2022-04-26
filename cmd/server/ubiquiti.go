package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/a23667788/ui-assignment/internal/client/postgres"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Ubiquiti struct {
	Router *mux.Router
	DB     *sql.DB
}

func (ui *Ubiquiti) Initialize() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	ui.Router = mux.NewRouter()
	ui.initializeRoutes()
}

func (ui *Ubiquiti) Run() {
	log.Fatal(http.ListenAndServe(":8000", ui.Router))
}

func (ui *Ubiquiti) initializeRoutes() {

	getR := ui.Router.Methods(http.MethodGet).Subrouter()
	// list all users.
	getR.HandleFunc("/users", ui.listUsers)
}

func (ui *Ubiquiti) listUsers(w http.ResponseWriter, r *http.Request) {
	log.Info("listUsers start")
	defer log.Info("listUsers done")

	client := postgres.DBClient{}
	client.Connect()

	res, err := client.List()
	if err != nil {
		log.Error(err)
		panic(err)
	}

	log.Debug(res)

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	respondWithJSON(w, http.StatusOK, jsonResponse)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
