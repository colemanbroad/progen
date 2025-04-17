package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/fzipp/canvas"
)

type Map struct {
	bytes  []int
	nx, ny int
}

type State struct {
	x, y int
}

// clips to edge
func (m Map) get(x, y int) int {
	if x < 0 || x >= m.nx || y < 0 || y >= m.ny {
		return -1
	}
	idx := x + m.nx*y
	return m.bytes[idx]
}

func (m Map) set(x, y int, b int) {
	m.bytes[x+m.nx*y] = b
}

var maphistory Map

func init_maphistory() {
	maphistory = Map{
		bytes: make([]int, 10*10),
		nx:    10,
		ny:    10,
	}
	for i := range maphistory.bytes {
		maphistory.bytes[i] = 0
	}
	pitvisit = make([]Site, 0, 1_000)
}

// TODO
func draw_line_scalar(ctx *canvas.Context, x0, y0, xn, yn float64) {

}

func draw_pitsworld_state(ctx *canvas.Context) {
	Nx := 700.0
	Ny := 700.0
	ctx.ClearRect(0, 0, Nx, Ny)

	// vertical bars
	for i := 0; i <= maphistory.nx; i++ {
		x0 := float64(i) * Nx / float64(maphistory.nx)
		y0 := 0.0
		yn := Ny
		// TODO: fixme:
		draw_line_scalar(ctx, x0, y0, x0, yn)
	}

	// horizontal bars
	for j := 0; j <= maphistory.ny; j++ {
		y0 := float64(j) * Ny / float64(maphistory.ny)
		x0 := 0.0
		xn := Nx
		// TODO: fixme:
		draw_line_scalar(ctx, x0, y0, xn, y0)
	}

	// text of count per square
	ctx.SetTextAlign(canvas.AlignCenter)
	ctx.SetTextBaseline(canvas.BaselineMiddle)
	ctx.SetFillStyle(color.Black)
	for i := range maphistory.nx {
		for j := range maphistory.nx {
			x0 := (float64(i) + 0.5) * Nx / float64(maphistory.nx)
			y0 := (float64(j) + 0.5) * Ny / float64(maphistory.ny)
			pitvisit = append(pitvisit, Site{x0, y0})
			ctx.FillText(fmt.Sprint(maphistory.get(i, j)), x0, y0)
		}
	}
}

type Site struct{ x, y float64 }

var pitvisit []Site

func reward_from_pitsworld_visits(s State) float64 {
	count := maphistory.get(s.x, s.y)
	maphistory.set(s.x, s.y, count+1)
	rew := 1.0 / float64(count+1)
	return rew
}

func reward_from_pitsworld_distance(s State) float64 {
	fx := float64(s.x)
	fy := float64(s.y)
	return fx + fy
	// return math.Sqrt(math.Pow(fx, 2) + math.Pow(fy, 2))
}

// Pitsworld Test Fragments Below!!

func newMap() Map {
	// fmt.Println("newMap")
	m := Map{
		nx:    10,
		ny:    10,
		bytes: make([]int, 10*10),
	}

	for x := range m.nx {
		for y := range m.ny {

			m.set(x, y, 0)

			// random walls
			if rand.Float64() < 0.05 {
				m.set(x, y, 1)
			}
			// a line of wall
			if x == 4 {
				m.set(x, y, 1)
			}
			// a single opening
			if y == 8 {
				m.set(x, y, 0)
			}
		}
	}

	return m
}

func start(m Map) State {
	// fmt.Println("start")
	return State{x: 2, y: 2}
}

// only allowed to move if into empty space
func move(s State, m Map, v int) State {
	// fmt.Println("move")

	Found_reward(reward_from_pitsworld_visits(s), "move")

	// up
	if v%4 == 0 {
		if s.y == 0 {
			return s
		}
		if m.get(s.x, s.y-1) == 0 {
			return State{s.x, s.y - 1}
		}
	}
	// right
	if v%4 == 1 {
		if s.x == m.nx-1 {
			return s
		}
		if m.get(s.x+1, s.y) == 0 {
			return State{s.x + 1, s.y}
		}
	}
	// down
	if v%4 == 2 {
		if s.y == m.ny-1 {
			return s
		}
		if m.get(s.x, s.y+1) == 0 {
			return State{s.x, s.y + 1}
		}
	}
	// left
	if v%4 == 3 {
		if s.x == 0 {
			return s
		}
		if m.get(s.x-1, s.y) == 0 {
			return State{s.x - 1, s.y}
		}
	}

	//unreachable
	// panic(unreachable)
	return State{}
}

// int-encoded observation of the adjacent squares
func observe(m Map, s State) int {
	// we would prefer an encoding that was symmetric in the four directions
	// this may require returning a vector and having vector ops in the stdlib.
	r := 0
	r += 1 * m.get(s.x, s.y+1)
	r += 2 * m.get(s.x+1, s.y)
	r += 4 * m.get(s.x, s.y-1)
	r += 8 * m.get(s.x-1, s.y)
	return r
}

func add_pitsworld(lib Library) {
	// map () -> Map -- gen new map
	// start (map) -> State -- set starting position on map
	// move (map, state, move) -> state -- update position, get reward, maybe die
	lib.addFuncToLibrary(start, "start", []Type{"Map"}, "State")
	lib.addFuncToLibrary(move, "move", []Type{"State", "Map", "int"}, "State")
	lib.addFuncToLibrary(newMap, "newMap", []Type{}, "Map")
	lib.addFuncToLibrary(observe, "observe", []Type{"Map", "State"}, "int")
}
