package gof

import (
	"fmt"
	"image"
	"testing"
	"time"
)

func TestGofPlayerCreation(t *testing.T) {
	_, err := New("./sample.gif", func(img *image.Image) {})

	if err != nil {
		t.Errorf("Failed to create GofPlayer instance")
	}
}

func TestGofPlayerRender(t *testing.T) {
	ch := make(chan string)
	gp, _ := New("./sample.gif", func(img *image.Image) {
		time.Sleep(1 * time.Second)
		ch <- "rendered"
	})

	select {
	case <-ch:
		// Do nothing
	case <-time.After(gp.delay):
		t.Errorf(fmt.Sprintf("Expected render to be called within %s", gp.delay))
	}
}
