package routes

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"

	"github.com/jayden-chan/ctl-server/db"
	"github.com/jayden-chan/ctl-server/util"
)

// Items gets the items for a user or adds a new item
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

// ItemsID deletes or updates the given item
// URI: /items/:id
func ItemsID(res http.ResponseWriter, req *http.Request) {
	authSuccess, user, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}

	itemID := mux.Vars(req)["itemID"]
	if itemID == "" {
		util.HTTPRes(res, "'Item ID' field not found in request URI", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodDelete:
		query := `DELETE FROM items WHERE user_id = $1 AND id = $1`
		result, err := db.Exec(query, user, itemID)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error has occurred", http.StatusInternalServerError)
			return
		}

		if r, _ := result.RowsAffected(); r == 0 {
			util.HTTPRes(res, "Item not found or does not belong to user", http.StatusNotFound)
			return
		}

		util.HTTPRes(res, "Item deleted", http.StatusOK)
		return

	case http.MethodPatch:
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			util.HTTPRes(res, "Malformed request data", http.StatusBadRequest)
			return
		}

		var (
			folder      db.JSONNullString
			status      string
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
				if vt != jsonparser.Null {
					folder.String = string(value)
					folder.Valid = true
				}
			case 1:
				status = string(value)
			case 2:
				description = string(value)
			case 3:
				if vt != jsonparser.Null {
					due.String = string(value)
					due.Valid = true
				}
			}
		}, paths...)

		if status == "" || description == "" {
			util.HTTPRes(res, "One or more fields missing from request", http.StatusBadRequest)
			return
		}

		query := `UPDATE items SET folder_id = $1, status = $2, description = $3, due = $4
		WHERE user_id = $5 AND id = $6`
		result, err := db.Exec(query, folder, status, description, due, user, itemID)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
			return
		}

		if r, _ := result.RowsAffected(); r == 0 {
			util.HTTPRes(res, "Item not found or does not belong to user", http.StatusBadRequest)
			return
		}

		util.HTTPRes(res, "Item updated", http.StatusOK)
		return
	}
}
