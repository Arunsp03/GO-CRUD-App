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
	db.SetMaxOpenConns(1)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		owner TEXT,
		taskId TEXT PRIMARY KEY,
		taskName TEXT,
		createdDate TEXT DEFAULT (datetime('now'))
	)
`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
    userName TEXT PRIMARY KEY,
    password TEXT
)
`)
	if err != nil {
		log.Fatal(err)
	}

	taskDb := TaskDb{db: db}
	userDb := UserDB{db: db}

	router.Route("/api/tasks", func(r chi.Router) {
		r.Use(AuthMiddleWare)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value("username").(string)

			tasksData, err := taskDb.GetTasks(username)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, tasksData)
		})

		r.Post("/addtask", func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value("username").(string)
			var newTask Task
			err := json.NewDecoder(r.Body).Decode(&newTask)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			err = taskDb.AddTask(newTask, username)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}

			respondWithJSON(w, http.StatusOK, newTask)
		})

		r.Put("/edittask", func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value("username").(string)
			var task Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			err = taskDb.EditTask(task, username)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondWithJSON(w, http.StatusOK, task)
		})

		r.Delete("/deletetask/{taskId}", func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value("username").(string)
			var taskId string
			taskId = chi.URLParam(r, "taskId")
			err := taskDb.DeleteTask(taskId, username)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, map[string]string{
				"message": "Successfully deleted",
				"taskId":  taskId,
			})

		})

	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		respondWithJSON(w, http.StatusOK, "Hello world")
	})

	router.Post("/registeruser", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		err = userDb.RegisterUser(user.UserName, user.Password)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, "Successfullly registered user")

	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		isValidUser, err := userDb.ValidateLogin(user.UserName, user.Password)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if isValidUser {
			tokenString, err := CreateToken(user.UserName)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged in", "token": tokenString})
			return
		}

		respondWithJSON(w, http.StatusUnauthorized, "Login failed . Wrong password")

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
