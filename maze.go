package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
	"math"
	"math/rand"
	"time"
)

type Maze struct {
    cells   []*Cell
    cols    int
    rows    int
    scale   int
}

type Cell struct {
    x       int
    y       int
    walls   uint8
    visited bool
    current bool
}

func MazeAnimation(out io.Writer, cols int, rows int, scale int) {
    // Initialize the maze
    maze := &Maze{
        cells: make([]*Cell, cols * rows),
        cols: cols,
        rows: rows,
        scale: scale,
    }

    // Initialize the cells
    for y := 0; y < rows; y++ {
        for x := 0; x < cols; x++ {
            maze.cells[cols * y + x] = &Cell{x, y, 15, false, false}
        }
    }

    anim := &gif.GIF{}

    maze.generateMaze(anim, 0, 0, 0, NewStack())

    gif.EncodeAll(out, anim)
}

func NewMaze(cols int, rows int, scale int) *Maze {
    // init maze
    maze := &Maze{
        cells: make([]*Cell, cols * rows),
        cols: cols,
        rows: rows,
        scale: scale,
    }

    // init cells
    for y := 0; y < rows; y++ {
        for x := 0; x < cols; x++ {
            maze.cells[cols * y + x] = &Cell{x, y, 15, false, false}
        }
    }

    //maze.generateMaze(0, 0, 0, NewStack())

    return maze
}

func (m *Maze) Png(w io.Writer) {
    // dimensions
    width := m.cols * (2 * m.scale) + m.scale
    height := m.rows * (2 * m.scale) + m.scale

    // setup image
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    // draw cells
    for _, c := range m.cells {
        c.DrawPNG(img, m.scale)
    }
    png.Encode(w, img)
}

func (c *Cell) DrawGIF(img *image.Paletted, scale int) {
    // weight of cell wall in pixels
    weight := int(math.Ceil(float64(scale) / 4))

    // cell corners (including walls)
    x1 := c.x * (scale + weight)
    x2 := x1 + (2 * weight) + scale - 1
    y1 := c.y * (scale + weight)
    y2 := y1 + (2 * weight) + scale - 1

    black := color.RGBA{0, 0, 0, 255}
    white := color.RGBA{255, 255, 255, 255}
    green := color.RGBA{0, 255, 0, 0}

    for x := x1; x <= x2; x++ {
        for y := y1; y <= y2; y++ {
            // initialize all pixels to white
            img.Set(x, y, white)
            // set pixels to black (open) where needed
            if x > x1 && x < x2 {
                // cell body
                if y > y1 && y < y2 {
                    if c.current {
                        img.Set(x, y, green)
                    } else {
                        img.Set(x, y, black)
                    }
                }
                // top wall
                if y < y1 + weight && c.walls & 8 == 0 {
                    img.Set(x, y, black)
                }
                // bottom wall
                if y > y2 - weight && c.walls & 2 == 0 {
                    img.Set(x, y, black)
                }
            }
            if y > y1 && y < y2 {
                // left wall
                if x < x1 + weight && c.walls & 1 == 0 {
                    img.Set(x, y, black)
                }
                // right wall
                if x > x2 - weight && c.walls & 4 == 0 {
                    img.Set(x, y, black)
                }
            }
        }
    }
}

func (c *Cell) DrawPNG(img *image.RGBA, scale int) {
    // weight of cell wall in pixels
    weight := int(math.Ceil(float64(scale) / 4))

    // cell corners (including walls)
    x1 := c.x * (scale + weight)
    x2 := x1 + (2 * weight) + scale - 1
    y1 := c.y * (scale + weight)
    y2 := y1 + (2 * weight) + scale - 1

    black := color.RGBA{0, 0, 0, 255}
    white := color.RGBA{255, 255, 255, 255}

    for x := x1; x <= x2; x++ {
        for y := y1; y <= y2; y++ {
            // initialize all pixels to black
            img.Set(x, y, black)
            // set pixels to white (open) where needed
            if x > x1 && x < x2 {
                // cell body
                if y > y1 && y < y2 {
                    img.Set(x, y, white)
                }
                // top wall
                if y < y1 + weight && c.walls & 8 == 0 {
                    img.Set(x, y, white)
                }
                // bottom wall
                if y > y2 - weight && c.walls & 2 == 0 {
                    img.Set(x, y, white)
                }
            }
            if y > y1 && y < y2 {
                // left wall
                if x < x1 + weight && c.walls & 1 == 0 {
                    img.Set(x, y, white)
                }
                // right wall
                if x > x2 - weight && c.walls & 4 == 0 {
                    img.Set(x, y, white)
                }
            }
        }
    }
}

func (m *Maze) generateMaze(a *gif.GIF, x int, y int, count int, seen *Stack) *Cell {
    c := m.cellAt(x, y)
    c.current = true

    // animation
    weight := int(math.Ceil(float64(m.scale) / 4))
    width := (m.cols - 1) * (m.scale + weight) + (2 * weight) + m.scale
    height := (m.rows - 1) * (m.scale + weight) + (2 * weight) + m.scale
    //width := m.cols * (2 * m.scale) + m.scale
    //height := m.rows * (2 * m.scale) + m.scale
    palette := []color.Color{
        color.RGBA{255, 255, 255, 255},
        color.RGBA{0, 0, 0, 255},
        color.RGBA{0, 255, 0, 255},
    }
    img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
    
    for _, c := range m.cells {
        c.DrawGIF(img, m.scale)
    }

    a.Image = append(a.Image, img)
    a.Delay = append(a.Delay, 0)

    c.current = false
    neighbors := []*Cell{ 
        m.cellAt(x, y - 1),
        m.cellAt(x + 1, y), 
        m.cellAt(x, y + 1), 
        m.cellAt(x - 1, y),
    }
    rand.Seed(time.Now().UnixNano())
    random := rand.Intn(4)
    for i := range neighbors {
        randNeighbor := neighbors[(random + i) % 4]
        if randNeighbor != nil && !randNeighbor.visited {
            randNeighbor.visited = true
            randNeighbor.current = true
           
            c.removeWall(randNeighbor)
            seen.Push(c)
            m.generateMaze(a, randNeighbor.x, randNeighbor.y, count + 1, seen)
            return randNeighbor
        }
    }
    if len(seen.cells) > 0 {
        c, _ := seen.Pop()
        m.generateMaze(a, c.x, c.y, count + 1, seen)
    }
    return nil
}

func (m *Maze) cellAt(x int, y int) *Cell {
    // return nil if index is invalid
    if x < 0 || y < 0 || x > m.cols - 1 || y > m.rows - 1 {
        return nil
    }
    return m.cells[m.cols * y + x]
}

func (c *Cell) direction(n *Cell) uint {
    // direction for getting from Cell(c) to Cell(n)
    if c.x > n.x {          // left
        return 1
    } else if c.x < n.x {   // right
        return 4
    } else if c.y > n.y {   // top
        return 8
    } else {                // bottom
        return 2
    }
}

func (c *Cell) removeWall(n *Cell) {
    // remove wall between Cell(c) and Cell(n)
    dir := uint8(c.direction(n))
    c.walls = c.walls & ^dir
    n.walls = n.walls & ^(dir >> 2 | dir << 2)
}
