# gof

✨ **a GUI agnostic GIF player with playback control**

## 🚀 Install

```sh
go get -u github.com/claytn/gof
```

## 💡 Usage

```go
gp, err := gof.New("./sample.gif", func(img *image.Image) {
	// GUI rendering logic
})

go func () {
	for {
		// Update playback direction to "rewind"
		gp.SetFrameDirection(gof.REWIND)
		time.Sleep(time.Second * 5)

		// Pause playback
		gp.SetFrameDirection(gof.PAUSE)
		time.Sleep(time.Second * 1)

		// Update playback frame rate delay
		gp.SetDelay(time.Millisecond * 20)

		// Begin playing GIF in "forward" direction
		gp.SetFrameDirection(gof.PLAY)
	}
}()
```

## Spec

* [New](#New)
* [SetFrameDirection](#SetFrameDirection)
* [SetDelay](#SetDelay)

### New 

> Creates and starts a new GIF player

`gof.New` is the entry point to creating a GIF player. Given a valid path to a GIF file, `New` will load the file and begin calling the render function on a delay specified by the GIF file. By default the player will call your custom render function with incrementing frames (i.e in a "forward" direction).

```go
gp, err := gof.New("path/to/gif", func (img *image.Image) {
	// ADD GUI RENDERING LOGIC TO UPDATE CANVAS
})
```

### SetFrameDirection

> Gives you the ability to change the direction frames are "rendered"

Once called, `SetFrameDirection` will change the receiving order of each GIF frame in your custom render function or will halt render calls all-together depending on the provided `frameDirection` argument. 

Options include: `gof.PLAY`, `gof.PAUSE`, `gof.REWIND`

```go
gp, err := gof.New("path/to/gif", func (img *image.Image) {
	// ADD GUI RENDERING LOGIC TO UPDATE CANVAS
})

// Call custom render function with decrementing frame index
gp.SetFrameDirection(gof.REWIND)
```

### SetDelay

> Change the speed of your GIF

Once called, `SetDelay` will change the tick interval between calls to your custom render function.

```go
gp, err := gof.New("path/to/gif", func (img *image.Image) {
	// ADD GUI RENDERING LOGIC TO UPDATE CANVAS
})

// Change tick interval between render calls to 25ms
gp.SetDelay(time.Millisecond * 25)
```





