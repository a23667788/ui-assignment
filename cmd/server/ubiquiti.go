package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Ubiquiti struct {
	Router *mux.Router
	DB     *sql.DB
}

func (ui *Ubiquiti) Initialize(user, password, dbname string) {

	ui.Router = mux.NewRouter()
	ui.initializeRoutes()
}

func (ui *Ubiquiti) Run() {
	log.Fatal(http.ListenAndServe(":8000", ui.Router))
}

func (ui *Ubiquiti) initializeRoutes() {
	//list all users.
	ui.Router.HandleFunc("/users", ui.getUsersHandler).Methods(http.MethodGet)

	// search an user by fullname.
	ui.Router.HandleFunc("/username/fullname:[0-9+]}", ui.getUsernameHandler).Methods(http.MethodGet)

	ui.Router.HandleFunc("/username/id:[0-9]+", ui.updateUsername).Methods(http.MethodPatch)

	// get the user’s detailed information.
	ui.Router.HandleFunc("/user/{id:[0-9]+}", ui.getUsertHandler).Methods(http.MethodGet)

	// create the user (user sign up).
	ui.Router.HandleFunc("/user", ui.createUserHandler).Methods(http.MethodPost)

	// delete the user.
	ui.Router.HandleFunc("/user", ui.deleteUserHandler).Methods(http.MethodDelete)

	// update the user.
	ui.Router.HandleFunc("/user", ui.updateUserHandler).Methods(http.MethodPatch)

	// generate the token to the user (user sign in).
	ui.Router.HandleFunc("/userSessions", ui.UserSessionsHandler)
}

func (ui *Ubiquiti) getUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func (ui *Ubiquiti) getUsernameHandler(w http.ResponseWriter, r *http.Request) {

}

func (ui *Ubiquiti) updateUsername(w http.ResponseWriter, r *http.Request) {

}

func (ui *Ubiquiti) getUsertHandler(w http.ResponseWriter, r *http.Request) {

}

func (ui *Ubiquiti) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// decode json first
	// _ := json.NewDecoder(r.Body)

}

func (ui *Ubiquiti) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
func (ui *Ubiquiti) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (ui *Ubiquiti) UserSessionsHandler(w http.ResponseWriter, r *http.Request) {

}
