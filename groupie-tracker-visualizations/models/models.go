package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ARtists Artists
var RElations Relations

type Relations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relationsurl string   `json:"relations"`
	Relations    map[string][]string
}

func HandleJson() {
	ReadArtistJson()
	ReadRelationJson()
}

// this reads the artists file and turns it into a global var called artists
func ReadArtistJson() {

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println(err)
	}

	artist_text, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	jsonerr := json.Unmarshal(artist_text, &ARtists)
	if jsonerr != nil {
		fmt.Println(jsonerr)
	}
}

// reads the relations json
func ReadRelationJson() {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Println(err)
	}

	rawfile, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var temp_stuct struct {
		Data []Relations `json:"index"`
	}

	jsonerr := json.Unmarshal(rawfile, &temp_stuct)
	if jsonerr != nil {
		fmt.Println(jsonerr)
	}
	for index, poo := range temp_stuct.Data {
		ARtists[index].Relations = poo.DatesLocations
	}
}
