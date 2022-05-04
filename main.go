package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

var db *gorm.DB
var err error

type Todo struct {
	TodoId    uint      `gorm:"primaryKey" json:"todo_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateTodoDTO struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UpdateTodoDTO struct {
	Completed bool `json:"completed"`
}

// use godot package to load/read the .env file and
// return the value of the key
func dotenvGetenv(key string) string {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	value := os.Getenv(key)
	return value
}

func GetAllTodos(w http.ResponseWriter, _ *http.Request) {
	var todos []Todo
	result := db.Find(&todos)
	if result.Error != nil {
		log.Println("Error when querying for todos")
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Message    string
		Success    bool
		StatusCode uint
		Data       []Todo
	}{
		Message:    "Todos fetched successfully",
		Success:    true,
		StatusCode: 201,
		Data:       todos,
	})
}

func GetOneTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoId := params["todoId"]

	var todo Todo
	db.First(&todo, todoId)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Message    string
		Success    bool
		StatusCode uint
		Data       Todo
	}{
		Message:    "Todo fetched successfully",
		Success:    true,
		StatusCode: 201,
		Data:       todo,
	})

}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t CreateTodoDTO
	err = decoder.Decode(&t)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	todo := Todo{
		Title:     t.Title,
		Completed: t.Completed,
	}
	db.Create(&todo)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Message    string
		Success    bool
		StatusCode uint
	}{
		Message:    "Todo created successfully",
		Success:    true,
		StatusCode: 201,
	})
	if err != nil {
		log.Fatalf("Cannot encode to json")
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoId := params["todoId"]
	var t UpdateTodoDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		log.Fatalf("Error decoding the update todo %v", err)
	}
	var todo Todo
	db.First(&todo, todoId)
	todo.Completed = t.Completed
	db.Save(&todo)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Message    string
		Success    bool
		StatusCode uint
	}{
		Message:    "Todo updated successfully",
		Success:    true,
		StatusCode: 201,
	})
	if err != nil {
		log.Fatalf("Cannot encode to json")
	}
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	todoId := params["todoId"]

	db.Delete(&Todo{}, todoId)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Message    string
		Success    bool
		StatusCode uint
	}{
		Message:    "Todo deleted successfully",
		Success:    true,
		StatusCode: 201,
	})
	if err != nil {
		log.Fatalf("Cannot encode to json")
	}
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, err = fmt.Fprintf(w, "Hello World!")
	if err != nil {
		log.Fatalf("Internal Server Error")
	}
}

func main() {
	// Connect to database
	dbHost := dotenvGetenv("DB_HOST")
	dbUser := dotenvGetenv("DB_USER")
	dbPassword := dotenvGetenv("DB_PASSWORD")
	dbName := dotenvGetenv("DB_NAME")
	dbPort := dotenvGetenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database")
	}
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatalf("Error migrating schema in the database")
	}

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/", HelloWorld).Methods("GET")
	r.HandleFunc("/todos", GetAllTodos).Methods("GET")
	r.HandleFunc("/todos/{todoId}", GetOneTodo).Methods("GET")
	r.HandleFunc("/todos", CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{todoId}", UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{todoId}", DeleteTodo).Methods("DELETE")

	// Server on port 8080
	port := dotenvGetenv("PORT")
	log.Printf("Listening on port %v", port)
	addr := fmt.Sprintf(":%v", port)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Internal server error")
	}
}
