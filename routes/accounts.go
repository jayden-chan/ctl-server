package routes

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/jayden-chan/ctl-server/db"
	"github.com/jayden-chan/ctl-server/util"
)

// Register registers a user in the database
// Path: /register
func Register(res http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HTTPRes(res, "Malformed request data", http.StatusBadRequest)
		return
	}

	var (
		email    string
		password string
	)

	paths := [][]string{
		[]string{"email"},
		[]string{"password"},
	}
	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		switch idx {
		case 0:
			email = string(value)
		case 1:
			password = string(value)
		}
	}, paths...)

	if email == "" || password == "" {
		util.HTTPRes(res, "One or more fields missing", http.StatusBadRequest)
		return
	}

	queryEmail := `SELECT * FROM users WHERE email = $1`
	rows, err := db.Query(queryEmail, email)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error has occurred", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	if rows.Next() {
		util.HTTPRes(res, "Email is already registered", http.StatusBadRequest)
		return
	}

	queryInsert := `INSERT INTO users(email, password, access) VALUES($1, crypt($2, gen_salt('bf', 8)), $3)`
	_, err = db.Exec(queryInsert, email, password, "normal")
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error has occurred", http.StatusInternalServerError)
		return
	}
}

// Deregister deletes a user from the database
// URI: /deregister
func Deregister(res http.ResponseWriter, req *http.Request) {
	authSuccess, user, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "User authorization failed.", http.StatusUnauthorized)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, user)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
		return
	}

	util.HTTPRes(res, "User deleted.", http.StatusOK)
	return
}

// Login verifies a user's credentials and issues a JWT auth token
// Path: /login
func Login(res http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HTTPRes(res, "Malformed request data", http.StatusBadRequest)
		return
	}

	var (
		email    string
		password string
	)

	paths := [][]string{
		[]string{"email"},
		[]string{"password"},
	}
	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		switch idx {
		case 0:
			email = string(value)
		case 1:
			password = string(value)
		}
	}, paths...)

	if email == "" || password == "" {
		util.HTTPRes(res, "One or more fields missing", http.StatusBadRequest)
		return
	}

	query := `SELECT id, access FROM users WHERE email = $1 AND password = crypt($2, password)`
	rows, err := db.Query(query, email, password)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred, please try again later", http.StatusInternalServerError)
		return
	}

	results := 0
	userID := ""
	userAccess := ""

	defer rows.Close()
	for rows.Next() {
		results++
		rows.Scan(&userID, &userAccess)
	}

	if results < 1 {
		util.HTTPRes(res, "Incorrect email or password.", http.StatusUnauthorized)
		return
	} else if results > 1 {
		util.HTTPRes(res, "Internal server error has occurred.", http.StatusInternalServerError)
		return
	}

	type JWTRes struct {
		Token string `json:"token"`
	}

	token, err := util.GenerateJWT(userID, userAccess)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
		return
	}
	util.HTTPJSONRes(res, JWTRes{Token: token}, http.StatusOK)
}
