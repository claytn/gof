package gof

import (
	"image"
	"testing"
	"time"
)

func TestGofPlayerCreation(t *testing.T) {
	_, err := New("./assets/test.gif", func(img *image.Image) {})

	if err != nil {
		t.Errorf("Failed to create GofPlayer instance")
	}
}

func TestGofPlayerRender(t *testing.T) {
	ch := make(chan string)
	gp, _ := New("./assets/test.gif", func(img *image.Image) {
		ch <- "rendered"
	})

	select {
	case <-ch:
		// Do nothing
	case <-time.After(gp.delay):
		t.Errorf("Expected render to be called within %s", gp.delay)
	}
}

/**
 * The below tests rely on the test.gif file specifically which contains a
 * single pixel gif that increases it's R value in the RGBA color of the pixel.
 * This makes it easy to test if the playback controls are working as expected
 */
func TestGofRenderOrderPlay(t *testing.T) {
	ch := make(chan *image.Image)
	gp, _ := New("./assets/test.gif", func(img *image.Image) {
		ch <- img
	})
	// Redundant direction set for readability
	gp.SetFrameDirection(PLAY)

	var prev *image.Image = nil
	frame := 0
	frameCount := len(gp.frames)
	for img := range ch {
		frame++
		if frame >= frameCount {
			// Consider this test done after looping through the full GIF
			break
		}

		if prev == nil {
			prev = img
			continue
		}

		r, _, _, _ := (*img).At(0, 0).RGBA()
		oldR, _, _, _ := (*prev).At(0, 0).RGBA()
		if r <= oldR {
			t.Errorf("Expected frames to be played in an increasing order")
		}

		prev = img
	}
}

func TestGofRenderOrderRewind(t *testing.T) {
	ch := make(chan *image.Image)
	gp, _ := New("./assets/test.gif", func(img *image.Image) {
		ch <- img
	})
	gp.SetFrameDirection(REWIND)

	var prev *image.Image = nil
	frame := 0
	frameCount := len(gp.frames)
	for img := range ch {
		frame++
		if frame >= frameCount {
			// Consider this test done after looping through the full GIF
			break
		}

		if prev == nil {
			prev = img
			continue
		}

		r, _, _, _ := (*img).At(0, 0).RGBA()
		oldR, _, _, _ := (*prev).At(0, 0).RGBA()
		/**
		 * Only enforce this once passed the first two frames
		 * On reverse, the first frame is still the 0th frame so going from 0th to (len - 1)
		 * will result in this test failing.
		 */
		if r >= oldR && frame > 2 { // The first frame is the 0th frame so going in reverse ...
			t.Errorf("Expected frames to be played in an decreasing order")
		}

		prev = img
	}
}
