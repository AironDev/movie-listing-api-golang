package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func index(w http.ResponseWriter, r *http.Request) {
	movies := getMovies()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Movies retrieved",
		Data:    movies,
	})
}

func show(w http.ResponseWriter, r *http.Request, id string) {
	movie, err := getMovie(id)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Movie retrieved",
		Data:    movie,
	})
}

func update(w http.ResponseWriter, r *http.Request, id string) {
	movie, err := getMovie(id)
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusNotFound,
			Message: "Movie not found",
			Data:    nil,
		})
	}
	json.NewDecoder(r.Body).Decode(&movie)
	movie.updateMovie(movie)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Movie updated",
		Data:    movie,
	})
}

func delete(w http.ResponseWriter, r *http.Request, id string) {
	movie, _ := getMovie(id)
	deleteMovie(movie)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Movie deleted",
		Data:    nil,
	})
}

func store(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)
	movie.Id = strconv.Itoa(rand.Intn(100000000000))
	model.Movies = append(model.Movies, movie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusCreated,
		Message: "Movie created",
		Data:    movie,
	})
}

func main() {
	fetchMovies("datastore")
	http.HandleFunc("/movies/", MoviesRouter)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}

func MoviesRouter(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/movies/")
	if len(id) > 0 {
		switch r.Method {
		case http.MethodGet:
			show(w, r, id)
		case http.MethodPatch:
			update(w, r, id)
		case http.MethodDelete:
			delete(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			index(w, r)
		case http.MethodPost:
			store(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

}
