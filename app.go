// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}
type statusRes struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(us-cdbr-iron-east-02.cleardb.net:3306)/%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	//passenger_details
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{ffn:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{ffn:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{ffn:[0-9]+}", a.deleteUser).Methods("DELETE")
	a.Router.HandleFunc("/userLogin", a.getUserLogin).Methods("POST")
	//lounge-login
	a.Router.HandleFunc("/loungeLogin", a.getLoungeLogin).Methods("POST")
	a.Router.HandleFunc("/createLoungeLogin", a.createLoungeLogin).Methods("POST")
	//lounge-booking-get-create
	a.Router.HandleFunc("/getLoungeBooking/{ticket_id}", a.getloungebooking).Methods("GET")
	a.Router.HandleFunc("/getLoungeBookings/{ffn:[0-9]+}", a.getloungebookingsbyffn).Methods("GET")
	a.Router.HandleFunc("/getUpcomingLoungeBookings/{lounge_id:[0-9]+}", a.getupcomingloungebookings).Methods("GET")
	a.Router.HandleFunc("/getCurrentLoungeBookings/{lounge_id:[0-9]+}", a.getcurrentloungebookings).Methods("GET")
	a.Router.HandleFunc("/createLoungeBooking", a.createLoungeBooking).Methods("POST")
	a.Router.HandleFunc("/getLoungeDetail/{lounge_id:[0-9]+}", a.getLoungeDetail).Methods("GET")
	a.Router.HandleFunc("/getLoungeDetails", a.getLoungeDetails).Methods("GET")
	a.Router.HandleFunc("/checkin", a.checkin).Methods("POST")
	a.Router.HandleFunc("/checkout", a.checkout).Methods("POST")
	//card-check
	a.Router.HandleFunc("/cardCheck", a.getcardetails).Methods("POST")
	//pnr-check
	a.Router.HandleFunc("/getpnr/{ffn:[0-9]+}", a.getpnr).Methods("GET")
}

func (a *App) getUserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t reqlogin
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	ffn := string(t.FFN)
	u := dblogin{FFN: ffn}
	if err := u.getUserLogin(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Passenger not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// compare password
	var passwordtes = bcrypt.CompareHashAndPassword([]byte(u.Pass), []byte(t.Pass))

	fmt.Println(string(t.Pass))
	fmt.Println(passwordtes)
	if passwordtes == nil {
		//login success

		res := statusRes{Status: 200, Result: "success"}
		json.NewEncoder(w).Encode(res)
	} else {
		//login failed
		res := statusRes{Status: 400, Result: "fail"}
		json.NewEncoder(w).Encode(res)
	}

}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ffn, err := strconv.Atoi(vars["ffn"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid FFN ID")
		return
	}
	u := passenger{FFN: ffn}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ffn, err := strconv.Atoi(vars["ffn"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid FFN ")
		return
	}
	u := passenger{FFN: ffn}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
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

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	users, err := getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u passenger
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u.FFN)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ffn, err := strconv.Atoi(vars["ffn"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var u passenger
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	u.FFN = ffn
	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getLoungeLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t reqloginlounge
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	loungeid := string(t.Loungeid)
	fmt.Println(loungeid)
	u := dbloginlounge{Loungeid: loungeid}
	if err := u.getLoungeLogin(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			res := statusRes{Status: 401, Result: "Lounge not found"}
			json.NewEncoder(w).Encode(res)
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// compare password
	var passwordtes = bcrypt.CompareHashAndPassword([]byte(u.Pass), []byte(t.Pass))

	fmt.Println(string(t.Pass))
	fmt.Println(passwordtes)
	if passwordtes == nil {
		//login success

		res := statusRes{Status: 200, Result: "Success"}
		json.NewEncoder(w).Encode(res)
	} else {
		//login failed
		res := statusRes{Status: 400, Result: "Failed"}
		json.NewEncoder(w).Encode(res)
	}

}

func (a *App) createLoungeLogin(w http.ResponseWriter, r *http.Request) {
	var u loginlounge
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.createLoungeLogin(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u.Loungeid)
}

//get lounge booking details

func (a *App) getloungebooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketid := vars["ticket_id"]

	u := loungebooking{BookingID: ticketid}
	if err := u.getloungebooking(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Ticket ID not found in our database")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getupcomingloungebookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loungeID := vars["lounge_id"]
	u := loungebooking{LoungeID: loungeID}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	loungebookings, err := u.getupcomingloungebookings(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, loungebookings)
}

func (a *App) getcurrentloungebookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loungeID := vars["lounge_id"]
	u := loungebooking{LoungeID: loungeID}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	loungebookings, err := u.getcurrentloungebookings(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, loungebookings)
}

func (a *App) createLoungeBooking(w http.ResponseWriter, r *http.Request) {
	var u loungebooking
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.createLoungeBooking(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u.BookingID)
}

func (a *App) getLoungeDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loungeid := vars["lounge_id"]

	u := loungedetail{LoungeID: loungeid}
	if err := u.getLoungeDetail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Lounge Detail not found.")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) checkin(w http.ResponseWriter, r *http.Request) {
	var u checkinout
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.checkin(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "checkin successful"})
}

func (a *App) checkout(w http.ResponseWriter, r *http.Request) {
	var u checkinout
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.checkout(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "checkout successful"})
}

func (a *App) getLoungeDetails(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	loungedetails, err := getLoungeDetails(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, loungedetails)
}

func (a *App) getloungebookingsbyffn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ffn := vars["ffn"]

	u := loungebooking{FFN: ffn}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	loungebookingsbyffn, err := u.getloungebookingsbyffn(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, loungebookingsbyffn)
}

func (a *App) getcardetails(w http.ResponseWriter, r *http.Request) {

	var u cardcheck
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.getcardetails(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getpnr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ffn := vars["ffn"]
	u := flightbooking{FFN: ffn}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	flightbookings, err := u.getpnr(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, flightbookings)
}
