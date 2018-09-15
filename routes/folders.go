package routes

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"

	"github.com/jayden-chan/ctl-server/db"
	"github.com/jayden-chan/ctl-server/util"
)

// URI: /folders
// Folders returns a list of the user's folders or adds a new folder
func Folders(res http.ResponseWriter, req *http.Request) {
	authSuccess, user, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		query := `SELECT id, name, subfolder FROM folders WHERE user_id = $1`
		rows, err := db.Query(query, user)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error has occurred", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type row struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Subfolder string `json:"subfolder"`
		}

		type results struct {
			Results []row `json:"folders"`
		}

		var ret results
		for rows.Next() {
			var r row
			rows.Scan(&r.ID, &r.Name, &r.Subfolder)
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
			name   string
			parent sql.NullString
		)

		paths := [][]string{
			[]string{"name"},
			[]string{"parent"},
		}
		jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			switch idx {
			case 0:
				name = string(value)
			case 1:
				asString := string(value)
				if asString != "" {
					parent.String = asString
					parent.Valid = true
				}
			}
		}, paths...)

		if name == "" {
			util.HTTPRes(res, "'Name' field not found in request body", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO folders(user_id, parent, name) VALUES($1, $2, $3)`
		_, err = db.Exec(query, user, parent, name)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error occurred", http.StatusInternalServerError)
			return
		}

		util.HTTPRes(res, "Folder added", http.StatusCreated)
		return
	}
}

// URI: /folders/:id
func FoldersID(res http.ResponseWriter, req *http.Request) {
	authSuccess, user, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}
	folderID := mux.Vars(req)["folderID"]
	if folderID == "" {
		util.HTTPRes(res, "'Folder ID' field not found in request URI", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodDelete:
		query := `DELETE from folders WHERE user_id = $1 AND id = $2`
		results, err := db.Exec(query, user, folderID)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error has occurred", http.StatusInternalServerError)
			return
		}

		if r, _ := results.RowsAffected(); r <= 0 {
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
			id     string
			name   string
			parent sql.NullString
		)

		paths := [][]string{
			[]string{"id"},
			[]string{"name"},
			[]string{"parent"},
		}

		jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			switch idx {
			case 0:
				id = string(value)
			case 1:
				name = string(value)
			case 2:
				asString := string(value)
				if asString != "" {
					parent.String = asString
					parent.Valid = true
				}
			}
		}, paths...)

		if id == "" {
			util.HTTPRes(res, "'ID' field not found in request body", http.StatusBadRequest)
			return
		}

		var query string
		if name != "" {
			query = `UPDATE folders SET name = $1, parent = $2 WHERE user_id = $3 AND id = $4`
		} else {
			query = `UPDATE folders parent = $2 WHERE user_id = $3 AND id = $4`

		}

		results, err := db.Exec(query, name, parent, user, id)
		if err != nil {
			log.Println(err)
			util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
			return
		}

		if r, _ := results.RowsAffected(); r == 0 {
			util.HTTPRes(res, "Folder does not exist or does not belong to user", http.StatusNotFound)
			return
		}

		util.HTTPRes(res, "Folder updated", http.StatusOK)
		return
	}
}
