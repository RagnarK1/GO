package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	nrAnts     int
	startName  string
	endName    string
	totalPaths [][]string
	goodPaths  [][][]string
	shortPath  int
	setIndex   int
	indexNr    int
)

type node struct {
	Roomname    string
	Connections map[string]*node
	Xcord       int
	Ycord       int
	Visited     bool
}

func AssignRoom(name string, x int, y int) *node {
	return &node{
		Roomname:    name,
		Connections: map[string]*node{},
		Xcord:       x,
		Ycord:       y,
		Visited:     false,
	}
}

func (room *node) AddConnection(connection *node) {
	if _, connected := room.Connections[connection.Roomname]; !connected {
		room.Connections[connection.Roomname] = connection
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: missing filename")
		fmt.Println("USAGE: go run main.go testing/example00.txt")
		return
	}

	testCase := os.Args[1]
	readCase, _ := os.ReadFile(testCase)

	tempG := InputHandler(string(readCase))

	fmt.Printf("Start room: %v \n", startName)
	fmt.Printf("End room: %v \n", endName)
	fmt.Printf("Number of ants: %v \n", nrAnts)

	trail := []string{}
	DfsSearch(startName, trail, tempG)

	Combo(MinorFilter(totalPaths))

	SelectPath()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

}

func CtrAnts(line string) {
	if line != "" {
		a, err := strconv.Atoi(line)
		if err != nil {
			HandleError("Whole ants please")
		}
		if a < 1 {
			HandleError("No ants to move?")
		} else {
			nrAnts = a
			return
		}
		HandleError("First line not formatted correctly")
	}
}

func InputHandler(tempS string) map[string]*node {
	chart := make(map[string]*node)
	lines := strings.Split(tempS, "\n")

	for l, line := range lines {

		if line != "" && line[0] == 35 {
			if line == "##start" {
				d := strings.Split(lines[l+1], " ")
				startName = d[0]
			}
			if line == "##end" {
				f := strings.Split(lines[l+1], " ")
				endName = f[0]
			}
			continue
		}

		if l == 0 {
			CtrAnts(line)
		}

		room := strings.Split(line, " ")
		if len(room) == 3 {
			if _, cache := chart[room[0]]; !cache {
				x, er1 := strconv.Atoi(room[1])
				y, er2 := strconv.Atoi(room[2])
				if er1 != nil || er2 != nil {
					fmt.Println("coordinate formatted wrong, " + room[0] + " skipped")
					continue
				}
				chart[room[0]] = AssignRoom(room[0], x, y)
			} else {
				fmt.Println("already exists,  " + room[0] + " skipped")
				continue
			}
		}

		conns := strings.Split(line, "-")
		if len(conns) == 2 {
			if _, err1 := chart[conns[0]]; !err1 {
				HandleError("Not pointed from a Node " + conns[0])
			}
			if _, err2 := chart[conns[1]]; !err2 {
				HandleError("Not pointed to a Node " + conns[1])
			}
			if conns[0] != conns[1] {
				chart[conns[0]].AddConnection(chart[conns[1]])
				chart[conns[1]].AddConnection(chart[conns[0]])
			} else {
				HandleError("Connected to room " + conns[0])
			}
		}
	}
	return chart
}

func MinorFilter(all [][]string) [][]string {
	var arrSlc []string
	for l := range all {
		if len(all[l]) == 0 {
			return sortPathsByLength(all)
		}
		arrSlc = append(arrSlc, strings.Join(all[l], " . "))
	}
	for l := range arrSlc {
		for k := range arrSlc {
			if k == l {
				continue
			} else if strings.Contains(arrSlc[k], arrSlc[l]) {
				arrSlc[k] = "pathwontwork"
			}
		}
	}

	var e2 []string
	for l := range arrSlc {
		if arrSlc[l] != "pathwontwork" {
			e2 = append(e2, arrSlc[l])
		}
	}

	for l := 0; l < len(e2)-1; l++ {
		if len(e2[l]) > len(e2[l+1]) {
			e2[l], e2[l+1] = e2[l+1], e2[l]
			l = -1
		}
	}

	var e3 [][]string
	for l := range e2 {
		e3 = append(e3, strings.Split(e2[l], " . "))
	}

	return e3
}

func sortPathsByLength(e2 [][]string) [][]string {
	for l := 0; l < len(e2)-1; l++ {
		if len(e2[l]) < len(e2[l+1]) {
			e2[l], e2[l+1] = e2[l+1], e2[l]
			l = -1
		}
	}
	return e2
}

func SelectPath() {
	line := make([][]int, 0)

	for r, prev := range goodPaths {

		var antRow [][]int
		for range prev {
			antRow = append(antRow, make([]int, 0))
		}

		for l := 1; l <= nrAnts; l++ {
			k := 10000
			for z := range prev {
				token := len(prev[z]) + len(antRow[z])
				if k > token {
					k = token
					indexNr = z
				} else {
					continue
				}
			}
			antRow[indexNr] = append(antRow[indexNr], l)
		}
		move := 0
		for l := range antRow {
			token := len(prev[l]) + len(antRow[l])
			if move == 0 || move < token {
				move = token
			}
		}
		if shortPath == 0 || move < shortPath {
			shortPath = move
			setIndex = r
			line = antRow
		}
	}
	OutputHandler(goodPaths[setIndex], shortPath, line)
}

func Combo(paths [][]string) {
	for i := range paths {
		var set [][]string
		set = append(set, paths[i])
		for k := range paths {
			if i == k || CrashCheck(set, paths[k]) {
				continue
			} else {
				set = append(set, paths[k])
			}
		}
		goodPaths = append(goodPaths, set)
	}
}

func CrashCheck(set [][]string, X []string) bool {
	for l := range set {
		for x := range X {
			for ii := range set[l] {
				if X[x] == set[l][ii] {
					return true
				}
			}
		}
	}
	return false
}

func DfsSearch(cue string, trail []string, chart map[string]*node) {
	chart[cue].Visited = true
	if cue != startName && cue != endName {
		trail = append(trail, chart[cue].Roomname)
	}
	if cue == endName { 
		temp := make([]string, len(trail))
		copy(temp, trail)
		totalPaths = append(totalPaths, temp)
	} else {
		for k, cont := range chart[cue].Connections {
			if !cont.Visited {
				DfsSearch(k, trail, chart)
			}
		}
	}
	chart[cue].Visited = false
}

func OutputHandler(solution [][]string, turns int, antID [][]int) {
	start := time.Now()
	set := make([][]string, turns)
	copy(set, solution)
	if solution[0][0] == "" {
		var last string
		for a := range antID[0] {
			last += "L" + strconv.Itoa(antID[0][a]) + "-" + endName + " "
		}
		fmt.Println(last)
		return
	}
	for l := range solution {
		set[l] = append(set[l], endName)
	}

	for r := 0; r < turns; r++ {
		var last string
		for d := range solution {
			for a := range antID[d] {
				plc := r - a
				if plc < 0 || plc > len(solution[d]) {
					continue
				} else {
					last += "L" + strconv.Itoa(antID[d][a]) + "-" + set[d][plc] + " "
				}
			}
		}
		fmt.Println(last)
	}
	fmt.Printf("Found %v paths in %v.\n", len(totalPaths), time.Since(start))
	fmt.Printf("Used quickest path possible with %v turns.\n", turns)
}

func HandleError(message string) {
	fmt.Printf("Error: %s\n", message)
	os.Exit(1)
}
