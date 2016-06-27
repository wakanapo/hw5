package transit

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"html/template"
)

type Line struct {
	Name string `json:"Name"`
	Stations []string `json:"Stations"`
}

type Rails struct {
	Lines []Line
}

func HandleTrinsit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, err := ioutil.ReadFile("resource/line.json")
	var rails Rails
	json_err := json.Unmarshal(file, &rails.Lines)
	if err != nil {
		fmt.Fprintln(w, "Format Error: ", json_err)
	}
	t, err := template.ParseFiles("view/transit.html")
   	if err := t.Execute(w, rails); err != nil {
		fmt.Println("Failed to build page", err)
	}
}
