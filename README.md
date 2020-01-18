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

 #### Controls

 Game:
 - ESC to pause
 - Space to unpause or resume
 - First to 10 points wins and the game automatically resets

 Player 1 (left paddle):
 - W/S to move paddle up or down

 Player 2 (right paddle):
 - Up/down arrow keys to move paddle up or down (or I/K)

 #### Planned features
 - Win screen
 - Improving AI difficulties
 - Improve select screen
 - OSX binary for release - you should be able to self compile currently; I just
 haven't figured out the cross compiling yet

 #### Known bugs
 - Paddles can move off screen
    - This actually provides some utility for changing ball velocity, but it may
    make sense to set a limit to how far off they can move
- Collision detection could use some tweaking
- AI sucks
