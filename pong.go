package main

/* TODO
 * ai more realistic
 * win screen
 * improve player select
 */

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type gameState int

const (
	START gameState = iota
	PLAY
	SELECT
)

type gameMode int

const (
	EASY gameMode = iota
	HARD
	PVP
)

var state = START

var letters = [][]byte{
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

var nums = [][]byte{
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
				startX -= size * 3
			}
		}
		startX += size*3 + size
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

	if ball.Y-ball.Radius < 0 {
		ball.YVel = -ball.YVel
		ball.Y = ball.Radius
	} else if ball.Y+ball.Radius > float32(winHeight) {
		ball.YVel = -ball.YVel
		ball.Y = float32(winHeight) - ball.Radius
	}

	if ball.X-ball.Radius < 0 {
		rightPaddle.Score++
		ball.XVel = -300
		if rand.Intn(2) > 0 {
			ball.YVel = float32((rand.Intn(301)))
		} else {
			ball.YVel = float32((rand.Intn(301))) * -1
		}
		ball.pos = getCenter()
		state = START
	} else if ball.X+ball.Radius > float32(winWidth) {
		leftPaddle.Score++
		ball.XVel = 300
		if rand.Intn(2) > 0 {
			ball.YVel = float32((rand.Intn(301)))
		} else {
			ball.YVel = float32((rand.Intn(301))) * -1
		}
		ball.pos = getCenter()
		state = START
	}

	if ball.X < leftPaddle.X+leftPaddle.Width/2 {
		if ball.Y > leftPaddle.Y-leftPaddle.Height/2 && ball.Y < leftPaddle.Y+leftPaddle.Height/2 {
			ball.XVel = -ball.XVel
			ball.X = leftPaddle.X + leftPaddle.Width/2.0 + ball.Radius
			switch y := ball.Y; {
			case y <= leftPaddle.Y+leftPaddle.Height/2 && y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*1):
				ball.YVel += 300
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*1) && y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*2):
				ball.YVel += 150
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*2) && y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*3):
				ball.YVel -= 0
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*3) && y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*4):
				ball.YVel -= 150
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*4) && y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*5):
				ball.YVel -= 300
			default:
				fmt.Println("Collision error, contact dev if you get this")
			}
		}
	}

	if ball.X > rightPaddle.X-rightPaddle.Width/2 {
		if ball.Y > rightPaddle.Y-rightPaddle.Height/2 && ball.Y < rightPaddle.Y+rightPaddle.Height/2 {
			ball.XVel = -ball.XVel
			ball.X = rightPaddle.X - rightPaddle.Width/2.0 - ball.Radius
			switch y := ball.Y; {
			case y <= rightPaddle.Y+rightPaddle.Height/2 && y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*1):
				ball.YVel += 300
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*1) && y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*2):
				ball.YVel += 150
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*2) && y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*3):
				ball.YVel -= 0
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*3) && y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*4):
				ball.YVel -= 150
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*4) && y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*5):
				ball.YVel -= 300
			default:
				fmt.Println("Collision error, contact dev if you get this")
			}
		}
	}
}

type paddle struct {
	pos
	Width  float32
	Height float32
	Speed  float32
	Score  int
	Player int
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
	if paddle.Player == 0 {
		if keyState[sdl.SCANCODE_W] != 0 {
			paddle.Y -= paddle.Speed * elapsedTime
		}
		if keyState[sdl.SCANCODE_S] != 0 {
			paddle.Y += paddle.Speed * elapsedTime
		}
	} else {
		if keyState[sdl.SCANCODE_UP] != 0 {
			paddle.Y -= paddle.Speed * elapsedTime
		}
		if keyState[sdl.SCANCODE_DOWN] != 0 {
			paddle.Y += paddle.Speed * elapsedTime
		}
	}
}

