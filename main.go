package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type room struct {
	name      string
	id        int
	visited   bool
	connexion []*room
}

type graphe struct {
	rooms     []*room
	path      []path
	nmbr      int
	start     *room
	end       *room
	interconn [][]int
}

type path struct {
	step   int
	fin    string
	chemin []string
	ants   int
}

func readimput(filename string) int {
	file, err := os.Open(filename)
	reader := bufio.NewScanner(file)
	var input string
	var fourmis int
	for reader.Scan() {
		line := reader.Text()
		splittedline := strings.Split(line, " ")
		if len(splittedline) == 1 {

			input = splittedline[0]
			input = strings.TrimSpace(input)
			
			fourmis, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Not correct number")
				continue
			} else {
				break
			}
		}
	}

	return fourmis
}

func (g *graphe) readfile() []string {
	readFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()
	var maps []string
	var connexions []string
	isconnexions := false
	var start = false
	idnmbr := 0
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#comment") || strings.HasPrefix(line, "#another comment") {
			continue
		}
		if strings.HasPrefix(line, "##start") {
			start = true
			continue
		}
		if strings.HasPrefix(line, "##end") {
			isconnexions = true
			continue
		}
		if len(strings.Split(line, " ")) == 3 {
			maps = append(maps, line)
			id := (strings.Split(line, " ")[0])
			g.rooms = append(g.rooms, &room{name: id, id: idnmbr})
			if isconnexions {
				g.end = g.rooms[idnmbr]
				isconnexions = false
			}
			if start {
				g.start = g.rooms[idnmbr]
				start = false
			}
			idnmbr++
			continue
		}
		if len(strings.Split(line, "-")) == 2 {
			connexions = append(connexions, line)
			continue
		}
	}
	return connexions
}

func (g *graphe) allco(conex []string) {
	for i := 0; i <= g.nmbr-1; i++ {
		var tab []int
		g.interconn = append(g.interconn, tab)
		for j := g.nmbr - 1; j >= 0; j-- {
			g.interconn[i] = append(g.interconn[i], 0)
		}
	}
	for i := g.nmbr - 1; i >= 0; i-- {
		for j := g.nmbr - 1; j >= 0; j-- {
			for k := 0; k < len(conex); k++ {
				a := strings.Split(conex[k], "-")[0]
				v := strings.Split(conex[k], "-")[1]
				if v == g.rooms[i].name && a == g.rooms[j].name {
					g.interconn[i][j] = 1
					g.interconn[j][i] = 1
				}
			}
		}
	}
	for i := 0; i < len(g.rooms); i++ {
		for j := 0; j < len(g.interconn[i]); j++ {
			if g.interconn[i][j] == 1 {
				for m, _ := range g.rooms {
					if g.rooms[m].id == j {
						g.rooms[i].connexion = append(g.rooms[i].connexion, g.rooms[m])
					}
				}
			}
		}
	}
}
func (g *graphe) appendnewpath() bool {
	a := false
	g.start.visited = true
	for i := 0; i < len(g.start.connexion); i++ {
		if !g.start.connexion[i].visited {
			g.path = append(g.path, path{})
			g.path[i].chemin = append(g.path[i].chemin, g.start.connexion[i].name)
			g.start.connexion[i].visited = true
			a = true
		}
	}
	return a
}

func (g *graphe) foundroumbyname(name string) *room {
	for i := 0; i < len(g.rooms); i++ {
		if g.rooms[i].name == name {
			return g.rooms[i]
		}
	}
	return nil
}

