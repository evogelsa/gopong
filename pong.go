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

// define window size -- this affects gameplay so 800x600 is a good standard size
const winWidth, winHeight int = 800, 600

// game state enum for screen type
type gameState int

const (
	START gameState = iota
	PLAY
	SELECT
)

// select screen game mode enum
type gameMode int

const (
	EASY gameMode = iota
	HARD
	PVP
)

// game state starts on start screen
var state = START

// byte arrays containing the letters for "space"
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

// byte array containing numbers 0 - 9
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

// struct for storing pixel color data
type color struct {
	R, G, B byte
}

// struct for storing pixel position data
type pos struct {
	X, Y float32
}

// setPixel sets an individual pixel at the specified location inside pixels array to provided color
func setPixel(x int, y int, c color, pixels []byte) {
	// pixels array is a flattened array of the screen data which has four bits for each pixel
	// convert x y position to index in this array
	index := (y*winWidth + x) * 4

	// set the index and following indecies to the desired color, ignore alpha channel
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}

}

// getCenter returns the position of the center pixel based on the winWidth and winHeight
func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

// clear iteratively zeros all pixels in given pixel array
func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

// screenDraw updates rendered textures and draws to window It takes the frame times in order to preserve the physics
// time independence from framerate
func screenDraw(tex *sdl.Texture, renderer *sdl.Renderer, frameStart time.Time, elapsedTime *float32, pixels []byte) {
	// update texture with newly drawn pixels array
	tex.Update(nil, pixels, winWidth*4)
	// copy texture to renderer
	renderer.Copy(tex, nil, nil)
	// push to window
	renderer.Present()

	// physics-framerate independence
	// cap max framerate at 200 but ensure that any physics calculations are scaled to accomodate for changes in timing
	// between frames
	*elapsedTime = float32(time.Since(frameStart).Seconds())
	if *elapsedTime < .005 {
		sdl.Delay(5 - uint32(*elapsedTime*1000.0))
		*elapsedTime = float32(time.Since(frameStart).Seconds())
	}
}

// drawNumber takes a position and color and draws the number from nums array to screen. The size sets a multiplier
// of how big to scale the number
func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	// get starting x and y positions based on size
	startX := int(pos.X) - (size*3)/2
	startY := int(pos.Y) - (size*5)/2

	// draw one pixel in nums as a square of given size
	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		// move to next pixel
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

// drawSpace draws the word space at the given location with given color and scales it with size parameter
func drawSpace(pos pos, color color, size int, pixels []byte) {
	// get starting x and y pos based on size and adjust so the word space is centered
	startX := int(pos.X) - (5*(size*3+size))/2
	startY := int(pos.Y) - (size*5)/2 - (winHeight)/4

	// 5 letters
	for letter := 0; letter < 5; letter++ {
		// draw one pixel as a square of given size
		for i, v := range letters[letter] {
			if v == 1 {
				for y := startY; y < startY+size; y++ {
					for x := startX; x < startX+size; x++ {
						setPixel(x, y, color, pixels)
					}
				}
			}
			// move to next pixel
			startX += size
			if (i+1)%3 == 0 {
				startY += size
				startX -= size * 3
			}
		}
		// move to next letter
		startX += size*3 + size
		startY = int(pos.Y) - (size*5)/2 - (winHeight)/4
	}
}

// lerp is linear interpolation
func lerp(a float32, b float32, percent float32) float32 {
	return a + percent*(b-a)
}

// ball struct stores information relevant to the pong ball
type ball struct {
	pos
	Radius float32
	XVel   float32
	YVel   float32
	Color  color
}

// draw is a method acting on ball type which draws the ball based on its stored information
func (ball *ball) draw(pixels []byte) {
	// define ball with a square but only draw to pixel if its within desired radius (given by ball struct)
	for y := -ball.Radius; y < ball.Radius; y++ {
		for x := -ball.Radius; x < ball.Radius; x++ {
			if x*x+y*y < ball.Radius*ball.Radius {
				setPixel(int(ball.X+x), int(ball.Y+y), ball.Color, pixels)
			}
		}
	}
}

