package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
)

const (
	planRoom    = "0"
	planWall    = "-1"
	planVisited = "-2"
)

type FloorPlan [][]string

type RoomChairs map[string]int

type RoomPlanner struct {
	floorPlan  FloorPlan
	rooms      map[string]RoomChairs
	roomPos    map[string]string
	chairPos   map[string]string
	chairTypes map[string]bool
}

func NewRoomPlanner(fileName string) *RoomPlanner {
	fp := &RoomPlanner{
		floorPlan: nil,
		rooms:     make(map[string]RoomChairs),
		roomPos:   make(map[string]string),
		chairPos:  make(map[string]string),
		chairTypes: map[string]bool{
			"W": true,
			"P": true,
			"S": true,
			"C": true,
		},
	}

	if err := fp.readFloorPlan(fileName); err != nil {
		log.Errorf("Can't read file %s, %s", fileName, err)
		return nil
	}

	return fp
}

func (rp *RoomPlanner) readFloorPlan(fileName string) error {
	var (
		row int
		col int
	)
	row = 0

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	rp.floorPlan = make(FloorPlan, 0)

	for scanner.Scan() {
		var (
			nameStarted bool
			planLine    []string
			val         string
			roomName    []string
		)

		nameStarted = false
		roomName = make([]string, 0)

		col = 0
		for _, s := range scanner.Text() {
			ss := string(s)

			if nameStarted {
				roomName = append(roomName, ss)
			}
			if !nameStarted && len(roomName) > 0 {
				roomNameStr := strings.Join(roomName[:len(roomName)-1], "")
				rp.rooms[roomNameStr] = nil
				rp.roomPos[fmt.Sprintf("%d-%d", row, col-1)] = roomNameStr
				roomName = nil
			}

			switch ss {
			case "+", "-", "|", "\\", "/":
				val = planWall
			case " ":
				val = planRoom
			case "W", "P", "S", "C":
				rp.chairPos[fmt.Sprintf("%d-%d", row, col)] = ss
				val = planRoom
			case "(":
				nameStarted = true
				val = planRoom
			case ")":
				nameStarted = false
				val = planRoom
			default:
				val = "0"
			}
			planLine = append(planLine, val)
			col += 1
		}
		rp.floorPlan = append(rp.floorPlan, planLine)
		row += 1
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (rp *RoomPlanner) traverse() {
	var roomName string

	for i := 0; i < len(rp.floorPlan); i++ {
		for j := 0; j < len(rp.floorPlan[i]); j++ {
			if rp.floorPlan[i][j] != planRoom {
				continue
			}
			currentRoom := RoomChairs{
				"W": 0,
				"P": 0,
				"S": 0,
				"C": 0,
			}
			roomName = ""
			queue := [][]int{{i, j}}
			for len(queue) > 0 {
				node := queue[0]
				queue = queue[1:]

				row := node[0]
				col := node[1]

				rp.floorPlan[row][col] = planVisited

				if v, ok := rp.roomPos[fmt.Sprintf("%d-%d", row, col)]; ok {
					roomName = v
				}

				if v, ok := rp.chairPos[fmt.Sprintf("%d-%d", row, col)]; ok {
					currentRoom[v] += 1
				}

				for _, m := range [][]int{{-1, 0}, {1, 0}, {0, 1}, {0, -1}} {
					nI := row + m[0]
					nJ := col + m[1]
					if 0 <= nI && nI < len(rp.floorPlan) && 0 <= nJ && nJ < len(rp.floorPlan[nI]) {
						if rp.floorPlan[nI][nJ] == planRoom {
							rp.floorPlan[nI][nJ] = planVisited
							queue = append(queue, []int{nI, nJ})
						}
					}
				}
			}

			// if we outside of the room
			if roomName != "" {
				rp.rooms[roomName] = currentRoom
			}
		}
	}

	total := RoomChairs{
		"W": 0,
		"P": 0,
		"S": 0,
		"C": 0,
	}

	for _, room := range rp.rooms {
		for c, v := range room {
			total[c] += v
		}
	}

	rp.rooms["_total"] = total
}

func (rp *RoomPlanner) print() {
	roomNames := make([]string, 0)

	for k := range rp.rooms {
		roomNames = append(roomNames, k)
	}

	sort.Strings(roomNames)

	for _, room := range roomNames {
		switch room {
		case "_total":
			fmt.Println("total:")
		default:
			fmt.Printf("%s: \n", room)
		}

		ch := make([]string, 0)
		for chair, count := range rp.rooms[room] {
			ch = append(ch, fmt.Sprintf("%s: %d", chair, count))
		}
		fmt.Printf("%s \n", strings.Join(ch, ", "))
	}

}

func main() {
	if len(os.Args) != 1 {
		panic(fmt.Errorf("USAGE: room-planner <floor plan file>"))
	}

	rp := NewRoomPlanner(os.Args[0])
	rp.traverse()
	rp.print()
}
