package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/jayden-chan/ctl-server/db"
	"github.com/jayden-chan/ctl-server/util"
)

// URI: /folders
func Folders(res http.ResponseWriter, req *http.Request) {
	authSuccess, user, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, subfolder FROM folders WHERE user_id = $1", user)
		if err != nil {
			util.HTTPRes(res, "An internal server error occurred.", http.StatusInternalServerError)
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
			name      string
			subfolder string
		)

		paths := [][]string{
			[]string{"name"},
			[]string{"subfolder"},
		}
		jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			switch idx {
			case 0:
				name = string(value)
			case 1:
				subfolder = string(value)
			}
		}, paths...)

		if name == "" {
			util.HTTPRes(res, "'Name' field not found in request body", http.StatusBadRequest)
			return
		}

		if subfolder == "" {
			query := `INSERT INTO folders(user_id, name) VALUES($1, $2)`
			_, err := db.Exec(query, user, name)
			if err != nil {
				util.HTTPRes(res, "An internal server error occurred", http.StatusInternalServerError)
				return
			}

		} else {
			query := `INSERT INTO folders(user_id, subfolder, name) VALUES($1, $2, $3)`
			_, err := db.Exec(query, user, subfolder, name)
			if err != nil {
				util.HTTPRes(res, "An internal server error occurred", http.StatusInternalServerError)
				return
			}
		}
	}
}

// URI: /folders/:id
func FoldersID(res http.ResponseWriter, req *http.Request) {
	authSuccess, _, _ := util.Authenticate(req)
	if !authSuccess {
		util.HTTPRes(res, "Customer authorization failed.", http.StatusUnauthorized)
		return
	}

	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HTTPRes(res, "Malformed request data", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodDelete:
		util.HTTPRes(res, "Not implemented", http.StatusNotImplemented)
		return
	}
}
