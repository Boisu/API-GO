package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type Word struct {
	ID      int    `form:"id" json:"id"`
	Text    string `form:"text" json:"text"`
	Counter int64  `form:"counter" json:"counter"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Main Go
func main() {
	db, err = gorm.Open("mysql", "root:@/boisu?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connection failder", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Word{})
	handleRequests()
}

// Routing
func handleRequests() {
	log.Println("Start the developtment server at http://127.0.0.1:8081")
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/words", allWords).Methods("GET")
	myRouter.HandleFunc("/words/new", createWords).Methods("POST")
	myRouter.HandleFunc("/words/{id}", detailWords).Methods("GET")
	myRouter.HandleFunc("/words/{id}", updateWords).Methods("PUT")
	myRouter.HandleFunc("/words/{id}", deleteWords).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

// Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

// All Data
func allWords(w http.ResponseWriter, r *http.Request) {

	word := []Word{}
	db.Find(&word)

	res := Result{Code: 200, Data: word, Message: "Success get all data"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

// Create Data
func createWords(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var word Word
	json.Unmarshal(payloads, &word)

	db.Create(&word)

	res := Result{Code: 200, Data: word, Message: "Data Created"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Find Data
func detailWords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get words")

	vars := mux.Vars(r)
	wordID := vars["id"]

	var word Word
	db.First(&word, wordID)

	if result := db.First(&word, wordID); result.Error != nil {
		// Not Found
		res := Result{Code: 204, Message: "No Data Found"}
		results, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(results)
	} else {
		// Found
		res := Result{Code: 200, Data: word, Message: "Data Founded"}
		results, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(results)
	}
}

// Update Data
func updateWords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wordID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var wordNew Word
	json.Unmarshal(payloads, &wordNew)

	var word Word
	db.First(&word, wordID)

	if result := db.First(&word, wordID); result.Error != nil {
		// Not Found
		res := Result{Code: 204, Message: "Failed update words"}
		results, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(results)
	} else {
		// Found
		db.Model(&word).Updates(wordNew)

		res := Result{Code: 200, Data: word, Message: "Data Updated"}
		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

// Delete Data
func deleteWords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wordID := vars["id"]

	var word Word

	db.First(&word, wordID)

	if result := db.First(&word, wordID); result.Error != nil {
		// Not Found
		res := Result{Code: 204, Message: "No Data Found"}
		results, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(results)
	} else {
		// Found
		db.Delete(&word)

		res := Result{Code: 200, Message: "Data Deleted"}
		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
