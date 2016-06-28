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

type lineStatus struct {
	lineNames []string
	relatedStation []*lineStatus
}


func HandleTrinsit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, err := ioutil.ReadFile("resource/line.json")
	var rails Rails
	json_err := json.Unmarshal(file, &rails.Lines)
	if err != nil {
		fmt.Fprintln(w, "Format Error: ", json_err)
	}
	fromStation := r.FormValue("fromStation")
	toStation := r.FormValue("toStation")
	route := searchRoute(fromStation, toStation, setLineStatusFromJson(rails.Lines))
	printRoute(w, route)
	t, err := template.ParseFiles("view/transit.html")
   	if err := t.Execute(w, rails); err != nil {
		fmt.Println("Failed to build page", err)
	}
}

func makeRoute(route []string, next string) []string {
    k := len(route)
    newRoute := make([]string, k + 1)
    copy(newRoute, route)
    newRoute[k] = next
    return newRoute
}

func member(n string, xs []string) bool {
    for _, x := range xs {
        if n == x { return true }
    }
    return false
}

func searchRoute(fromStation, toStation string, railsMap map[string][]string) []string{
    que := make([][]string, 0)
    front := 0
    que = append(que, []string{fromStation})
    for ; front < len(que); front++ {
        route := que[front]
        here := route[len(route) - 1]
        if here == toStation {
            return route
        } else {
            for _, next := range railsMap[here] {
                if !member(next, route) {
                    que = append(que, makeRoute(route, next))
                }
            }
        }
    }
	notFound := []string{"ルートが見つかりません"}
	return notFound
}

func setLineStatusFromJson(lines []Line) map[string][]string {
	var railsMap map[string][]string = make(map[string][]string)
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i].Stations); j++ {
			if j <  len(lines[i].Stations) - 1 {
				railsMap[lines[i].Stations[j]] = append(railsMap[lines[i].Stations[j]], lines[i].Stations[j+1])
			}
			if j > 0 {
				railsMap[lines[i].Stations[j]] = append(railsMap[lines[i].Stations[j]], lines[i].Stations[j-1])
			}
		}
	}
	return railsMap
}

func printRoute(w http.ResponseWriter, route []string) {
	fmt.Fprintf(w, route[0])
	for _, station := range route[1:] {
		fmt.Fprintf(w, "=>")
		fmt.Fprintf(w, station)
	}
}
