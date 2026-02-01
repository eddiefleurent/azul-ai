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
