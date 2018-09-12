package routes

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/jayden-chan/robotender-server/db"
	"github.com/jayden-chan/robotender-server/util"
)

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
		case 0: // email
			email = string(value)
		case 1: // password
			password = string(value)
		}
	}, paths...)

	if email == "" || password == "" {
		util.HTTPRes(res, "One or more fields missing", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred. Please try again.", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	if rows.Next() {
		util.HTTPRes(res, "Email is already registered", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users(email, password, access) VALUES($1, crypt($2, gen_salt('bf', 8)), $3)", email, password, "normal")
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred. Please try again.", http.StatusInternalServerError)
		return
	}
}

// URI: /deregister
func Deregister(res http.ResponseWriter, req *http.Request) {
	authSuccess, user := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, user)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
		return
	}

	util.HTTPRes(res, "Customer deleted.", http.StatusOK)
	return
}

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

	rows, err := db.Query("SELECT id FROM customers WHERE email = $1 AND password = crypt($2, password)", email, password)

	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred, please try again later", http.StatusInternalServerError)
		return
	}

	results := 0
	userID := ""

	defer rows.Close()
	for rows.Next() {
		results++
		var (
			id string
		)
		rows.Scan(&id)
		userID = id
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

	token, err := util.GenerateJWT(userID)
	if err != nil {
		log.Println(err)
		util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
		return
	}
	util.HTTPJSONRes(res, JWTRes{Token: token}, http.StatusOK)
}
