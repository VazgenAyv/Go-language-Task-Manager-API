package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ht21992/go-task-manager/database"
)

func TaskMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		}
		defer r.Body.Close()

		if _, ok := payload["completed"]; !ok {
			next.ServeHTTP(w, r)
			return
		}

		vars := mux.Vars(r)
		idStr := vars["id"]

		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var count int
		err = database.DB.QueryRow("SELECT count(id) FROM tasks WHERE maintask = ? AND completed = 0;", taskID).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if count > 0 {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}