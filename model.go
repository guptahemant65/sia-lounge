// model.go
package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type passenger struct {
	FFN         int    `json:"ffn"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
	Mobile      string `json:"mobile"`
	TierStatus  string `json:"tier_status"`
	Pass        string `json:"pass"`
}

type dblogin struct {
	FFN  string `json:"ffn"`
	Pass []byte `json:"pass"`
}

type reqlogin struct {
	FFN  string `json:"ffn"`
	Pass string `json:"pass"`
}

type loginlounge struct {
	Loungeid string `json:"lounge_id"`
	Pass     string `json:"pass"`
}

type dbloginlounge struct {
	Loungeid string `json:"lounge_id"`
	Pass     []byte `json:"pass"`
}

type reqloginlounge struct {
	Loungeid string `json:"lounge_id"`
	Pass     string `json:"pass"`
}

func (u *passenger) getUser(db *sql.DB) error {

	statement := fmt.Sprintf("SELECT email,name,country_code,mobile,tier_status FROM passenger_details WHERE ffn=%d", u.FFN)
	return db.QueryRow(statement).Scan(&u.Email, &u.Name, &u.CountryCode, &u.Mobile, &u.TierStatus)
}

func (u *dbloginlounge) getLoungeLogin(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT pass FROM lounge_login WHERE lounge_id='%s'", u.Loungeid)
	return db.QueryRow(statement).Scan(&u.Pass)

}

func (u *dblogin) getUserLogin(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT pass FROM passenger_details WHERE ffn='%s'", u.FFN)
	return db.QueryRow(statement).Scan(&u.Pass)
}

func (u *passenger) updateUser(db *sql.DB) error {

	statement := fmt.Sprintf("UPDATE passenger_details SET name='%s', email='%s', country_code='%s', mobile='%s',tier_status='%s',pass= '%s' WHERE ffn=%d", u.Name, u.Email, u.CountryCode, u.Mobile, u.TierStatus, u.Pass, u.FFN)
	_, err := db.Exec(statement)
	return err
}
func (u *passenger) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM passenger_details WHERE ffn=%d", u.FFN)
	_, err := db.Exec(statement)
	return err
}
func (u *passenger) createUser(db *sql.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	u.Pass = string(hash)
	statement := fmt.Sprintf("INSERT INTO passenger_details(name,email,country_code,mobile,tier_status,pass) VALUES('%s','%s','%s','%s','%s','%s')", u.Name, u.Email, u.CountryCode, u.Mobile, u.TierStatus, u.Pass)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.FFN)
	if err != nil {
		return err
	}
	return nil
}
func getUsers(db *sql.DB, start, count int) ([]passenger, error) {
	statement := fmt.Sprintf("SELECT ffn,name,email,country_code,mobile,tier_status FROM passenger_details LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []passenger{}
	for rows.Next() {
		var u passenger
		if err := rows.Scan(&u.FFN, &u.Name, &u.Email, &u.CountryCode, &u.Mobile, &u.TierStatus); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (u *loginlounge) createLoungeLogin(db *sql.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	u.Pass = string(hash)
	statement := fmt.Sprintf("INSERT INTO lounge_login(pass) VALUES('%s')", u.Pass)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.Loungeid)
	if err != nil {
		return err
	}
	return nil
}

// GET LOUNGE BOOKING DETAILS

type loungebooking struct {
	BookingID     string `json:"ticket_id"`
	FFN           string `json:"ffn"`
	LoungeID      string `json:"lounge_id"`
	LoungeName    string `json:"lounge_name"`
	LoungeAddress string `json:"lounge_address"`
	Num           string `json:"no_of_guests"`
	Names         string `json:"guest_name"`
	Checkin       string `json:"checkin"`
	Checkout      string `json:"checkout"`
	PNR           string `json:"pnr"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	AmountPaid    string `json:"amount_paid"`
}

