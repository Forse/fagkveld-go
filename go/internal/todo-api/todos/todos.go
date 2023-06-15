package todos

import (
	"encoding/json"
	tododb "fagkveld/internal/todo-api/todos/db"
	"net/http"
)

type Todo struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"isCompleted"`
}

func CreateTodo(db *tododb.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		var todo Todo

		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todoEntity := tododb.Todo{
			ID:          0,
			Title:       todo.Title,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
		}

		db.Create(&todoEntity)

		todo.ID = todoEntity.ID

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
}
