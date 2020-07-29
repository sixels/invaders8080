package invaders

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/protoshark/invaders8080/cpu"
)

// Invaders game struct
type Invaders struct {
	cpu         cpu.CPU
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	frameBuffer []uint8

	port1 uint8
	port2 uint8

	shift0      uint8
	shift1      uint8
	shiftOffset uint8
}

// InvadersRoms path and memory offset
var InvadersRoms = map[uint16]string{
	0x0000: "./test/invaders.h",
	0x0800: "./test/invaders.g",
	0x1000: "./test/invaders.f",
	0x1800: "./test/invaders.e",
}

// Screen dimensions
const (
	ScreenWidth  int32 = 224
	ScreenHeight int32 = 256
)

// Machine config
var (
	FPS             = 60
	Frames          = 1000. / FPS
	CyclesPerFrames = 2000 * Frames
)

// New Invaders
func New() Invaders {
	game := Invaders{
		cpu: cpu.New(),
	}
	game.frameBuffer = make([]uint8, ScreenWidth*ScreenHeight*4)

	game.port1 = 1 << 3
	game.port2 = 0

	return game
}

// Load space invaders into cpu memory
func (game *Invaders) loadInvaders() {
	for offset, path := range InvadersRoms {
		file := sdl.RWFromFile(path, "rb")
		defer file.Close()

		fmt.Printf("Loading %s %s 0x%04x\n", path, "into offset", offset)

		file.Read(game.cpu.Memory[offset:])
	}
}

// Run space invaders
func (game *Invaders) Run() {
	game.setup()
	defer sdl.Quit()
	defer game.renderer.Destroy()
	defer game.texture.Destroy()

	// start the timer
	timer := sdl.GetTicks()

	running := true
	for running {
		// running = game.handleEvents()

		if sdl.GetTicks()-timer >= uint32(Frames)/2 {
			game.update()
			if game.cpu.IntEnable {
				game.cpu.Interrupt(0x08)
			}
			game.update()
			if game.cpu.IntEnable {
				game.cpu.Interrupt(0x10)
			}

			game.draw()

			running = game.handleEvents()
			timer = sdl.GetTicks()
		}
	}
}

func (game *Invaders) setup() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	// defer sdl.Quit()

	// create the window
	window, err := sdl.CreateWindow("Space Invaders", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth*2, ScreenHeight*2, sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	// defer window.Destroy()

	// set minimum size
	window.SetMinimumSize(ScreenWidth, ScreenHeight)

	// hide cursor
	sdl.ShowCursor(sdl.DISABLE)

	game.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	// defer game.renderer.Destroy()

	err = game.renderer.SetLogicalSize(ScreenWidth, ScreenHeight)
	if err != nil {
		panic(err)
	}

	game.texture, err = game.renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)
	if err != nil {
		panic(err)
	}
	// defer game.texture.Destroy()

	game.updateScreen()
	game.loadInvaders()
}

func (game *Invaders) updateScreen() {
	game.texture.Update(nil, game.frameBuffer, 4*int(ScreenWidth))
	game.renderer.Clear()
	game.renderer.Copy(game.texture, nil, nil)
	game.renderer.Present()
}

func (game *Invaders) handleEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			key := e.Keysym.Scancode
			if e.Type == sdl.KEYDOWN {
				switch key {
				case sdl.SCANCODE_ESCAPE:
					return false
				case sdl.SCANCODE_C:
					fmt.Print(".")
					game.port1 |= 1 << 0
				case sdl.SCANCODE_RETURN:
					fmt.Print(",")
					game.port1 |= 1 << 2
				case sdl.SCANCODE_SPACE:
					fmt.Print("^")
					game.port1 |= 1 << 4
				case sdl.SCANCODE_LEFT:
					fmt.Print("<")
					game.port1 |= 1 << 5
				case sdl.SCANCODE_RIGHT:
					fmt.Print(">")
					game.port1 |= 1 << 6
				}
			}
			if e.Type == sdl.KEYUP {
				switch key {
				case sdl.SCANCODE_C:
					game.port1 &= ^uint8(1 << 0)
				case sdl.SCANCODE_RETURN:
					game.port1 &= ^uint8(1 << 2)
				case sdl.SCANCODE_SPACE:
					game.port1 &= ^uint8(1 << 4)
				case sdl.SCANCODE_LEFT:
					game.port1 &= ^uint8(1 << 5)
				case sdl.SCANCODE_RIGHT:
					game.port1 &= ^uint8(1 << 6)
				}
			}
		}
	}

	return true
}

func (game *Invaders) handleInOut(opcode uint8) {
	switch opcode {
	case 0xdb: // IN
		port := game.cpu.NextByte()
		game.cpu.PC++

		result := &game.cpu.A

		switch port {
		case 1:
			*result = game.port1
		case 2:
			*result = game.port2
		case 3:
			value := (uint16(game.shift1) << 8) | uint16(game.shift0)
			*result = uint8(value >> (8 - game.shiftOffset))
		}

	case 0xd3: // OUT
		port := game.cpu.NextByte()
		game.cpu.PC++

		switch port {
		case 2:
			game.shiftOffset = game.cpu.A
		case 4:
			game.shift0 = game.shift1
			game.shift1 = game.cpu.A
		}
	}
}

func (game *Invaders) update() {
	game.cpu.Cycles = 0
	for game.cpu.Cycles < uint32(CyclesPerFrames/2) {
		opcode := game.cpu.NextByte()
		game.cpu.Step(false)
		game.handleInOut(opcode)
	}
}

func (game *Invaders) draw() {
	for i := 0; i < 256*224/8; i++ {
		x := i * 8 % 256
		y := i * 8 / 256

		pix := game.cpu.Memory[cpu.VRAMOffset+i]

		for b := 0; b < 8; b++ {
			// rotate 90 deg
			px := y
			py := -(x + b) + int(ScreenHeight) - 1

			index := (py*int(ScreenWidth) + px) * 4
			if (pix>>b)&1 != 0 {
				game.frameBuffer[index], game.frameBuffer[index+1], game.frameBuffer[index+2] = 0xff, 0xff, 0xff
			} else {
				game.frameBuffer[index], game.frameBuffer[index+1], game.frameBuffer[index+2] = 0x00, 0x00, 0x00
			}
		}
	}

	game.updateScreen()
}
