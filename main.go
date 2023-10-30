package main

import (
	"math"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const GRID_WIDTH = 10
const GRID_HEIGHT = 10

type gridArrT [GRID_WIDTH][GRID_HEIGHT]bool

type vec2 struct {
	x int
	y int
}

func main() {
	window, cnv, err := sdlcanvas.CreateWindow(500, 500, "Hello")
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	cellSize := calcCellSize(cnv)
	grid := gridArrT{}
	editing := true

	mousePos := vec2{0, 0}
	window.MouseMove = func(x, y int) {
		mousePos.x = x
		mousePos.y = y
	}

	window.MouseDown = func(btn, x, y int) {
		cell := getCellFromMouse(cnv, mousePos, cellSize)
		grid = onMouseClick(grid, cell)
	}

	window.KeyDown = func(scancode int, rn rune, name string) {
		if name != "KeyP" {
			return
		}

		editing = false
		cnv.ClearRect(0, 0, float64(cnv.Width()), float64(cnv.Height()))
	}

	window.MainLoop(func() {
		w, h := float64(cnv.Width()), float64(cnv.Height())
		cnv.ClearRect(0, 0, w, h)

		drawGrid(cnv, grid, cellSize)

		if editing {
			highlightMouse(cnv, mousePos, cellSize)
			return
		}

		time.Sleep(1 * time.Second)

		grid = step(grid)
	})
}

func calcCellSize(cnv *canvas.Canvas) vec2 {
	width := cnv.Width() / GRID_WIDTH
	height := cnv.Height() / GRID_HEIGHT
	return vec2{width, height}
}

func drawGrid(cnv *canvas.Canvas, grid gridArrT, cellSize vec2) {
	for i, row := range grid {
		for j, item := range row {
			if !item {
				continue
			}
			cnv.SetFillStyle("#FFF")
			cnv.FillRect(
				float64(i*cellSize.x), float64(j*cellSize.y),
				float64(cellSize.x), float64(cellSize.y),
			)
		}
	}
}

func getCellFromMouse(cnv *canvas.Canvas, mousePos vec2, cellSize vec2) vec2 {
	// need to divide by 2 because of retina screen issues
	return vec2{
		int(math.Floor(float64(mousePos.x) / float64(cellSize.x/2))),
		int(math.Floor(float64(mousePos.y) / float64(cellSize.y/2))),
	}
}

func highlightMouse(cnv *canvas.Canvas, mousePos vec2, cellSize vec2) {
	if mousePos.x == 0 || mousePos.y == 0 {
		return
	}

	cell := getCellFromMouse(cnv, mousePos, cellSize)

	cnv.SetFillStyle("#FF0000")
	cnv.FillRect(
		float64(cell.x*cellSize.x),
		float64(cell.y*cellSize.y),
		float64(cellSize.x),
		float64(cellSize.y),
	)
}

func onMouseClick(grid gridArrT, cell vec2) gridArrT {
	grid[cell.x][cell.y] = !grid[cell.x][cell.y]
	return grid
}

func step(grid gridArrT) gridArrT {

	newGrid := grid

	for i, row := range grid {
		for j, cell := range row {
			liveNeighbors := getNeighbors(grid, i, j)

			// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			if cell && liveNeighbors < 2 {
				newGrid[i][j] = false
			}
			// Any live cell with two or three live neighbours lives on to the next generation.
			if cell && liveNeighbors == 2 || liveNeighbors == 3 {
				newGrid[i][j] = true
			}

			// Any live cell with more than three live neighbours dies, as if by overpopulation.
			if cell && liveNeighbors > 3 {
				newGrid[i][j] = false
			}

			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			if !cell && liveNeighbors == 3 {
				newGrid[i][j] = true
			}
		}
	}

	return newGrid
}

func getNeighbors(grid gridArrT, i, j int) int {
	neighbors := 0

	positions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, pos := range positions {
		x, y := i+pos[0], j+pos[1]
		if x >= 0 && x < GRID_HEIGHT && y >= 0 && y < GRID_WIDTH && grid[x][y] {
			neighbors++
		}
	}

	return neighbors
}
