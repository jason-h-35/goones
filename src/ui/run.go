package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"go_nes/src/nes"
	"log"
	"runtime"
)

const(
	Width = 256
	Height = 240
	title = "Nes Emulator"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	runtime.GOMAXPROCS(2)
	runtime.LockOSThread()
}

func RunUi(n *nes.Nes){
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(Width, Height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	d := newDirector(n, window)
	d.start()
}

