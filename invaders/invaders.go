package invaders

import (
	"fmt"

	"github.com/protoshark/invaders8080/cpu"
	"github.com/veandco/go-sdl2/sdl"
)

// Invaders game struct
type Invaders struct {
	cpu         cpu.CPU
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	frameBuffer []uint8

	ports [9]uint8

	shiftOffset   uint8
	shiftRegister uint16
}

const InitialROMOffset uint32 = 0x0000

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

	// game.ports[1] = 1 << 3

	return game
}

// Load space invaders into cpu memory
func (game *Invaders) loadROM(romPath string) {
	file := sdl.RWFromFile(romPath, "rb")
	defer file.Close()

	fmt.Printf("Loading %s\n", romPath)

	file.Read(game.cpu.Memory[InitialROMOffset:])
}

// Run space invaders
func (game *Invaders) Run(romPath string) {
	game.setup(romPath)
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
			running = game.handleEvents()
			game.draw()

			if game.cpu.IntEnable {
				game.cpu.Interrupt(0x10)
			}

			timer = sdl.GetTicks()
		}
	}
}

func (game *Invaders) setup(romPath string) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	// create the window
	window, err := sdl.CreateWindow("Space Invaders", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth*2, ScreenHeight*2, sdl.WINDOW_RESIZABLE)

	if err != nil {
		sdl.Quit()
		panic(err)
	}

	// set minimum size
	window.SetMinimumSize(ScreenWidth, ScreenHeight)

	// hide cursor
	sdl.ShowCursor(sdl.DISABLE)

	game.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		sdl.Quit()
		panic(err)
	}

	if game.renderer.SetLogicalSize(ScreenWidth, ScreenHeight) != nil {
		game.renderer.Destroy()
		window.Destroy()
		sdl.Quit()
		panic(err)
	}

	game.texture, err = game.renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)

	if err != nil {
		game.renderer.Destroy()
		window.Destroy()
		sdl.Quit()
		panic(err)
	}

	game.updateScreen()
	game.loadROM(romPath)
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
				case sdl.SCANCODE_C: // INSERT COIN
					game.ports[1] |= 0x01

				case sdl.SCANCODE_S: // P1 START
					game.ports[1] |= 1 << 2
				case sdl.SCANCODE_RETURN: // P2 START
					game.ports[1] |= 1 << 1

				case sdl.SCANCODE_W: // P1 SHOOT
					game.ports[1] |= 1 << 4
				case sdl.SCANCODE_A: // P1 LEFT
					game.ports[1] |= 1 << 5
				case sdl.SCANCODE_D: // P1 RIGHT
					game.ports[1] |= 1 << 6

				case sdl.SCANCODE_UP: // P2 SHOOT
					game.ports[2] |= 1 << 4
				case sdl.SCANCODE_LEFT: // P2 LEFT
					game.ports[2] |= 1 << 5
				case sdl.SCANCODE_RIGHT: // P2 RIGHT
					game.ports[2] |= 1 << 6
				}
			}
			if e.Type == sdl.KEYUP {
				switch key {
				case sdl.SCANCODE_C:
					game.ports[1] &= ^uint8(1 << 0)

				case sdl.SCANCODE_S:
					game.ports[1] &= ^uint8(1 << 2)
				case sdl.SCANCODE_RETURN:
					game.ports[1] &= ^uint8(1 << 1)

				case sdl.SCANCODE_W:
					game.ports[1] &= ^uint8(1 << 4)
				case sdl.SCANCODE_A:
					game.ports[1] &= ^uint8(1 << 5)
				case sdl.SCANCODE_D:
					game.ports[1] &= ^uint8(1 << 6)

				case sdl.SCANCODE_UP:
					game.ports[2] &= ^uint8(1 << 4)
				case sdl.SCANCODE_LEFT:
					game.ports[2] &= ^uint8(1 << 5)
				case sdl.SCANCODE_RIGHT:
					game.ports[2] &= ^uint8(1 << 6)
				}
			}
		}
	}

	return true
}

func (game *Invaders) handleIn(port uint8) {
	if port == 2 || port == 4 {
		return
	}

	game.ports[port] = game.cpu.A
}

func (game *Invaders) handleOut(port uint8) {
	if port == 3 {
		return
	}

	game.cpu.A = game.ports[port]
}

func (game *Invaders) update() {
	game.cpu.Cycles = 0
	for game.cpu.Cycles < uint32(CyclesPerFrames/2) {
		opcode := game.cpu.NextByte()

		// emulate shift register
		if opcode == 0xd3 {
			nextByte := game.cpu.NextByte()
			if nextByte == 2 {
				game.shiftOffset = game.cpu.A
			} else if nextByte == 4 {
				game.shiftRegister = (uint16(game.cpu.A) << 8) | (uint16(game.shiftRegister) >> 8)
			}
		} else if opcode == 0xdb {
			if game.cpu.NextByte() == 3 {
				game.cpu.A = uint8(game.shiftRegister >> (8 - game.shiftOffset))
			}
		}

		game.cpu.Step(false)

		if opcode == 0xd3 {
			game.handleIn(game.cpu.NextByte())
			game.cpu.PC++
		}
		if opcode == 0xdb {
			game.handleOut(game.cpu.NextByte())
			game.cpu.PC++
		}
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
