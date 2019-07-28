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
	BookingID string `json:"ticket_id"`
	FFN       string `json:"ffn"`
	Num       string `json:"no_of_guests"`
	Names     string `json:"guest_names"`
	Checkin   string `json:"checkin"`
	Checkout  string `json:"checkout"`
	PNR       string `json:"pnr"`
	Status    string `json:"status"`
}

func (u *loungebooking) getloungebooking(db *sql.DB) error {

	statement := fmt.Sprintf("SELECT ffn,no_of_guests,guest_names,checkin,checkout,pnr,status FROM lounge_booking WHERE ticket_id='%s'", u.BookingID)
	return db.QueryRow(statement).Scan(&u.FFN, &u.Num, &u.Names, &u.Checkin, &u.Checkout, &u.PNR, &u.Status)
}

func getloungebookings(db *sql.DB, start, count int) ([]loungebooking, error) {
	statement := fmt.Sprintf("SELECT ticket_id,ffn,no_of_guests,guest_names,checkin,checkout,pnr,status FROM lounge_booking where status != 'completed' && TIMESTAMPDIFF(HOUR,checkin,CONVERT_TZ( current_timestamp(),'GMT','+08:00' ))<=12 ")
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loungebookings := []loungebooking{}
	for rows.Next() {
		var u loungebooking
		if err := rows.Scan(&u.BookingID, &u.FFN, &u.Num, &u.Names, &u.Checkin, &u.Checkout, &u.PNR, &u.Status); err != nil {
			return nil, err
		}
		loungebookings = append(loungebookings, u)
	}
	return loungebookings, nil
}

func (u *loungebooking) createLoungeBooking(db *sql.DB) error {
	rand.Seed(time.Now().UnixNano())
	u.BookingID = randSeq(25)
	statement := fmt.Sprintf("INSERT INTO lounge_booking(ticket_id,ffn,no_of_guests,guest_names,checkin,checkout,pnr,status) VALUES('%s','%s','%s','%s','%s','%s','%s','%s')", u.BookingID, u.FFN, u.Num, u.Names, u.Checkin, u.Checkout, u.PNR, u.Status)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}
