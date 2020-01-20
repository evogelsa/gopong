## GoPong

This was written alongside the Games With Go Pong tutorial created by Jack Mott.
I used this as a way to learn and experiment with Golang, and I have expanded
the Pong game to add some new features such as:
- 2 player mode
- Paddle edge detection
- Changing ball velocity
- Random start state
- Unpause screen
- Game pause with escape

#### Gameplay

The game consists of two paddles and a ball. The goal is to prevent the ball from
hitting your side of the screen. Paddles can give velocity to the ball based on
where the ball contacts the paddle. There are seven contactable segments of the
paddle. When the ball detects a collision it calculates its new velocity in order
to give it a predetermined angle. The angle of the ball will either be 75, 45, 15,
or 0 degrees depending on which segment it contacts. After the 4th and 12th
collisions the ball will speed up by 50%. 

#### Controls

 Game:
 - ESC to pause
 - Space to unpause or resume
 - First to 9 points wins and the game automatically resets

 Player 1 (left paddle):
 - W/S to move paddle up or down

 Player 2 (right paddle):
 - Up/down arrow keys to move paddle up or down (or I/K)

#### Planned features
 - Win screen
 - Improved difficulty/game mode select screen
 - Improving AI difficulties
 - OSX binary for release - you should be able to self compile currently; I just
 haven't figured out the cross compiling yet

#### Known bugs
- Collision detection could use some tweaking