func (g *graphe) appendpathstep() bool {
	a := false
	long := len(g.path)
	for i := 0; i < long; i++ {
		room := g.foundroumbyname(g.path[i].chemin[len(g.path[i].chemin)-1])
		first := 0
		for j := 0; j < len(room.connexion); j++ {
			if !room.connexion[j].visited {
				a = true
				if first == 0 {
					g.path[i].chemin = append(g.path[i].chemin, room.connexion[j].name)
					g.path[i].step = len(g.path[i].chemin)
					g.path[i].fin = g.path[i].chemin[len(g.path[i].chemin)-1]
					g.path[i].ants = 0
				} else {
					g.path = append(g.path, path{})
					for k, v := range g.path[i].chemin {
						if k < len(g.path[i].chemin)-1 {
							g.path[len(g.path)-1].chemin = append(g.path[len(g.path)-1].chemin, v)
							g.path[len(g.path)-1].step = len(g.path[len(g.path)-1].chemin)
							g.path[len(g.path)-1].fin = g.path[len(g.path)-1].chemin[len(g.path[len(g.path)-1].chemin)-1]
							g.path[len(g.path)-1].ants = 0
						}
					}
				}
				first++
			}
		}
	}
	for i := 0; i < len(g.path); i++ {
		for j := 0; j < len(g.path[i].chemin); j++ {
			room := g.foundroumbyname(g.path[i].chemin[j])
			if !(room.name == g.end.name) {
				room.visited = true
			}
		}
	}
	return a
}

func (g *graphe) clean() {
	var a []path
	for _, v := range g.path {
		if v.fin == g.end.name {
			if !exist(a, v) {
				a = append(a, v)
			}
		}
	}
	a = sort(a)
	a = removedoublr(a)
	g.path = a
}
func sort(a []path) []path {
	for i := 0; i < len(a); i++ {
		for j := i; j < len(a); j++ {
			if a[i].step > a[j].step {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
	return a
}

func removedoublr(paths []path) []path {
	var shortest []path
	allow := true
	shortest = append(shortest, paths[0])
	for i := 1; i < len(paths); i++ {
		allow = true
		for j := 0; j < len(shortest); j++ {
			if intersect(paths[i].chemin, shortest[j].chemin) {
				allow = false
				if shortest[j].step > paths[i].step {
					shortest[j] = paths[i]
					i = 0
				}
			}
		}
		if allow {
			shortest = append(shortest, paths[i])
		}
	}
	return shortest
}

func intersect(path1, path2 []string) bool {
	for i, c1 := range path1 {
		for j, c2 := range path2 {
			if c1 == c2 {
				if i != len(path1)-1 && j != len(path2)-1 {
					return true
				}
			}
		}
	}

	return false
}

func remove(a []path, v path) []path {
	var x []path
	for _, y := range a {
		if !(strings.Join(v.chemin, "") == strings.Join(y.chemin, "")) {
			x = append(x, y)
		}
	}
	return x
}

func exist(s []path, str path) bool {
	for _, v := range s {
		if strings.Join(v.chemin, "") == strings.Join(str.chemin, "") {
			return true
		}
	}

	return false
}

func (g *graphe) solution(f int) {
	if len(g.path) == 0 {
		fmt.Println("No link")
		return
	}
	for i := 0; i < f; i++ {
		g.path[0].ants++
		g.path[0].step++
		g.path = sort(g.path)
	}
	var result [][]string
	var a []string
	fourmis := 1
	path := 0
	index := 0
	for i := 0; i < len(g.path); i++ {
		for j := 0; j < g.path[i].ants; j++ {
			result = append(result, a)
			for k := 0; k < j; k++ {
				result[index] = append(result[index], "")
			}
			index++
		}
	}
	for i := 0; i < len(g.path); i++ {
		for j := 0; j < f; j++ {
			if g.path[i].ants > 0 {
				for k := 0; k < len(g.path[i].chemin); k++ {

					result[path] = append(result[path], "L"+strconv.Itoa(fourmis)+"-"+g.path[i].chemin[k])

				}
				fourmis++
				path++
				g.path[i].ants--
			}
		}
	}
	fmt.Println(result)
	long := len(result[0])
	for i := 1; i < len(result); i++ {
		if long < len(result[i]) {
			long = len(result[i])
		}
	}
	for i := 0; i < long; i++ {
		for j := 0; j < len(result); j++ {
			if i < len(result[j]) {
				if result[j][i] != "" {
					fmt.Printf(result[j][i] + " ")
				}
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: [./lem-in <file>]")
		os.Exit(0)
	}
	filename := os.Args[1]
	fourmis := readimput(filename)
	var g graphe
	conex := g.readfile()
	g.nmbr = len(g.rooms)
	g.allco(conex)
	g.appendnewpath()
	fin := true
	for fin {
		fin = g.appendpathstep()
	}
	g.clean()
	g.solution(fourmis)
	fmt.Println(g.path)
}
