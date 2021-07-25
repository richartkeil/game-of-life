package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

var size = 200
var maxAge = 50
var initialDensity = 0.2

// If a cell is healthy, it will remain alive or become alive.
// If a cell is not healthy, it will die or remain dead.
type Cell struct {
	alive      bool
	healthy    bool
	neighbours int
	age        int
}

type Grid struct {
	width, height int
	cells         [][]Cell
}

func main() {
	window, canvas, err := sdlcanvas.CreateWindow(500, 500, "Test")
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	rand.Seed(time.Now().UnixNano())
	grid := getRandomGrid(size, size, initialDensity)

	window.MouseMove = func(x, y int) {
		cellX := size * x / canvas.Width() * 2  // Retina x2
		cellY := size * y / canvas.Height() * 2 // Retina x2
		for i := 0; i < grid.width; i++ {
			for j := 0; j < grid.height; j++ {
				if math.Sqrt(math.Pow(float64(cellX-i), 2)+math.Pow(float64(cellY-j), 2)) <= 10 {
					grid.cells[i][j].alive = rand.Float64() < initialDensity
				}
			}
		}
	}

	previousTime := time.Now()

	window.MainLoop(func() {
		draw(canvas, &grid)
		process(&grid)
		fmt.Printf("Took %vms\n", time.Since(previousTime).Milliseconds())
		previousTime = time.Now()
	})
}

func getRandomGrid(width, height int, initialDensity float64) Grid {
	grid := Grid{width: width, height: height, cells: make([][]Cell, width)}
	for i := 0; i < width; i++ {
		grid.cells[i] = make([]Cell, height)
		for j := 0; j < height; j++ {
			grid.cells[i][j] = Cell{alive: rand.Float64() < initialDensity, healthy: false, age: 0}
		}
	}
	return grid
}

func process(grid *Grid) {
	// Calculate health of cells.
	for i := 0; i < grid.width; i++ {
		for j := 0; j < grid.height; j++ {
			aliveNeighbours := getAliveNeighbours(grid, i, j)
			cell := &grid.cells[i][j]
			cell.healthy = cell.alive
			cell.neighbours = aliveNeighbours
			cell.age++
			if cell.alive && (aliveNeighbours < 2 || aliveNeighbours > 3) {
				cell.healthy = false
			}
			if !cell.alive && aliveNeighbours == 3 {
				cell.healthy = true
			}
			// If a cell changes its state, reset its age.
			if cell.healthy != cell.alive {
				cell.age = 0
			}
		}
	}

	// Update live of cells based on health.
	for i := 0; i < grid.width; i++ {
		for j := 0; j < grid.height; j++ {
			if grid.cells[i][j].age >= maxAge {
				grid.cells[i][j].alive = false
				continue
			}

			if grid.cells[i][j].healthy {
				grid.cells[i][j].alive = true
			} else {
				grid.cells[i][j].alive = false
			}
		}
	}
}

func getAliveNeighbours(grid *Grid, x, y int) int {
	count := 0
	// Left neighbours
	if x > 0 && y > 0 && grid.cells[x-1][y-1].alive {
		count++
	}
	if x > 0 && grid.cells[x-1][y].alive {
		count++
	}
	if x > 0 && y < grid.height-1 && grid.cells[x-1][y+1].alive {
		count++
	}
	// Mid neighbours
	if y > 0 && grid.cells[x][y-1].alive {
		count++
	}
	if y < grid.height-1 && grid.cells[x][y+1].alive {
		count++
	}
	// Right neighbours
	if x < grid.width-1 && y > 0 && grid.cells[x+1][y-1].alive {
		count++
	}
	if x < grid.width-1 && grid.cells[x+1][y].alive {
		count++
	}
	if x < grid.width-1 && y < grid.height-1 && grid.cells[x+1][y+1].alive {
		count++
	}
	return count
}

func draw(canvas *canvas.Canvas, grid *Grid) {
	w, h := float64(canvas.Width()), float64(canvas.Height())
	cellWidth := w / float64(grid.width)
	cellHeight := h / float64(grid.height)

	canvas.SetFillStyle("#41475b")
	canvas.FillRect(0, 0, w, h)

	for i := 0; i < grid.width; i++ {
		for j := 0; j < grid.height; j++ {
			if grid.cells[i][j].alive {
				color := 1 - float64(grid.cells[i][j].age)/float64(maxAge)
				canvas.SetFillStyle(1.0, 1.0, 1.0, color)
				canvas.FillRect(float64(i)*cellWidth, float64(j)*cellHeight, cellWidth, cellHeight)
			}
		}
	}
}
