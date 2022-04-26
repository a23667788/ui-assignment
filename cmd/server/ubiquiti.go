package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/a23667788/ui-assignment/internal/client/postgres"
	"github.com/a23667788/ui-assignment/internal/entity"
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
	// search an user by fullname.
	getR.HandleFunc("/user/{fullname}", ui.getUser)
	// get the user’s detailed information.
	getR.HandleFunc("/userDetail/{account}", ui.getUserDetail)

	postR := ui.Router.Methods(http.MethodPost).Subrouter()
	// create the user (user sign up).
	postR.HandleFunc("/user", ui.createUser)

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

func (ui *Ubiquiti) getUser(w http.ResponseWriter, r *http.Request) {
	log.Info("getUser start")
	defer log.Info("getUser done")

	client := postgres.DBClient{}
	client.Connect()

	vars := mux.Vars(r)
	fullname := vars["fullname"]

	res, err := client.Get(fullname)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Debug(res)

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (ui *Ubiquiti) getUserDetail(w http.ResponseWriter, r *http.Request) {
	log.Info("getUserDetail start")
	defer log.Info("getUserDetail done")

	client := postgres.DBClient{}
	client.Connect()

	vars := mux.Vars(r)
	account := vars["account"]

	res, err := client.GetUserDetail(account)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Debug(res)

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (ui *Ubiquiti) createUser(w http.ResponseWriter, r *http.Request) {
	log.Info("createUser start")
	defer log.Info("createUser done")

	client := postgres.DBClient{}
	client.Connect()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Debug(string(body))

	var req entity.CreateUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = client.Insert(req)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
