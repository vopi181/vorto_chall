package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var PROB_PATH = "probs/problem1.txt"
var loads []Load

const MAX_DIST = 12 * 60

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Point struct {
	x float64
	y float64
}

type Load struct {
	loadno         int
	start          *Point
	end            *Point
	distStartEnd   float64
	distStartDepot float64
	distEndDepot   float64
}

// An Item is something we manage in a priority queue.

func string2Point(line string) *Point {
	tmp := line[1 : len(line)-1]

	raw_point_str := strings.Split(tmp, ",")
	x, err := strconv.ParseFloat(raw_point_str[0], 64)
	check(err)
	y, err := strconv.ParseFloat(raw_point_str[1], 64)
	check(err)
	return &Point{x: x, y: y}
}

func pointDist(p1 *Point, p2 *Point) float64 {
	return math.Sqrt(
		math.Pow(p2.x-p1.x, 2) + math.Pow(p2.y-p1.y, 2),
	)
}

// Sum point dist from depot to first load, from first load to second load, etc.
// until the last load to the depot
func pathDist(path []Load) float64 {
	dist := 0.0

	for i, load := range path {
		dist += load.distStartEnd
		if i == 0 {
			dist += load.distStartDepot
		}
		if i == len(path)-1 {
			dist += load.distEndDepot
		}
		// add distance from load to next load
		if i < len(path)-1 {
			dist += pointDist(load.end, path[i+1].start)
		}

	}

	return dist
}

func successor(path []Load, loads []Load) [][]Load {
	ret := make([][]Load, 0)

	for _, load := range loads {
		tmp := make([]Load, len(path))
		copy(tmp, path)
		tmp = append(tmp, load)
		tmpDist := pathDist(tmp)

		if tmpDist < MAX_DIST {
			found := false
			for _, p := range path {
				if p.loadno == load.loadno {
					found = true
					break
				}
			}

			if !found {
				ret = append(ret, tmp)
			}
		}
	}

	return ret
}

func heuristic(drivers [][]Load) float64 {
	ret := 0.0

	// total cost is 500 * number of drivers + distance
	for _, driver := range drivers {
		ret += pathDist(driver)
		ret += float64(500 * len(drivers))
	}

	// weight how close starting points are to each other and drop off points are to each other
	// heuristicilly finding neighbors
	ret -= drivers[0][0].distStartDepot * 5
	ret -= drivers[len(drivers)-1][len(drivers[len(drivers)-1])-1].distEndDepot * 5

	return ret
}

func main() {
	if len(os.Args) < 2 {
		log.Panicln("Need path to problem text file as argument")
	}
	PROB_PATH = os.Args[1]
	depot_point := &Point{0.0, 0.0}

	file, err := os.Open(PROB_PATH)
	check(err)

	scanner := bufio.NewScanner(file)
	linenum := 0

	for scanner.Scan() {
		if linenum != 0 {
			line := scanner.Text()
			raw_points := strings.Split(line, " ")
			p1 := string2Point(raw_points[1])
			p2 := string2Point(raw_points[2])
			loads = append(loads, Load{
				loadno: linenum,
				start:  p1, end: p2,
				distStartEnd:   pointDist(p1, p2),
				distStartDepot: pointDist(depot_point, p1),
				distEndDepot:   pointDist(p2, depot_point),
			})
		}

		linenum += 1
	}

	queueLoads := make([]Load, len(loads))
	copy(queueLoads, loads)

	var drivers [][]Load
	currDriver := make([]Load, 0)

	startingPath := []Load{queueLoads[0]}
	queueLoads = queueLoads[1:]
	currDriver = append(currDriver, startingPath[0])

	for len(queueLoads) >= 0 {

		// keep track of best heuristic
		bestHeuristic := math.MaxFloat64
		bestPath := make([]Load, 0)
		bestPathIndex := 0

		for _, v := range successor(startingPath, queueLoads) {
			tmpDrivers := make([][]Load, len(drivers))
			copy(tmpDrivers, drivers)
			tmpDrivers = append(tmpDrivers, v)

			tmpHeuristic := heuristic(tmpDrivers)
			if tmpHeuristic < bestHeuristic {
				bestHeuristic = tmpHeuristic
				bestPath = v
				bestPathIndex = v[len(v)-1].loadno
			}
		}

		if len(bestPath) == 0 {

			drivers = append(drivers, currDriver)
			currDriver = nil

			if len(queueLoads) == 0 {
				break
			}

			startingPath = []Load{queueLoads[0]}
			currDriver = append(currDriver, queueLoads[0])

			// remove first load from queue_loads
			queueLoads = queueLoads[1:]
			continue
		} else {
			currDriver = bestPath
		}

		// remove best_path from queue_loads
		for i, load := range queueLoads {
			if load.loadno == bestPathIndex {
				queueLoads = append(queueLoads[:i], queueLoads[i+1:]...)
				break
			}
		}
		// set starting_path to best_path
		startingPath = bestPath
	}

	for _, v := range drivers {
		fmt.Print("[")
		if len(v) == 1 {
			fmt.Print(v[0].loadno)
		} else {
			for i, load := range v {
				if i == len(v)-1 {
					fmt.Print(load.loadno)
				} else {
					fmt.Print(load.loadno, ",")
				}
			}
		}

		fmt.Print("]\n")

		// assert pathDist(v) <= MAX_DIST
		// basic sanity check
		if pathDist(v) > MAX_DIST {
			panic("pathDist(v) > MAX_DIST")
		}

	}

}
