package main

/* TODO
 * paddle edge detection
 * 2 player
 * ai more realistic
 * resizing of window
 * load bitmaps
 */

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type gameState int

const (
	START gameState = iota
	PLAY
)

var state = START

var letters = [][]byte {
	{
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
	},
	{
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
	},
	{
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
}

var nums = [][]byte {
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	{
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
	},
	{
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	{
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

type color struct {
	R, G, B byte
}

type pos struct {
	X, Y float32
}

func setPixel(x int, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.R
	}

}

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func screenDraw(tex *sdl.Texture, renderer *sdl.Renderer, frameStart time.Time, elapsedTime *float32, pixels []byte) {
	tex.Update(nil, pixels, winWidth*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	*elapsedTime = float32(time.Since(frameStart).Seconds())
	if *elapsedTime < .005 {
		sdl.Delay(5 - uint32(*elapsedTime/1000.0))
		*elapsedTime = float32(time.Since(frameStart).Seconds())
	}
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.X) - (size*3)/2
	startY := int(pos.Y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func drawSpace(pos pos, color color, size int, pixels []byte) {
	startX := int(pos.X) - (5*(size*3+size))/2
	startY := int(pos.Y) - (size*5)/2 - (winHeight)/4

	for letter := 0; letter < 5; letter++ {
		for i, v := range letters[letter] {
			if v == 1 {
				for y := startY; y < startY+size; y++ {
					for x := startX; x < startX+size; x++ {
						setPixel(x, y, color, pixels)
					}
				}
			}
			startX += size
			if (i+1)%3 == 0 {
				startY += size
				startX -= size*3
			}
		}
		startX += size*3+size
		startY = int(pos.Y) - (size*5)/2 - (winHeight)/4
	}
}

func lerp(a float32, b float32, percent float32) float32 {
	return a + percent*(b-a)
}

type ball struct {
	pos
	Radius float32
	XVel   float32
	YVel   float32
	Color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.Radius; y < ball.Radius; y++ {
		for x := -ball.Radius; x < ball.Radius; x++ {
			if x*x+y*y < ball.Radius*ball.Radius {
				setPixel(int(ball.X+x), int(ball.Y+y), ball.Color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	ball.X += ball.XVel * elapsedTime
	ball.Y += ball.YVel * elapsedTime

	if ball.Y-ball.Radius < 0 || ball.Y+ball.Radius > float32(winHeight) {
		ball.YVel = -ball.YVel
	}

	if ball.X-ball.Radius < 0 {
		rightPaddle.Score++
		ball.pos = getCenter()
		state = START
	} else if ball.X+ball.Radius > float32(winWidth) {
		leftPaddle.Score++
		ball.pos = getCenter()
		state = START
	}

	if ball.X < leftPaddle.X+leftPaddle.Width/2 {
		if ball.Y > leftPaddle.Y-leftPaddle.Height/2 && ball.Y < leftPaddle.Y+leftPaddle.Height/2 {
			ball.XVel = -ball.XVel
			ball.X = leftPaddle.X + leftPaddle.Width/2.0 + ball.Radius
		}
	}

	if ball.X > rightPaddle.X-rightPaddle.Width/2 {
		if ball.Y > rightPaddle.Y-rightPaddle.Height/2 && ball.Y < rightPaddle.Y+rightPaddle.Height/2 {
			ball.XVel = -ball.XVel
			ball.X = rightPaddle.X - rightPaddle.Width/2.0 - ball.Radius
		}
	}
}

type paddle struct {
	pos
	Width  float32
	Height float32
	Speed  float32
	Score  int
	Color  color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.X - paddle.Width/2)
	startY := int(paddle.Y - paddle.Height/2)

	for y := 0; y < int(paddle.Height); y++ {
		for x := 0; x < int(paddle.Width); x++ {
			setPixel(startX+x, startY+y, paddle.Color, pixels)
		}
	}

	numX := lerp(paddle.X, getCenter().X, 0.2)
	drawNumber(pos{numX, 35}, paddle.Color, 10, paddle.Score, pixels)
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.Y -= paddle.Speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.Y += paddle.Speed * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	paddle.Y = ball.Y
}

func main() {

	// initialize the event checker
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	// create a window with name
	window, err := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	// create a renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	// create texture
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	ball := ball{getCenter(), 20, 400, 400, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	// check for any events (mouse, keeb, etc) and close when quit event (hit x) is seen
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if state == PLAY {
			player1.update(keyState, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == START {
			drawSpace(getCenter(), color{255, 255, 255}, 5, pixels)
			screenDraw(tex, renderer, frameStart, &elapsedTime, pixels)
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.Score == 9 || player2.Score == 9 {
					player1.Score = 0
					player2.Score = 0
				}
				state = PLAY
			}
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		screenDraw(tex, renderer, frameStart, &elapsedTime, pixels)
	}

}
