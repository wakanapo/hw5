package app

import (
	"net/http"
	"wordshuffle"
	"transit"
)

func init() {
	http.HandleFunc("/pata", wordshuffle.HandlePata)
	http.HandleFunc("/transit", transit.HandleTrinsit)
}


