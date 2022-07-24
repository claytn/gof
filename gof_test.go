package gof

import (
	"image"
	"testing"
)

func TestGofPlayerCreation(t *testing.T) {
	_, err := New("./sample.gif", func(img *image.Image) {

	})

	if err != nil {
		t.Errorf("Failed to create GofPlayer instance")
	}
}
