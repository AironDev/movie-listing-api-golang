package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

type CustomError struct{}

type Movie struct {
	Id       string   `json:"id"`
	Title    string   `json:"title"`
	Director Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Movies struct {
	Movies []Movie `json:"model"`
}

var model Movies

func fetchMovies(directory string) {
	data, err := ioutil.ReadFile("datastore/june.json")
	if err != nil {
		log.Fatal("ReadFile Error", err)
	}

	err = json.Unmarshal(data, &model)
	if err != nil {
		log.Fatal("Unmarshal: ", err)
	}

}

func getMovies() []Movie {
	return model.Movies
}

func getMovie(Id string) (Movie, error) {
	for _, movie := range model.Movies {
		if movie.Id == Id {
			return movie, nil
		}
	}
	err := errors.New("something bad")
	return Movie{}, err
}

func addNew(m Movie) {
	model.Movies = append(model.Movies, m)
}

func (m Movie) updateMovie(d Movie) {
	for index, movie := range model.Movies {
		if m.Id == movie.Id {
			model.Movies = append(model.Movies[:index], model.Movies[index+1:]...)
			model.Movies[index] = d
		}
	}
}

func deleteMovie(d Movie) {
	for index, movie := range model.Movies {
		if movie.Id == d.Id {
			model.Movies = append(model.Movies[:index], model.Movies[index+1:]...)
			break
		}
	}
}
