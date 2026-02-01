# Azul Official Rules

Based on the official rulebook from Plan B Games / Next Move Games.

## Overview

Azul is a tile-laying strategy game for 2-4 players designed by Michael Kiesling. Players take turns selecting colored tiles from factory displays to complete patterns on their player boards and score points.

## Components

- **100 tiles** (20 of each color: Blue, Yellow, Red, Black, White)
- **Player boards** with pattern lines, wall, floor line, and score track
- **Factory displays** (circular discs)
- **1 Starting player marker**
- **Tile bag**

## Setup

1. Each player takes a player board (colored wall side up)
2. Place scoring marker on "0" of score track
3. Place Factory displays in a circle:
   - **2 players**: 5 Factory displays
   - **3 players**: 7 Factory displays
   - **4 players**: 9 Factory displays
4. Fill the bag with all 100 tiles
5. The starting player fills each Factory display with exactly **4 tiles** randomly drawn from the bag
6. Place the starting player marker in the center of the table

## Objective

Score the most points by the end of the game. The game ends after the round in which at least one player completes a horizontal line of 5 consecutive tiles on their wall.

## Gameplay

Each round consists of three phases:

### Phase 1: Factory Offer

The starting player places the starting player marker in the center of the table, then takes the first turn. Play proceeds clockwise.

**On your turn, you MUST pick tiles in one of two ways:**

#### Option A: Take from a Factory Display
- Pick **ALL tiles of the same color** from any one Factory display
- Move the **remaining tiles** from that Factory display to the center of the table

#### Option B: Take from the Center
- Pick **ALL tiles of one color** from the center of the table
- **IMPORTANT**: If you are the **first player this round** to pick from the center, you must also take the starting player marker and place it on the **leftmost free space** of your floor line

**After picking tiles, place them on your player board:**

1. Add tiles to **one of the 5 pattern lines** (row 1 holds 1 tile, row 5 holds 5 tiles)
2. Place tiles **from right to left** in your chosen pattern line
3. If the pattern line already holds tiles, you may only add tiles of **the same color**
4. Once all spaces of a pattern line are filled, that line is **complete**
5. **Excess tiles** that don't fit must go to the **floor line**

**Critical Rule**: You are **NOT allowed** to place tiles of a certain color in a pattern line whose corresponding wall row **already has a tile of that color**.

#### Floor Line
- Any tiles you cannot or do not want to place go to the floor line (filled left to right)
- These tiles give **negative points** during scoring
- If the floor line is full (7 spaces), excess tiles are discarded to the game box

**Phase 1 ends** when the center of the table AND all Factory displays contain no more tiles.

### Phase 2: Wall-Tiling

All players can do this simultaneously:

1. Go through your pattern lines **from top to bottom**
2. For each **complete** pattern line:
   - Move the **rightmost tile** to the matching color space on the corresponding wall row
   - Score points **immediately** (see Scoring)
   - Remove remaining tiles from that pattern line and discard them to the game box
3. **Incomplete** pattern lines stay on your board for the next round

#### Scoring

**When placing a tile on the wall:**

- **If no adjacent tiles** (horizontally or vertically): gain **1 point**
- **If there are adjacent tiles**:
  - Count all **horizontally linked tiles** (including the new one) = X points
  - Count all **vertically linked tiles** (including the new one) = Y points
  - **Total = X + Y** (both directions score separately)

**Examples:**
- Tile with 2 tiles to its left horizontally = 3 points (horizontal only)
- Tile with 2 tiles above it vertically = 3 points (vertical only)
- Tile with 3 horizontal + 2 vertical adjacent = 4 + 3 = 7 points

**Floor Line Penalties:**
At the end of the Wall-tiling phase, lose points for tiles in your floor line:

| Position | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
|----------|---|---|---|---|---|---|---|
| Penalty  | -1| -1| -2| -2| -2| -3| -3|

- Score cannot drop below 0
- Remove all tiles from floor line and discard them
- If you have the starting player marker in your floor line, it counts as a tile for penalties, but keep it in front of you (don't discard)

### Phase 3: Preparing the Next Round

If no player has completed a horizontal wall row:

1. The player with the starting player marker refills each Factory display with 4 tiles from the bag
2. If the bag is empty, refill it with discarded tiles and continue
3. Place the starting player marker back in the center
4. Start the new round

## End of Game

The game ends **immediately after the Wall-tiling phase** in which at least one player completes a horizontal line of 5 consecutive tiles on their wall.

### Final Scoring Bonuses

After the final round, score additional points:

| Achievement | Bonus |
|-------------|-------|
| Each complete **horizontal line** (5 tiles in a row) | **+2 points** |
| Each complete **vertical line** (5 tiles in a column) | **+7 points** |
| Each complete **color set** (all 5 tiles of one color on wall) | **+10 points** |

### Determining the Winner

The player with the **most points** wins.

**Tiebreaker**: The tied player with **more complete horizontal lines** wins. If still tied, players share the victory.

## Wall Pattern (Colored Side)

The standard wall has a fixed color pattern where each row is shifted:

```
Row 1: Blue    Yellow  Red     Black   White
Row 2: White   Blue    Yellow  Red     Black
Row 3: Black   White   Blue    Yellow  Red
Row 4: Red     Black   White   Blue    Yellow
Row 5: Yellow  Red     Black   White   Blue
```

Each color appears exactly once per row and once per column.

## Variant: Gray Wall

Use the gray side of the player board for a different experience:
- When moving a tile from a pattern line to the wall, you may place it on **any space** in the corresponding wall row
- However, each color may only appear **once per row** and **once per column**
- If you cannot legally place a tile, all tiles from that pattern line go to the floor line

---

## Quick Reference

### Factory Count
- 2 players: 5 factories
- 3 players: 7 factories  
- 4 players: 9 factories

### Tiles per Color: 20 (100 total)

### Pattern Line Sizes
- Line 1: 1 tile
- Line 2: 2 tiles
- Line 3: 3 tiles
- Line 4: 4 tiles
- Line 5: 5 tiles

### Floor Penalties: -1, -1, -2, -2, -2, -3, -3

### End Game Bonuses
- Complete row: +2
- Complete column: +7
- Complete color: +10
