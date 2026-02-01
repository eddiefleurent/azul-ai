# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Game Rules

For complete official Azul rules, see [docs/AZUL_RULES.md](docs/AZUL_RULES.md).

## Build and Run Commands

### Using just

```bash
just play           # Play the game (build and run, you vs medium AI)
just play-human     # Play as human (you vs AI)
just play-ai        # Watch AI vs AI
just play-terminator # Play against hard AI (Terminator)
just build          # Build the binary (outputs ./azul)
just run            # Run with go run (default settings)
just clean          # Remove build artifacts
just start          # Build and run the binary
```

### Direct go commands

```bash
# Build the binary
go build -o azul-ai .

# Run with default settings (2 players, you vs medium AI)
./azul-ai

# Run with specific AI difficulty
./azul-ai -ai easy     # Random moves
./azul-ai -ai medium   # Heuristic-based (default)
./azul-ai -ai hard     # Minimax with alpha-beta pruning

# Run with multiple players
./azul-ai -players 3 -human 2    # 3 players, you're player 2
./azul-ai -human 0               # Watch AI vs AI
```

## Architecture Overview

This is a command-line implementation of the board game Azul with AI opponents. The codebase follows a clean separation between game logic, AI, and display rendering.

### Package Structure

**game/** - Core game rules and state management
- `tiles.go` - TileColor enum (Blue/Yellow/Red/Black/White) and utilities
- `bag.go` - Tile bag that manages drawing and discarding tiles
- `factory.go` - Factory displays and center area that hold tiles between rounds
- `player.go` - PlayerBoard containing pattern lines (1-5 capacity), wall (5x5 grid), and floor line
- `game.go` - Main Game struct orchestrating all components, move validation, and round/game end logic

**ai/** - AI opponent implementations
- `ai.go` - Three difficulty levels:
  - Easy: Random legal moves
  - Medium: Heuristic-based (prioritizes completing lines, avoids overflow, considers bonus scoring)
  - Hard: Minimax with alpha-beta pruning (3-4 move lookahead)

**display/** - Terminal rendering with ANSI colors
- `display.go` - Renders game state, player boards, factories, and interactive move selection

**main.go** - CLI entry point and game loop with interactive human input

### Key Game State Flow

1. **Game Setup**: NewGame() creates factories (2*numPlayers + 1), initializes bag with 20 tiles per color
2. **Round Loop**:
   - Factories filled with 4 tiles each from bag
   - Players take turns selecting tiles from factories/center
   - When all factories/center empty, round ends
3. **Round End Scoring**:
   - Completed pattern lines move one tile to wall, rest discarded
   - Wall placement scores: 1 point + adjacent tiles (horizontal + vertical)
   - Floor line penalties applied: -1, -1, -2, -2, -2, -3, -3
4. **Game End**: Triggered when any player completes a horizontal wall row
   - Bonuses: +2 per complete row, +7 per complete column, +10 per complete color set

### Wall Pattern

The wall uses a standard Azul pattern where each row contains all 5 colors in a shifted arrangement:
```text
Row 1: Blue, Yellow, Red, Black, White
Row 2: White, Blue, Yellow, Red, Black
Row 3: Black, White, Blue, Yellow, Red
Row 4: Red, Black, White, Blue, Yellow
Row 5: Yellow, Red, Black, White, Blue
```

This means each color appears exactly once per row and column. Note: internally the code uses 0-based indexing, but all user-facing displays use 1-based numbering (lines 1-5).

### Important Validation Rules

- Pattern lines can only hold one color at a time
- A color cannot be placed on a pattern line if that color is already on the wall in that row
- Taking from center includes the first-player marker (goes to floor, -1 penalty)
- Overflow tiles from pattern lines go to floor line
- Game state is deeply cloneable for AI move simulation

### AI Heuristics (Medium Difficulty)

The medium AI uses weighted scoring:
- +50 for completing a pattern line, +10 per line size
- -15 per overflow tile
- -30 for placing directly on floor
- -20 for taking first-player marker early (when center has <5 tiles)
- +10-20 bonuses for moves contributing to row/column/color completion

### AI Minimax (Hard Difficulty)

- Depth 3-4 depending on move count (reduces to 3 when >20 moves available)
- Alpha-beta pruning for efficiency
- State evaluation weighs: player score (×10), pattern line progress, wall bonus potential, opponent scores (×-8), floor penalties (×-5)
- Maximizing when it's AI's turn, minimizing for opponents
