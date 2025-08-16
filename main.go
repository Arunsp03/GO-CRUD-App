package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port not found in .env file")
	}
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	fmt.Printf("Server is running on Port %v", port)

	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		taskId TEXT PRIMARY KEY,
		taskName TEXT,
		createdDate TEXT DEFAULT (datetime('now'))
	)
`)
	if err != nil {
		log.Fatal(err)
	}
	taskDb := TaskDb{db: db}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		respondWithJSON(w, http.StatusOK, "Hello world")
	})

	router.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {

		tasksData, err := taskDb.GetTasks()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, tasksData)
	})

	router.Post("/addtask", func(w http.ResponseWriter, r *http.Request) {
		var newTask Task
		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		err = taskDb.AddTask(newTask)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		respondWithJSON(w, http.StatusOK, newTask)
	})

	router.Put("/edittask", func(w http.ResponseWriter, r *http.Request) {
		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		err = taskDb.EditTask(task)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, task)
	})

	router.Delete("/deletetask/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		var taskId string
		taskId = chi.URLParam(r, "taskId")
		err := taskDb.DeleteTask(taskId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Successfully deleted",
			"taskId":  taskId,
		})

	})

	http.ListenAndServe(":"+port, router)

}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