// update is a method acting on ball type which updates the balls values based on its position
// this method handles collision cases with the paddles as well
func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	// move ball based on curren velocities and elapsed time
	ball.X += ball.XVel * elapsedTime
	ball.Y += ball.YVel * elapsedTime

	// top and bottom screen collision detection
	// reverse balls y velocities on collision
	if ball.Y-ball.Radius < 0 {
		ball.YVel = -ball.YVel
		ball.Y = ball.Radius
	} else if ball.Y+ball.Radius > float32(winHeight) {
		ball.YVel = -ball.YVel
		ball.Y = float32(winHeight) - ball.Radius
	}

	// left and right ball collision detection
	// add to corresponding players score and reset ball location and game state
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

	// ball paddle collision detection, reverses ball x velocity on collision but variably changes
	// y velocity based on collision position
	if ball.X < leftPaddle.X+leftPaddle.Width/2 {
		if ball.Y > leftPaddle.Y-leftPaddle.Height/2 && ball.Y < leftPaddle.Y+leftPaddle.Height/2 {
			ball.XVel = -ball.XVel
			ball.X = leftPaddle.X + leftPaddle.Width/2.0 + ball.Radius
			// paddes are divided into 5 sections, check which section ball is colliding with and
			// add velocity in y component
			switch y := ball.Y; {
			// top most
			case y <= leftPaddle.Y+leftPaddle.Height/2 &&
				y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*1):
				ball.YVel += 300
			// second to top
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*1) &&
				y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*2):
				ball.YVel += 150
			// middle
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*2) &&
				y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*3):
				ball.YVel -= 0
			// second to bottom
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*3) &&
				y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*4):
				ball.YVel -= 150
			// bottom
			case y <= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*4) &&
				y >= leftPaddle.Y+leftPaddle.Height/2-((leftPaddle.Height/5)*5):
				ball.YVel -= 300
			// incase I did something wrong
			default:
				fmt.Println("Collision error, contact dev if you get this")
			}
		}
	}

	// same as above but for right paddle collisions
	if ball.X > rightPaddle.X-rightPaddle.Width/2 {
		if ball.Y > rightPaddle.Y-rightPaddle.Height/2 && ball.Y < rightPaddle.Y+rightPaddle.Height/2 {
			ball.XVel = -ball.XVel
			ball.X = rightPaddle.X - rightPaddle.Width/2.0 - ball.Radius
			switch y := ball.Y; {
			case y <= rightPaddle.Y+rightPaddle.Height/2 &&
				y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*1):
				ball.YVel += 300
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*1) &&
				y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*2):
				ball.YVel += 150
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*2) &&
				y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*3):
				ball.YVel -= 0
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*3) &&
				y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*4):
				ball.YVel -= 150
			case y <= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*4) &&
				y >= rightPaddle.Y+rightPaddle.Height/2-((rightPaddle.Height/5)*5):
				ball.YVel -= 300
			default:
				fmt.Println("Collision error, contact dev if you get this")
			}
		}
	}
}

// paddle struct to stored relevant information to paddles
type paddle struct {
	pos
	Width  float32
	Height float32
	Speed  float32
	Score  int
	Player int
	Color  color
}

// draw is a method acting on paddles type which draws the paddles to the screen
func (paddle *paddle) draw(pixels []byte) {
	// start paddles centered in y
	startX := int(paddle.X - paddle.Width/2)
	startY := int(paddle.Y - paddle.Height/2)

	// draw rectangle
	for y := 0; y < int(paddle.Height); y++ {
		for x := 0; x < int(paddle.Width); x++ {
			setPixel(startX+x, startY+y, paddle.Color, pixels)
		}
	}

	// use interpolation to get paddle position
	numX := lerp(paddle.X, getCenter().X, 0.2)
	// draw score associated with paddle
	drawNumber(pos{numX, 35}, paddle.Color, 10, paddle.Score, pixels)
}

// update is a method acting on paddles which handles paddle movement
func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	// check for key pressed associated with given player and move paddle when appropriate
	// muliply speed by elapsed time to preserve behavior during frame loss
	if paddle.Player == 0 {
		if keyState[sdl.SCANCODE_W] != 0 {
			paddle.Y -= paddle.Speed * elapsedTime
		}
		if keyState[sdl.SCANCODE_S] != 0 {
			paddle.Y += paddle.Speed * elapsedTime
		}
	} else {
		if keyState[sdl.SCANCODE_UP] != 0 || keyState[sdl.SCANCODE_I] != 0 {
			paddle.Y -= paddle.Speed * elapsedTime
		}
		if keyState[sdl.SCANCODE_DOWN] != 0 || keyState[sdl.SCANCODE_K] != 0 {
			paddle.Y += paddle.Speed * elapsedTime
		}
	}
}

