package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Developer is a struct that represents a Developer in our application
type Developer struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	Creation int    `json:"creation"`
}

// Publisher is a struct that represents a Publisher in our application
type Publisher struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// Director is a struct that represents a Director in our application
type Director struct {
	FullName    string `json:"fullName"`
	Nationality string `json:"nationality"`
	Age         int    `json:"age"`
}

// Game is a struct that represents a game in our application
type Game struct {
	Title     string    `json:"title"`
	Director  Director  `json:"director"`
	Developer Developer `json:"developer"`
	Publisher Publisher `json:"publisher"`
}

var games []Game = []Game{}

func main() {
	jsonFile, err := os.Open("games.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened games.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	_ = json.Unmarshal([]byte(byteValue), &games)

	router := mux.NewRouter()

	router.HandleFunc("/games", addItem).Methods("POST")

	router.HandleFunc("/games", getAllGames).Methods("GET")

	router.HandleFunc("/games/{id}", getGame).Methods("GET")

	router.HandleFunc("/games/{id}", updateGame).Methods("PUT")

	router.HandleFunc("/games/{id}", patchGame).Methods("PATCH")

	router.HandleFunc("/games/{id}", deleteGame).Methods("DELETE")

	http.ListenAndServe(":5000", router)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// there was an error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	// error checking
	if id >= len(games) {
		w.WriteHeader(404)
		w.Write([]byte("No game found with specified ID"))
		return
	}

	game := games[id]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var newGame Game
	json.NewDecoder(r.Body).Decode(&newGame)

	games = append(games, newGame)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(games)

	gameUpdates, _ := json.Marshal(games)

	_ = ioutil.WriteFile("games.json", gameUpdates, 0644)
}

func updateGame(w http.ResponseWriter, r *http.Request) {
	// get the ID of the game from the route parameters
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	// error checking
	if id >= len(games) {
		w.WriteHeader(404)
		w.Write([]byte("No game found with specified ID"))
		return
	}

	// get the value from JSON body
	var updatedGame Game
	json.NewDecoder(r.Body).Decode(&updatedGame)

	games[id] = updatedGame

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGame)

	gameUpdates, _ := json.Marshal(games)

	_ = ioutil.WriteFile("games.json", gameUpdates, 0644)
}

func patchGame(w http.ResponseWriter, r *http.Request) {
	// get the ID of the game from the route parameters
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	// error checking
	if id >= len(games) {
		w.WriteHeader(404)
		w.Write([]byte("No game found with specified ID"))
		return
	}

	// get the current value
	game := &games[id]
	json.NewDecoder(r.Body).Decode(game)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)

	gameUpdates, _ := json.Marshal(games)

	_ = ioutil.WriteFile("games.json", gameUpdates, 0644)
}

func deleteGame(w http.ResponseWriter, r *http.Request) {
	// get the ID of the game from the route parameters
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	// error checking
	if id >= len(games) {
		w.WriteHeader(404)
		w.Write([]byte("No game found with specified ID"))
		return
	}

	games = append(games[:id], games[id+1:]...)
	gameUpdates, _ := json.Marshal(games)

	_ = ioutil.WriteFile("games.json", gameUpdates, 0644)

	w.WriteHeader(200)
}
