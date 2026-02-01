# Azul AI

A command-line implementation of the board game Azul with AI opponents.

## Building

```bash
go build -o azul-ai .
```

## Playing

```bash
# Default: 2 players, you vs medium AI
./azul-ai

# Play against hard AI
./azul-ai -ai hard

# 3 players, you're player 2
./azul-ai -players 3 -human 2

# Watch AI vs AI
./azul-ai -human 0
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-players N` | Number of players (2-4) | 2 |
| `-ai LEVEL` | AI difficulty: easy, medium, hard | medium |
| `-human N` | Which player is human (1-4), 0 for AI vs AI | 1 |
| `-help` | Show help | - |

## AI Difficulty Levels

- **Easy**: Random legal moves
- **Medium**: Heuristic-based (prioritizes completing lines, avoids overflow)
- **Hard**: Minimax with alpha-beta pruning (looks ahead 3-4 moves)

## Architecture

```
azul-ai/
├── main.go           # CLI and game loop
├── game/
│   ├── tiles.go      # Tile colors and utilities
│   ├── bag.go        # Tile bag with draw/discard
│   ├── factory.go    # Factory displays and center
│   ├── player.go     # Player board, pattern lines, wall
│   └── game.go       # Game state and rules
├── ai/
│   └── ai.go         # AI players (random, heuristic, minimax)
└── display/
    └── display.go    # Terminal rendering with colors
```

## Game Rules Summary

1. **Drafting**: Take all tiles of one color from a factory or center
2. **Placing**: Put tiles on a pattern line (1-5 capacity)
3. **Wall Tiling**: Completed lines move one tile to the wall, score points
4. **Game End**: First player to complete a horizontal wall row triggers end

### Scoring
- Tile placement: 1 + adjacent tiles
- Complete row: +2
- Complete column: +7
- All 5 of one color: +10
- Floor penalties: -1, -1, -2, -2, -2, -3, -3

## License

MIT