// aiUpdate is a method acting on paddles which handles automated movement of paddles
func (paddle *paddle) aiUpdate(ball *ball, diff int, elapsedTime float32) {
	switch diff {
	// "hard ai" -- only way to beat is to get ball to move fast enough it clips through paddle
	case 0:
		paddle.Y = ball.Y
	// ai paddle moves at paddle speed but keeps ball centered
	case 1:
		if paddle.Y < ball.Y {
			paddle.Y += paddle.Speed * elapsedTime
		} else if paddle.Y > ball.Y {
			paddle.Y -= paddle.Speed * elapsedTime
		}
	// "easy ai" -- ai paddle moves at paddle speed but has some error given by ball radius
	// when keeping ball centered on the paddle
	case 2:
		if paddle.Y < ball.Y-ball.Radius {
			paddle.Y += paddle.Speed * elapsedTime
		} else if paddle.Y > ball.Y+ball.Radius {
			paddle.Y -= paddle.Speed * elapsedTime
		}
	// make the ai player slower for really easy games
	case 3:
		if paddle.Y < ball.Y-ball.Radius {
			paddle.Y += paddle.Speed / 1.5 * elapsedTime
		} else if paddle.Y > ball.Y+ball.Radius {
			paddle.Y -= paddle.Speed / 1.5 * elapsedTime
		}
	}
}

// main func contains window create with sdl and the game loop
func main() {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())
	// print exit message and wait for 1.2 seconds before closing (avoid immediate terminal exit on windows)
	defer time.Sleep(1200 * time.Millisecond)
	defer fmt.Println("Thanks for playing!")

	// rudimentary select screen, player inputs game mode into os.stdin
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
			// not an easter egg
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

	// create the sdl window
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

	// make pixels array which is a flattened array given by winWidth and winHeight. Each pixel has
	// 4 bits of data hence the *4
	pixels := make([]byte, winWidth*winHeight*4)

	// create two paddles
	player1 := paddle{pos{50, 100}, 20, 100, 500, 0, 0, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, 500, 0, 1, color{255, 255, 255}}
	// create ball
	ball := ball{getCenter(), 20, 0, 0, color{255, 255, 255}}
	// initialize the balls starting direction and velocity
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

	// create keyState which checks for keypresses
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

		// handles different game states, play for ball in motion, start for pause
		switch state {
		case PLAY:
			// handle case for initial select screen
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
			// slowly progress balls xvelocity as match duration grows
			if ball.XVel > 0 {
				ball.XVel += gameElapsed / 50
			} else {
				ball.XVel -= gameElapsed / 50
			}
			// update ball position (checks for collisions)
			ball.update(&player1, &player2, elapsedTime)
			// check for escape key to pause game
			if keyState[sdl.SCANCODE_ESCAPE] != 0 {
				state = START
				paused = true
			}
		case START:
			// start screen, draw space text
			drawSpace(getCenter(), color{255, 255, 255}, 5, pixels)
			screenDraw(tex, renderer, frameStart, &elapsedTime, pixels)
			// check for space key to resume
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.Score == 9 || player2.Score == 9 {
					player1.Score = 0
					player2.Score = 0
				}
				// start game timer if fresh match, otherwise gametimer is paused
				if !paused {
					gameStart = time.Now()
				} else {
					paused = false
				}
				state = PLAY
			}
		// wip / unused, may use for future improved select screen
		case SELECT:
			var mode int
			fmt.Println("Select mode:")
			fmt.Printf("\t(0) Player vs Easy Computer\n\t(1) Player vs Hard Computer\n\t(2)Player vs Player\n")
			fmt.Println("Enter selection (#): ")
			fmt.Scanln("%d", &mode)
		}

		// clear pixels array to refresh screen
		clear(pixels)
		// draw new positions
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// update screen
		screenDraw(tex, renderer, frameStart, &elapsedTime, pixels)
	}
}
