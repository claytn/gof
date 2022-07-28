package gof

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"sync"
	"time"
)

type frameDirection int

const PLAY frameDirection = 1
const REWIND frameDirection = -1
const PAUSE frameDirection = 0

type GofPlayer struct {
	frames    []image.Image
	render    func(*image.Image)
	delay     time.Duration
	frame     int
	direction frameDirection
	mu        sync.RWMutex
}

func New(path string, render func(img *image.Image)) (*GofPlayer, error) {
	g, err := loadGIF(path)
	if err != nil {
		return nil, err
	}

	frames := unlayerFrames(g)

	var delayMs int
	if len(g.Delay) > 0 {
		delayMs = g.Delay[0] * 10 // Gif Delay contains values in 100ths of a second
	} else {
		delayMs = 30 // If no delay specified in GIF file, default to 30ms
	}

	gp := &GofPlayer{
		frames:    frames,
		render:    render,
		delay:     time.Millisecond * time.Duration(delayMs),
		direction: PLAY,
	}

	go gp.tick()

	return gp, nil
}

func (g *GofPlayer) SetDelay(delay time.Duration) {
	g.mu.Lock()
	g.delay = delay
	g.mu.Unlock()
}

func (g *GofPlayer) SetFrameDirection(direction frameDirection) {
	g.mu.Lock()
	g.direction = direction
	g.mu.Unlock()
}

func (g *GofPlayer) getDirection() frameDirection {
	g.mu.RLock()
	defer g.mu.RUnlock()
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

	g.frame += int(direction)

	// Safely read the delay
	g.mu.RLock()
	delay := g.delay
	g.mu.RUnlock()

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
