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
	// Update playback direction on loop
	for {
		time.Sleep(time.Second * 5)
		gp.SetFrameDirection(gof.REWIND)
		time.Sleep(time.Second * 5)
		gp.SetFrameDirection(gof.PAUSE)
		time.Sleep(time.Second * 1)
		gp.SetFrameDirection(gof.PLAY)
	}
}()
```

