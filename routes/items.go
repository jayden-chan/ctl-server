package routes

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"

	"github.com/jayden-chan/ctl-server/db"
	"github.com/jayden-chan/ctl-server/util"
)

// URI: /items
func Items(res http.ResponseWriter, req *http.Request) {
	authSuccess, user, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		query := `SELECT id, folder_id, status, description, due FROM items WHERE user_id = $1`
		rows, err := db.Query(query, user)
		if err != nil {
			util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type row struct {
			ID          string            `json:"id"`
			FolderID    db.JSONNullString `json:"folder"`
			Status      string            `json:"status"`
			Description string            `json:"description"`
			Due         db.JSONNullString `json:"due"`
		}

		type results struct {
			Results []row `json:"items"`
		}

		var ret results
		for rows.Next() {
			var r row
			err := rows.Scan(&r.ID, &r.FolderID, &r.Status, &r.Description, &r.Due)
			if err != nil {
				log.Println(err)
			}
			ret.Results = append(ret.Results, r)
		}

		util.HTTPJSONRes(res, ret, http.StatusOK)
		return

	case http.MethodPost:
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			util.HTTPRes(res, "Malformed request data", http.StatusBadRequest)
			return
		}

		var (
			folder      db.JSONNullString
			status      db.JSONNullString
			description string
			due         db.JSONNullString
		)
		paths := [][]string{
			[]string{"folder"},
			[]string{"status"},
			[]string{"description"},
			[]string{"due"},
		}
		jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			switch idx {
			case 0:
				folder.String = string(value)
				folder.Valid = true
			case 1:
				status.String = string(value)
				status.Valid = true
			case 2:
				description = string(value)
			case 3:
				due.String = string(value)
				due.Valid = true
			}
		}, paths...)

		if description == "" {
			util.HTTPRes(res, "One or more fields missing in request body", http.StatusBadRequest)
			return
		}

		if !status.Valid {
			status.String = "pending"
			status.Valid = true
		}

		query := `INSERT INTO items(user_id, folder_id, status, description, due) VALUES ($1, $2, $3, $4, $5)`
		_, err = db.Exec(query, user, folder, status, description, due)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error has occurred", http.StatusInternalServerError)
			return
		}

		util.HTTPRes(res, "Item created", http.StatusCreated)
		return
	}
}
