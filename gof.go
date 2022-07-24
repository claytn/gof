package gof

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"sync"
	"time"
)

const PLAY = 1
const REWIND = -1
const PAUSE = 0

type GofPlayer struct {
	frames      []image.Image
	render      func(*image.Image)
	delay       time.Duration
	frame       int
	direction   int
	delayMu     sync.Mutex
	directionMu sync.Mutex
}

func New(path string, render func(img *image.Image)) (*GofPlayer, error) {
	g, err := loadGIF(path)
	if err != nil {
		return nil, err
	}

	frames := unlayerFrames(g)

	gp := &GofPlayer{
		frames:    frames,
		render:    render,
		delay:     100,
		direction: PAUSE, // default to pause
	}

	go gp.tick()

	return gp, nil
}

func (g *GofPlayer) Play() {
	g.setDirection(PLAY)
}

func (g *GofPlayer) Pause() {
	g.setDirection(PAUSE)
}

func (g *GofPlayer) Rewind() {
	g.setDirection(REWIND)
}

func (g *GofPlayer) SetDelay(delay time.Duration) {
	g.delayMu.Lock()
	g.delay = delay
	g.delayMu.Unlock()
}

func (g *GofPlayer) setDirection(direction int) {
	g.directionMu.Lock()
	g.direction = direction
	g.directionMu.Unlock()
}

func (g *GofPlayer) getDirection() int {
	g.directionMu.Lock()
	defer g.directionMu.Unlock()
	return g.direction
}

func (g *GofPlayer) tick() {
	frameCount := len(g.frames)
	if g.frame >= frameCount {
		g.frame = 0
	} else if g.frame < 0 {
		g.frame = frameCount - 1
	}
	direction := g.getDirection()
	// Avoid calling render function when paused
	if direction != 0 {
		go g.render(&g.frames[g.frame])
	}

	g.frame += direction

	// Safely read the delay
	g.delayMu.Lock()
	delay := g.delay
	g.delayMu.Unlock()

	time.Sleep(delay)
	g.tick() // Relying on tail call recursion here
}

func loadGIF(path string) (*gif.GIF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func unlayerFrames(g *gif.GIF) []image.Image {
	images := make([]image.Image, len(g.Image))
	images[0] = g.Image[0]

	bounds := g.Image[0].Bounds()
	for i := 1; i < len(g.Image); i++ {
		unlayered := image.NewRGBA(bounds)
		curr := g.Image[i]
		prev := images[i-1]
		for x := 0; x < bounds.Dx(); x++ {
			for y := 0; y < bounds.Dy(); y++ {
				/**
				 * Only use the current frame pixel if the pixel is not transparent
				 * and it is within the specified bounds of the current image
				 */
				if !isTransparent(curr.At(x, y)) && isInBounds(x, y, curr.Rect) {
					unlayered.Set(x, y, curr.At(x, y))
				} else {
					unlayered.Set(x, y, prev.At(x, y))
				}
			}
		}
		images[i] = unlayered
	}
	return images
}

func isTransparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0
}

func isInBounds(x int, y int, rect image.Rectangle) bool {
	return rect.Min.X <= x && x < rect.Max.X && rect.Min.Y <= y && y < rect.Max.Y
}