func (paddle *paddle) aiUpdate(ball *ball, diff int, elapsedTime float32) {
	switch diff {
	case 0:
		paddle.Y = ball.Y
	case 1:
		if paddle.Y < ball.Y {
			paddle.Y += paddle.Speed * elapsedTime
		} else if paddle.Y > ball.Y {
			paddle.Y -= paddle.Speed * elapsedTime
		}
	case 2:
		if paddle.Y < ball.Y-ball.Radius {
			paddle.Y += paddle.Speed * elapsedTime
		} else if paddle.Y > ball.Y+ball.Radius {
			paddle.Y -= paddle.Speed * elapsedTime
		}
	case 3:
		if paddle.Y < ball.Y-ball.Radius {
			paddle.Y += paddle.Speed / 1.5 * elapsedTime
		} else if paddle.Y > ball.Y+ball.Radius {
			paddle.Y -= paddle.Speed / 1.5 * elapsedTime
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var mode gameMode
	var validInput bool
	for validInput == false {
		fmt.Println("Select mode:")
		fmt.Printf("\t(0) Player vs Easy Computer\n\t(1) Player vs Hard Computer\n\t(2) Player vs Player\n")
		fmt.Printf("Enter selection (#): ")
		_, err := fmt.Scanf("%d\n", &mode)
		if err != nil {
			fmt.Println("Unrecognized input. Please make a selection by entering the number of the option")
		} else if mode < 0 || mode > 2 {
			if mode == -1 {
				fmt.Print("Entering ludicrous mode")
				for i := 0; i < 6; i++ {
					time.Sleep(400 * time.Millisecond)
					fmt.Print(".")
				}
				fmt.Printf("\nJust kidding. ")
			}
			fmt.Printf("Please select an option on the menu\n")
		} else {
			validInput = true
		}
	}

	// initialize the event checker
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	// create a window with name
	window, err := sdl.CreateWindow("GoPong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	player1 := paddle{pos{50, 100}, 20, 100, 500, 0, 0, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, 500, 0, 1, color{255, 255, 255}}
	ball := ball{getCenter(), 20, 0, 0, color{255, 255, 255}}
	startDirectionX, startDirectionY := rand.Intn(2), rand.Intn(2)
	if startDirectionX > 0 {
		ball.XVel = -300
	} else {
		ball.XVel = 300
	}
	if startDirectionY > 0 {
		ball.YVel = float32(rand.Intn(301))
	} else {
		ball.YVel = float32(rand.Intn(301)) * -1
	}

	keyState := sdl.GetKeyboardState()

	var gameStart time.Time
	var frameStart time.Time
	var elapsedTime float32
	var gameElapsed float32
	var paused bool = false

	// check for any events (mouse, keeb, etc) and close when quit event (hit x) is seen
	for {
		frameStart, gameElapsed = time.Now(), float32(time.Since(gameStart).Seconds())
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		switch state {
		case PLAY:
			switch mode {
			case EASY:
				player1.update(keyState, elapsedTime)   // human player
				player2.aiUpdate(&ball, 2, elapsedTime) //ai player
			case HARD:
				player1.update(keyState, elapsedTime)   // human player
				player2.aiUpdate(&ball, 0, elapsedTime) //ai player
			case PVP:
				player1.update(keyState, elapsedTime) // human player
				player2.update(keyState, elapsedTime) // human player
			}
			if ball.XVel > 0 {
				ball.XVel += gameElapsed / 50
			} else {
				ball.XVel -= gameElapsed / 50
			}
			ball.update(&player1, &player2, elapsedTime)
			if keyState[sdl.SCANCODE_ESCAPE] != 0 {
				state = START
				paused = true
			}
		case START:
			drawSpace(getCenter(), color{255, 255, 255}, 5, pixels)
			screenDraw(tex, renderer, frameStart, &elapsedTime, pixels)
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.Score == 9 || player2.Score == 9 {
					player1.Score = 0
					player2.Score = 0
				}
				if !paused {
					gameStart = time.Now()
				} else {
					paused = false
				}
				state = PLAY
			}
		case SELECT:
			var mode int
			fmt.Println("Select mode:")
			fmt.Printf("\t(0) Player vs Easy Computer\n\t(1) Player vs Hard Computer\n\t(2)Player vs Player\n")
			fmt.Println("Enter selection (#): ")
			fmt.Scanln("%d", &mode)
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		screenDraw(tex, renderer, frameStart, &elapsedTime, pixels)
	}
}
