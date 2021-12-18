package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Map struct {
	rows [][]byte
}

func (m *Map) size() uint {
	// The map is a square
	return uint(len(m.rows))
}

func (m *Map) riskAt(x uint, y uint) uint {
	return uint(m.rows[y][x])
}

func (m *Map) copyMap(source *Map, xTile uint, yTile uint) {
	delta := xTile + yTile
	startY := yTile * source.size()
	limitY := startY + source.size()
	for y := startY; y < limitY; y++ {
		startX := xTile * source.size()
		limitX := startX + source.size()
		for x := startX; x < limitX; x++ {
			risk := source.riskAt(x-startX, y-startY) + delta
			if risk >= 10 {
				risk -= 9
			}
			m.rows[y][x] = byte(risk)
		}
	}
}

func processInputLine(line string) []byte {
	var risks = make([]byte, 0, len(line))

	for _, c := range line {
		var riskLevel byte

		riskLevel = byte(c - '0')
		if riskLevel >= 0 && riskLevel <= 9 {
			risks = append(risks, riskLevel)
		}
	}
	return risks
}

func readMap(fileName string) (*Map, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines = make([][]byte, 0, 100)

	for scanner.Scan() {
		risks := processInputLine(scanner.Text())
		if len(risks) > 0 {
			lines = append(lines, risks)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Map{lines}, nil
}

type TotalRiskMap struct {
	totalRisk [][]uint
}

func min(x uint, y uint) uint {
	if x < y {
		return x
	}
	return y
}

func (trm *TotalRiskMap) setInitialRisk(m *Map, x uint, y uint) {
	if x == 0 {
		if y > 0 {
			trm.totalRisk[y][x] = m.riskAt(x, y) + trm.totalRisk[y-1][x]
		}
	} else if y == 0 {
		trm.totalRisk[y][x] = m.riskAt(x, y) + trm.totalRisk[y][x-1]
	} else {
		trm.totalRisk[y][x] = m.riskAt(x, y) + min(trm.totalRisk[y][x-1], trm.totalRisk[y-1][x])
	}
}

func (trm *TotalRiskMap) minimizeRisk(m *Map, x uint, y uint) bool {
	limit := m.size() - 1
	riskAtPos := m.riskAt(x, y)
	minRiskUpToNow := trm.totalRisk[y][x] - riskAtPos
	smallestRisk := uint(0xFFFFFFFF)
	if x > 0 {
		risk := trm.totalRisk[y][x-1]
		smallestRisk = min(smallestRisk, risk)
	}
	if x < limit {
		risk := trm.totalRisk[y][x+1]
		smallestRisk = min(smallestRisk, risk)
	}
	if y > 0 {
		risk := trm.totalRisk[y-1][x]
		smallestRisk = min(smallestRisk, risk)
	}
	if y < limit {
		risk := trm.totalRisk[y+1][x]
		smallestRisk = min(smallestRisk, risk)
	}
	if smallestRisk < minRiskUpToNow {
		trm.totalRisk[y][x] = smallestRisk + riskAtPos
		return true
	}
	return false
}

func findBestPath(m *Map) uint {
	var data = make([][]uint, 0, m.size())

	for i := uint(0); i < m.size(); i++ {
		var totalRiskLine = make([]uint, m.size())

		data = append(data, totalRiskLine)
	}

	trm := TotalRiskMap{data}
	limit := m.size()

	for y := uint(0); y < limit; y++ {
		for x := uint(0); x < limit; x++ {
			trm.setInitialRisk(m, x, y)
		}
	}
	var changes uint

LOOP:

	changes = 0
	for y := limit - 1; y >= 1; y-- {
		for x := limit - 1; x >= 1; x-- {
			if trm.minimizeRisk(m, x, y) {
				changes++
			}
		}
	}
	if changes > 0 {
		changes = 0
		for y := uint(0); y < limit; y++ {
			for x := uint(0); x < limit; x++ {
				if trm.minimizeRisk(m, x, y) {
					changes++
				}
			}
		}
		if changes > 0 {
			goto LOOP
		}
	}
	return trm.totalRisk[m.size()-1][m.size()-1]
}

func growMap(m *Map, factor uint) *Map {
	rows := make([][]byte, m.size()*5)
	for i := 0; i < len(rows); i++ {
		rows[i] = make([]byte, m.size()*5)
	}
	largeMap := &Map{rows}
	for y := uint(0); y < factor; y++ {
		for x := uint(0); x < factor; x++ {
			largeMap.copyMap(m, x, y)
		}
	}
	return largeMap
}

func findBestPathForInput(fileName string) (uint, uint) {
	m, err := readMap(fileName)
	if err != nil {
		log.Fatal(err)
	}
	firstResult := findBestPath(m)
	m = growMap(m, 5)
	secondResult := findBestPath(m)
	return firstResult, secondResult
}

func main() {
	lowestRisk1, lowestRisk2 := findBestPathForInput("input.txt")
	fmt.Printf("Lowest risk is: %d\n", lowestRisk1)
	fmt.Printf("Lowest risk on enloarged map is: %d\n", lowestRisk2)
}
