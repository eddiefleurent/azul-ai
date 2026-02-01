# Build the Azul game
build:
    go build -o azul

# Run the game
run:
    go run .

# Clean build artifacts
clean:
    rm -f azul

# Build and run
start: build
    ./azul

# Play the game (build and run)
play: build
    ./azul

# Watch AI vs AI
play-ai: build
    ./azul -human 0 -players 2

# Play as human (you vs AI)
play-human: build
    ./azul -human 1 -players 2

# Play against hard AI (Terminator)
play-terminator: build
    ./azul -ai hard
