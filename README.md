# GoJump Console

A small terminal platformer written in Go using `tcell`.

## Project Description

GoJump Console is a simple obstacle-dodging game that runs in the terminal.  
You control a player block, jump over walls, and move right to avoid collisions.

## Setup And Run

### Requirements

- Go `1.25.0` (or compatible)
- A terminal that supports interactive key input

### Install Dependencies

```bash
go mod download
```

### Run

```bash
go run .
```

When the home screen appears:

- Press `S` to start

During the game:
- Press `J` to jump
- Press `D` to move right

At any time:
- Press `Esc` or `Ctrl+C` to quit

Notes:

- Optimal terminal size is `75x25`
- Changing terminal size can alter gameplay behavior and may break layout/flow

## Future Improvements

- Enforce or better handle dynamic terminal resizing
- Add score tracking
- Add difficulty scaling (speed and obstacle spacing)
- Improve movement/input handling for smoother  movement

## Citation

Parts of the screen setup/rendering structure and box-drawing approach were adapted from the tcell tutorial:

- https://github.com/gdamore/tcell/blob/main/TUTORIAL.md
