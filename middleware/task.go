package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ht21992/go-task-manager/database"
)

func TaskMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var payload map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if _, ok := payload["completed"]; !ok {
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
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
		fmt.Printf("Passed")
		if count > 0 {
			http.Error(w, "You have pending subtasks", http.StatusForbidden)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		fmt.Println(r.Body)
		next.ServeHTTP(w, r)
	})
}
