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

type Suspend struct {
	From string `json:"From"`
	To string `json:"To"`
}

func HandleTrinsit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file_line, err_line := ioutil.ReadFile("resource/line.json")
	file_suspend, err_suspend := ioutil.ReadFile("resource/suspend.json")
	
	var rails Rails
	json_err_line := json.Unmarshal(file_line, &rails.Lines)
	if err_line != nil {
		fmt.Fprintln(w, "Format Error: ", json_err_line)
	}
	
	var suspends []Suspend
	json_err_suspend := json.Unmarshal(file_suspend, &suspends)
	if err_suspend != nil {
		fmt.Fprintln(w, "Format Error: ", json_err_suspend)
	}
	for _, suspend := range suspends {
		fmt.Fprintf(w, "%s - %s間で運転を見合わせています。</br>", suspend.From, suspend.To)
	}
	fromStation := r.FormValue("fromStation")
	toStation := r.FormValue("toStation")
	railsMap := setLineStatusFromJson(rails.Lines)
	fmt.Fprintln(w, railsMap["多摩川"])
	stationStatus := setStaitionStatus(rails.Lines, suspends)
	route := searchRoute(fromStation, toStation, railsMap, stationStatus)
	printRoute(w, route, rails.Lines)
	t,_ := template.ParseFiles("view/transit.html")
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

func searchRoute(fromStation, toStation string, railsMap map[string][]string, suspendsMap map[string]bool) []string{
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
					if (here == fromStation && suspendsMap[next] || next == toStation && suspendsMap[here] || suspendsMap[here] && suspendsMap[next]) {
						que = append(que, makeRoute(route, next))
					}
                }
            }
        }
    }
	notFound := []string{"ルートが見つかりません"}
	return notFound
}

func setLineStatusFromJson(lines []Line) map[string][]string {
	var railsMap map[string][]string = make(map[string][]string)
	for _, line := range lines {
		frontStation := ""
		for _, station := range line.Stations {
			railsMap[frontStation] = append(railsMap[frontStation], station)
			railsMap[station] = append(railsMap[station], frontStation)
			frontStation = station
		}
	}
	return railsMap
}

func printRoute(w http.ResponseWriter, route []string, lines []Line) {
	fmt.Fprint(w, "<br>")
	fmt.Fprintf(w, route[0])
	front := route[0]
	var frontline, line string
	for _, station := range route[1:] {
		line = nowLine(lines, front, station)
		if len(frontline) > 0 && line != frontline {
			fmt.Fprintf(w, "（%s）", frontline)
			fmt.Fprintf(w, "=> ")
			fmt.Fprintf(w, front)
		}
		front = station
		frontline = line
	}
	if len(line) > 0 {
		fmt.Fprintf(w, "（%s）", line)
		fmt.Fprintf(w, "=> ")
		fmt.Fprintf(w, route[len(route)-1])
		fmt.Fprintf(w, "（%s）", line)
	}
}
 



func nowLine(lines []Line, front string, now string) string{
	var nowline string
	for _, line := range lines {
		if isIncludingTheseStations(line, front, now) {
			nowline = line.Name
			return nowline
		}
	}
	return nowline
}

func isIncludingTheseStations(line Line, station1 string, station2 string) bool {
	for _, station := range line.Stations {
		if (station == station1) {
			for _, station := range line.Stations {
				if (station == station2) {
					return true
				}
			}
		}
	}
	return false
}

func setStaitionStatus(lines []Line, suspends []Suspend) map[string]bool {
	var suspendsMap map[string]bool = make(map[string]bool)
	for _, line := range lines {
		for _, station := range line.Stations {
			suspendsMap[station] = true
		}
	}
	for _, line := range lines {
		for _, suspend := range suspends {
			if isIncludingTheseStations(line, suspend.From, suspend.To) {
				frag := 0
				for _, station := range line.Stations{
					if station == suspend.From {
						frag++
					}
					if (frag > 0) {
						suspendsMap[station] = false
					}
					if station == suspend.To {
						frag--
					}
				}
			}
		}
	}
	return suspendsMap
}