func (u *loungebooking) getloungebooking(db *sql.DB) error {

	statement := fmt.Sprintf("SELECT ffn,lounge_id,lounge_name,lounge_address,no_of_guests,guest_name,checkin,checkout,pnr,status,payment_method,amount_paid FROM lounge_booking WHERE ticket_id='%s'", u.BookingID)
	return db.QueryRow(statement).Scan(&u.FFN, &u.LoungeID, &u.LoungeName, &u.LoungeAddress, &u.Num, &u.Names, &u.Checkin, &u.Checkout, &u.PNR, &u.Status, &u.PaymentMethod, &u.AmountPaid)
}

func (u *loungebooking) getupcomingloungebookings(db *sql.DB, start, count int) ([]loungebooking, error) {
	statement := fmt.Sprintf("SELECT ticket_id,ffn,lounge_name,lounge_address,no_of_guests,guest_name,checkin,checkout,pnr,status,payment_method,amount_paid FROM lounge_booking where lounge_id = '%s' && status = 'confirmed' && TIMESTAMPDIFF(HOUR,checkin,CONVERT_TZ( current_timestamp(),'GMT','+08:00' ))<=12 ", u.LoungeID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loungebookings := []loungebooking{}
	for rows.Next() {
		var u loungebooking
		if err := rows.Scan(&u.BookingID, &u.FFN, &u.LoungeName, &u.LoungeAddress, &u.Num, &u.Names, &u.Checkin, &u.Checkout, &u.PNR, &u.Status, &u.PaymentMethod, &u.AmountPaid); err != nil {
			return nil, err
		}
		loungebookings = append(loungebookings, u)
	}
	return loungebookings, nil
}

func (u *loungebooking) getcurrentloungebookings(db *sql.DB, start, count int) ([]loungebooking, error) {
	statement := fmt.Sprintf("SELECT ffn,no_of_guests,guest_name,checkin,pnr,payment_method,amount_paid FROM lounge_booking where lounge_id = '%s' && status = 'in progress' ", u.LoungeID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loungebookings := []loungebooking{}
	for rows.Next() {
		var u loungebooking
		if err := rows.Scan(&u.FFN, &u.Num, &u.Names, &u.Checkin, &u.PNR, &u.PaymentMethod, &u.AmountPaid); err != nil {
			return nil, err
		}
		loungebookings = append(loungebookings, u)
	}
	return loungebookings, nil
}

func (u *loungebooking) getloungebookingsbyffn(db *sql.DB, start, count int) ([]loungebooking, error) {
	statement := fmt.Sprintf("SELECT ticket_id,ffn,lounge_id,lounge_name,lounge_address,no_of_guests,guest_name,checkin,checkout,pnr,status,payment_method,amount_paid FROM lounge_booking where status != 'completed' && ffn='%s' order by checkin", u.FFN)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loungebookingsbyffn := []loungebooking{}
	for rows.Next() {
		var u loungebooking
		if err := rows.Scan(&u.BookingID, &u.FFN, &u.LoungeID, &u.LoungeName, &u.LoungeAddress, &u.Num, &u.Names, &u.Checkin, &u.Checkout, &u.PNR, &u.Status, &u.PaymentMethod, &u.AmountPaid); err != nil {
			return nil, err
		}
		loungebookingsbyffn = append(loungebookingsbyffn, u)
	}
	return loungebookingsbyffn, nil
}

func (u *loungebooking) createLoungeBooking(db *sql.DB) error {
	u.Status = "confirmed"
	rand.Seed(time.Now().UnixNano())
	u.BookingID = randSeq(25)
	loungedetailstmt := fmt.Sprintf("SELECT lounge_name,location from lounge_details where lounge_id = '%s'", u.LoungeID)
	db.QueryRow(loungedetailstmt).Scan(&u.LoungeName, &u.LoungeAddress)
	statement := fmt.Sprintf("INSERT INTO lounge_booking(ticket_id,ffn,lounge_id,lounge_name,lounge_address,no_of_guests,guest_name,checkin,checkout,pnr,status,payment_method,amount_paid) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", u.BookingID, u.FFN, u.LoungeID, u.LoungeName, u.LoungeAddress, u.Num, u.Names, u.Checkin, u.Checkout, u.PNR, u.Status, u.PaymentMethod, u.AmountPaid)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

type loungedetail struct {
	LoungeID            string `json:"lounge_id"`
	LoungeName          string `json:"lounge_name"`
	Amenities           string `json:"amenities"`
	Price               string `json:"price"`
	AcceptedCards       string `json:"accepted_cards"`
	PrivateRoomCapacity string `json:"private_room_capacity"`
	SofaCapacity        string `json:"sofa_capacity"`
	Location            string `json:"location"`
}

func (u *loungedetail) getLoungeDetail(db *sql.DB) error {

	statement := fmt.Sprintf("SELECT lounge_name,amenities,price,accepted_cards,private_room_capacity,sofa_capacity,location FROM lounge_details WHERE lounge_id='%s'", u.LoungeID)
	return db.QueryRow(statement).Scan(&u.LoungeName, &u.Amenities, &u.Price, &u.AcceptedCards, &u.PrivateRoomCapacity, &u.SofaCapacity, &u.Location)
}

func getLoungeDetails(db *sql.DB, start, count int) ([]loungedetail, error) {
	statement := fmt.Sprintf("SELECT lounge_id,lounge_name,amenities,price,accepted_cards,private_room_capacity,sofa_capacity,location FROM lounge_details")
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loungedetails := []loungedetail{}
	for rows.Next() {
		var u loungedetail
		if err := rows.Scan(&u.LoungeID, &u.LoungeName, &u.Amenities, &u.Price, &u.AcceptedCards, &u.PrivateRoomCapacity, &u.SofaCapacity, &u.Location); err != nil {
			return nil, err
		}
		loungedetails = append(loungedetails, u)
	}
	return loungedetails, nil
}

type checkinout struct {
	TicketID string `json:"ticket_id"`
}

func (u *checkinout) checkin(db *sql.DB) error {
	var statusin = "IN PROGRESS"
	statement := fmt.Sprintf("UPDATE lounge_booking SET status = '%s' WHERE ticket_id = '%s'", statusin, u.TicketID)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (u *checkinout) checkout(db *sql.DB) error {
	var statusout = "COMPLETED"
	statement := fmt.Sprintf("UPDATE lounge_booking SET status = '%s' WHERE ticket_id = '%s'", statusout, u.TicketID)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

type cardcheck struct {
	CardNumber      string `json:"card_number"`
	LoungeLeft      string `json:"lounge_left"`
	AvailableLounge string `json:"available_lounge"`
}

func (u *cardcheck) getcardetails(db *sql.DB) error {

	statement := fmt.Sprintf("SELECT lounge_left,available_lounge FROM card_details WHERE card_number='%s'", u.CardNumber)
	return db.QueryRow(statement).Scan(&u.LoungeLeft, &u.AvailableLounge)
}

type flightbooking struct {
	PNR          string `json:"pnr"`
	FFN          string `json:"ffn"`
	From         string `json:"dep_from"`
	To           string `json:"arr_to"`
	Time         string `json:"time"`
	ExpectedTime string `json:"updated_time"`
	Names        string `json:"names"`
	FlightCode   string `json:"flight_code"`
	Terminal     string `json:"terminal"`
	Gate         string `json:"gate"`
}

func (u *flightbooking) getpnr(db *sql.DB, start, count int) ([]flightbooking, error) {
	statement := fmt.Sprintf("SELECT pnr,dep_from,arr_to,time,updated_time,names,flight_code,terminal,gate FROM booking_table where ffn = '%s'", u.FFN)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	flightbookings := []flightbooking{}
	for rows.Next() {
		var u flightbooking
		if err := rows.Scan(&u.PNR, &u.From, &u.To, &u.Time, &u.ExpectedTime, &u.Names, &u.FlightCode, &u.Terminal, &u.Gate); err != nil {
			return nil, err
		}
		flightbookings = append(flightbookings, u)
	}
	return flightbookings, nil
}
