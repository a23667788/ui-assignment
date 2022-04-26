package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/a23667788/ui-assignment/internal/client/postgres"
	"github.com/a23667788/ui-assignment/internal/entity"
	"github.com/a23667788/ui-assignment/internal/token"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Ubiquiti struct {
	Router *mux.Router
	DB     *sql.DB
	Jwt    token.JWT
	Ws     *websocket.Conn
}

func (ui *Ubiquiti) Initialize() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	ui.Router = mux.NewRouter()
	ui.initializeRoutes()

	// jwt setting
	prvKey, err := ioutil.ReadFile("assets/jwtRS256.key")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	pubKey, err := ioutil.ReadFile("assets/jwtRS256.key.pub")
	if err != nil {
		log.Error(err)
		panic(err)
	}

	ui.Jwt = token.NewJWT(prvKey, pubKey)
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
	getR.Use(ui.authMiddleware)

	postR := ui.Router.Methods(http.MethodPost).Subrouter()
	// create the user (user sign up).
	postR.HandleFunc("/user", ui.createUser)
	// generate the token to the user (user sign in).
	postR.HandleFunc("/userSession", ui.userSession)

	deleteR := ui.Router.Methods(http.MethodDelete).Subrouter()
	// delete the user.
	deleteR.HandleFunc("/user/{account}", ui.deleteUser)
	deleteR.Use(ui.authMiddleware)

	patchR := ui.Router.Methods(http.MethodPatch).Subrouter()
	// update the user.
	patchR.HandleFunc("/user/{account}", ui.updateUser)
	// update user’s fullname.
	patchR.HandleFunc("/username/{account}", ui.updateFullname)
	patchR.Use(ui.authMiddleware)

	// handle websocket
	ui.Router.HandleFunc("/ws", ui.wsEndpoint)
}

// listUsers godoc
// @Summary List all users
// @Description This is the description for listing user.
// @Tags Thing
// @Param paging query int false "paging"
// @Param sorting query int false "sorting"
// @Success 200 {object} entity.ListUsersResponse
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /users [get]
func (ui *Ubiquiti) listUsers(w http.ResponseWriter, r *http.Request) {
	log.Info("listUsers start")
	defer log.Info("listUsers done")

	var pagingSlice, sortingSlice []string
	vars := r.URL.Query()
	pagingSlice = vars["paging"]
	sortingSlice = vars["sorting"]

	// use first number
	var paging string
	if pagingSlice != nil {
		paging = pagingSlice[0]
	}

	sorting := strings.Join(sortingSlice[:], " ")

	client := postgres.DBClient{}
	client.Connect()

	res, err := client.List(paging, sorting)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	log.Debug(res)

	respondWithJSON(w, http.StatusOK, res)
}

// getUser godoc
// @Summary Get an user by fullname.
// @Description This is the description for getting user.
// @Tags Thing
// @Param fullname path string true "fullname"
// @Success 200 {object} entity.GetUser
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /user/{fullname} [get]
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

// getUserDetail godoc
// @Summary Get user’s detailed information.
// @Description This is the description for getting user detail inform.
// @Tags Thing
// @Param account path string true "account"
// @Success 200 {object} entity.User
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /userDetail/{account} [get]
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

// createUser godoc
// @Summary Create a user
// @Description This is the description for creating a user.
// @Tags Thing
// @accept application/x-www-form-urlencoded
// @produce application/json
// @Param acct formData string true "account"
// @Param fullname formData string true "fullname"
// @Param pwd formData string true "password"
// @Success 200 {object} entity.CreateUserResponse
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /user [post]
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

	var req entity.User
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

// userSession godoc
// @Summary Create a usersession
// @Description This is the description for creating a usersession.
// @Tags Thing
// @accept application/x-www-form-urlencoded
// @produce application/json
// @Param acct formData string true "account"
// @Param pwd formData string true "password"
// @Success 200 {object} entity.UserSessionResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /userSession [post]
func (ui *Ubiquiti) userSession(w http.ResponseWriter, r *http.Request) {
	log.Info("userSession start")
	defer log.Info("userSession done")

	client := postgres.DBClient{}
	client.Connect()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// create session
	var req entity.UserSessionRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = client.Validate(req)
	if err != nil {
		log.Error(err)

		// Send login err to ws
		text := []byte("login error")
		if err := ui.Ws.WriteMessage(websocket.TextMessage, text); err != nil {
			log.Println(err)

		}

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tok, err := ui.Jwt.Create(time.Minute, req.Acct)
	if err != nil {
		log.Fatalln(err)
	}

	var resp entity.UserSessionResponse
	resp.Jwt = tok
	respondWithJSON(w, http.StatusOK, resp)
}

// deleteUser godoc
// @Summary Delete a user
// @Description This is the description for deleting a user.
// @Tags Thing
// @produce application/json
// @Param account path string true "account"
// @Success 200 {object} entity.DeleteUserResponse
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /user/{account} [delete]
func (ui *Ubiquiti) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Info("deleteUser start")
	defer log.Info("deleteUser done")

	client := postgres.DBClient{}
	client.Connect()

	vars := mux.Vars(r)
	account := vars["account"]

	err := client.Delete(account)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// updateUser godoc
// @Summary Update a user
// @Description This is the description for updating a user.
// @Tags Thing
// @accept json
// @produce application/json
// @Param account path string true "account"
// @Success 200 {object} entity.UpdateUserResponse
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /user/{account} [patch]
func (ui *Ubiquiti) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Info("updateUser start")
	defer log.Info("updateUser done")

	client := postgres.DBClient{}
	client.Connect()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	vars := mux.Vars(r)
	account := vars["account"]

	log.Debug(string(body))

	var user entity.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = client.Update(account, user)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// updateFullname godoc
// @Summary Update a user's fullname
// @Description This is the description for updating a user's fullname.
// @Tags Thing
// @accept json
// @Param account path string true "account"
// @Success 200 {object} entity.UpdateUserResponse
// @Success 401 {object} entity.ErrorResponse
// @Success 403 {object} entity.ErrorResponse
// @Router /username/{account} [patch]
func (ui *Ubiquiti) updateFullname(w http.ResponseWriter, r *http.Request) {
	log.Info("updateUsername start")
	defer log.Info("updateUsername done")

	client := postgres.DBClient{}
	client.Connect()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Debug(string(body))

	vars := mux.Vars(r)
	account := vars["account"]

	var req entity.UpdateFullnameRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = client.UpdateFullname(account, req)
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		// read a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// log.Println(err)
			return
		}

		log.Info("reader", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			// log.Println(err)
			return
		}
	}
}

func (ui *Ubiquiti) wsEndpoint(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	var err error
	ui.Ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}

	reader(ui.Ws)

}

func (ui *Ubiquiti) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := ui.Jwt.Validate(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}

		dat := claims.(string)

		r.Header.Set("dat", dat)

		next.ServeHTTP(w, r)
	})
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
