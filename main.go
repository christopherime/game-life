package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Cell struct {
	Alive bool
	Color *color.Color
}

type GameOfLife struct {
	Cells  [][]Cell
	Width  int
	Height int
}

func NewGameOfLife(width int, height int) *GameOfLife {
	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
	}
	return &GameOfLife{cells, width, height}
}

func (g *GameOfLife) Alive(x, y int) bool {
	return g.Cells[y][x].Alive
}

func (g *GameOfLife) LiveNeighbors(x, y int) int {
	var count int
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i != 0 || j != 0) && g.Alive((x+j+g.Width)%g.Width, (y+i+g.Height)%g.Height) {
				count++
			}
		}
	}
	return count
}

func (g *GameOfLife) Next() {
	newCells := make([][]Cell, g.Height)
	for i := range newCells {
		newCells[i] = make([]Cell, g.Width)
	}
	for y, row := range g.Cells {
		for x, cell := range row {
			liveNeighbors := g.LiveNeighbors(x, y)
			newCells[y][x].Alive = cell.Alive && liveNeighbors >= 2 && liveNeighbors <= 3 || !cell.Alive && liveNeighbors == 3
			newCells[y][x].Color = cell.Color
			if newCells[y][x].Alive && newCells[y][x].Color == nil {
				newCells[y][x].Color = randomColor()
			}
		}
	}
	g.Cells = newCells
}

func randomColor() *color.Color {
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
		color.New(color.FgWhite),
	}
	return colors[rand.Intn(len(colors))]
}

func main() {
    rand.Seed(time.Now().UnixNano())

    g := NewGameOfLife(60, 30)

    // Initialize some cells randomly
    for i := range g.Cells {
        for j := range g.Cells[i] {
            g.Cells[i][j].Alive = rand.Float32() > 0.5
            if g.Cells[i][j].Alive {
                g.Cells[i][j].Color = randomColor()
            }
        }
    }

    // Setup to listen for SIGINT (ctrl+c)
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("\n- Ctrl+C pressed in Terminal")
        os.Exit(0)
    }()

    // Clear the terminal
    fmt.Print("\033[H\033[2J")

    // Create a copy of the initial state of the grid
    prevCells := make([][]Cell, g.Height)
    for i := range prevCells {
        prevCells[i] = make([]Cell, g.Width)
        copy(prevCells[i], g.Cells[i])
    }

    // Run the game
    for {
        g.Next()

        // Move the cursor to the top-left corner of the terminal
        fmt.Print("\033[H")

        // Display the current state
        for y, row := range g.Cells {
            for x, cell := range row {
                // Only update the cell if it has changed
                if cell.Alive != prevCells[y][x].Alive || cell.Color != prevCells[y][x].Color {
                    // Move the cursor to the position of the cell
                    fmt.Printf("\033[%d;%dH", y+1, x+1)

                    if cell.Alive {
                        cell.Color.Print("*")
                    } else {
                        fmt.Print(" ")
                    }

                    // Remember the new state of the cell
                    prevCells[y][x] = cell
                }
            }
        }

        time.Sleep(time.Millisecond * 200)
    }
}

