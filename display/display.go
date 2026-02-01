package display

import (
	"fmt"
	"strings"

	"github.com/eddie/azul-ai/game"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
	Blue      = "\033[94m"
	Yellow    = "\033[93m"
	Red       = "\033[91m"
	Black     = "\033[37m"
	White     = "\033[97m"
	Cyan      = "\033[96m"
	Magenta   = "\033[95m"
	Green     = "\033[92m"
	Gray      = "\033[90m"

	BgBlue   = "\033[44m"
	BgYellow = "\033[43m"
	BgRed    = "\033[41m"
	BgBlack  = "\033[40m"
	BgWhite  = "\033[47m"
)

// ColorTile returns a colored block representation of a tile
func ColorTile(t game.TileColor) string {
	switch t {
	case game.Blue:
		return BgBlue + White + Bold + " B " + Reset
	case game.Yellow:
		return BgYellow + Black + Bold + " Y " + Reset
	case game.Red:
		return BgRed + White + Bold + " R " + Reset
	case game.Black:
		return BgBlack + White + Bold + " K " + Reset
	case game.White:
		return BgWhite + Black + Bold + " W " + Reset
	case game.FirstPlayerMarker:
		return Magenta + Bold + "[1]" + Reset
	case game.NoTile:
		return Gray + " · " + Reset
	default:
		return " ? "
	}
}

// ColorTileSmall returns a smaller colored representation
func ColorTileSmall(t game.TileColor) string {
	switch t {
	case game.Blue:
		return Blue + Bold + "●" + Reset
	case game.Yellow:
		return Yellow + Bold + "●" + Reset
	case game.Red:
		return Red + Bold + "●" + Reset
	case game.Black:
		return Black + Bold + "●" + Reset
	case game.White:
		return White + Bold + "○" + Reset
	case game.FirstPlayerMarker:
		return Magenta + Bold + "①" + Reset
	case game.NoTile:
		return Gray + "·" + Reset
	default:
		return "?"
	}
}

// ColorTileCompact returns a compact colored tile (2 chars visible)
func ColorTileCompact(t game.TileColor) string {
	switch t {
	case game.Blue:
		return BgBlue + White + Bold + "B " + Reset
	case game.Yellow:
		return BgYellow + Black + Bold + "Y " + Reset
	case game.Red:
		return BgRed + White + Bold + "R " + Reset
	case game.Black:
		return BgBlack + White + Bold + "K " + Reset
	case game.White:
		return BgWhite + Black + Bold + "W " + Reset
	case game.FirstPlayerMarker:
		return Magenta + Bold + "1 " + Reset
	case game.NoTile:
		return Gray + "· " + Reset
	default:
		return "? "
	}
}

// DimTileCompact returns a compact dimmed tile for empty wall slots (2 chars)
func DimTileCompact(t game.TileColor) string {
	switch t {
	case game.Blue:
		return Blue + Dim + "B " + Reset
	case game.Yellow:
		return Yellow + Dim + "Y " + Reset
	case game.Red:
		return Red + Dim + "R " + Reset
	case game.Black:
		return Gray + Dim + "K " + Reset
	case game.White:
		return White + Dim + "W " + Reset
	default:
		return "· "
	}
}

// DimTile returns a dimmed tile for empty wall slots (3 chars, matches ColorTile)
func DimTile(t game.TileColor) string {
	switch t {
	case game.Blue:
		return Blue + Dim + " B " + Reset
	case game.Yellow:
		return Yellow + Dim + " Y " + Reset
	case game.Red:
		return Red + Dim + " R " + Reset
	case game.Black:
		return Gray + Dim + " K " + Reset
	case game.White:
		return White + Dim + " W " + Reset
	default:
		return " · "
	}
}

// Box width constant for alignment
const boxWidth = 45

// RenderFactories displays all factories and center
func RenderFactories(g *game.Game) string {
	var sb strings.Builder

	// Box header
	sb.WriteString(Bold + Cyan + "╭" + strings.Repeat("─", boxWidth) + "╮" + Reset + "\n")
	title := "FACTORY DISPLAYS"
	padding := (boxWidth - len(title)) / 2
	sb.WriteString(Bold + Cyan + "│" + Reset + strings.Repeat(" ", padding) + Bold + title + Reset + strings.Repeat(" ", boxWidth-padding-len(title)) + Bold + Cyan + "│" + Reset + "\n")
	sb.WriteString(Bold + Cyan + "╰" + strings.Repeat("─", boxWidth) + "╯" + Reset + "\n\n")

	for i, f := range g.Factories {
		sb.WriteString(fmt.Sprintf("  %s[%d]%s ", Cyan, i+1, Reset))
		if f.IsEmpty() {
			sb.WriteString(Gray + "─ empty ─" + Reset)
		} else {
			for _, t := range f.Tiles {
				sb.WriteString(ColorTile(t))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s[C]%s ", Magenta+Bold, Reset))
	if g.Center.IsEmpty() && !g.Center.HasFirstPlayerTile {
		sb.WriteString(Gray + "─ empty ─" + Reset)
	} else {
		if g.Center.HasFirstPlayerTile {
			sb.WriteString(ColorTile(game.FirstPlayerMarker))
		}
		for _, t := range g.Center.Tiles {
			sb.WriteString(ColorTile(t))
		}
	}
	sb.WriteString(Dim + "  (center)" + Reset + "\n")

	return sb.String()
}

// RenderPlayerBoard displays a player's board
func RenderPlayerBoard(pb *game.PlayerBoard, playerNum int, isCurrentPlayer bool, playerName string) string {
	var sb strings.Builder

	// Header box
	borderColor := ""
	if isCurrentPlayer {
		borderColor = Green
	}

	sb.WriteString(Bold + borderColor + "╭" + strings.Repeat("─", boxWidth) + "╮" + Reset + "\n")

	// Player name and score line
	turnMarker := ""
	if isCurrentPlayer {
		turnMarker = " ◄ YOUR TURN"
	}
	content := fmt.Sprintf(" %-12s          Score: %3d%s", playerName, pb.Score, turnMarker)
	// Pad to box width
	if len(content) < boxWidth {
		content += strings.Repeat(" ", boxWidth-len(content))
	}
	sb.WriteString(Bold + borderColor + "│" + Reset + content + Bold + borderColor + "│" + Reset + "\n")
	sb.WriteString(Bold + borderColor + "╰" + strings.Repeat("─", boxWidth) + "╯" + Reset + "\n")

	// Pattern lines and wall header with yellow border
	sb.WriteString("\n")
	sb.WriteString("  " + Yellow + Bold + "┌─ LINES & WALL ────────────────────────┐" + Reset + "\n")

	for row := 0; row < 5; row++ {
		pl := pb.PatternLines[row]

		// Yellow left border
		sb.WriteString("  " + Yellow + "│" + Reset)

		// Row number: " 1│" = 3 chars
		sb.WriteString(fmt.Sprintf(" %s%d%s│", Cyan+Bold, row+1, Reset))

		// Pattern line (right-aligned) - each slot is 2 chars wide, max 10 chars
		padding := 5 - pl.Size
		sb.WriteString(strings.Repeat("  ", padding))

		for i := 0; i < pl.Size; i++ {
			if i < pl.Size-pl.Filled {
				// Empty slot
				sb.WriteString(White + "□ " + Reset)
			} else {
				sb.WriteString(ColorTileCompact(pl.Color))
			}
		}

		// Arrow: " → " = 3 chars
		if pl.IsFull() {
			sb.WriteString(Green + Bold + " ▶ " + Reset)
		} else {
			sb.WriteString(Gray + " → " + Reset)
		}
		// Separator: "│ " = 2 chars
		sb.WriteString("│ ")

		// Wall row
		for col := 0; col < 5; col++ {
			expectedColor := game.WallPattern[row][col]
			if pb.Wall[row][col] {
				sb.WriteString(ColorTileCompact(expectedColor))
			} else {
				sb.WriteString(DimTileCompact(expectedColor))
			}
		}

		// Yellow right border
		sb.WriteString(" " + Yellow + "│" + Reset + "\n")
	}

	// Yellow bottom border
	sb.WriteString("  " + Yellow + Bold + "└──────────────────────────────────────┘" + Reset + "\n")

	// Floor line
	sb.WriteString("\n  " + Red + "Floor:" + Reset + " ")
	if len(pb.FloorLine) == 0 {
		sb.WriteString(Gray + "empty" + Reset)
	} else {
		for i, t := range pb.FloorLine {
			sb.WriteString(ColorTileSmall(t))
			if i < len(game.FloorPenalties) {
				sb.WriteString(fmt.Sprintf("%s%d%s ", Red+Dim, game.FloorPenalties[i], Reset))
			} else {
				sb.WriteString(" ")
			}
		}
	}

	// Floor penalty slots
	sb.WriteString("\n         ")
	for i := 0; i < 7; i++ {
		if i < len(pb.FloorLine) {
			sb.WriteString(Red + "▼  " + Reset)
		} else {
			sb.WriteString(Gray + fmt.Sprintf("%d  ", game.FloorPenalties[i]) + Reset)
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

// RenderGame displays the full game state
func RenderGame(g *game.Game, playerNames []string) string {
	var sb strings.Builder

	// Clear screen
	sb.WriteString("\033[H\033[2J")

	// Title - same width as other boxes
	sb.WriteString("\n")
	sb.WriteString(Bold + Cyan + "╔" + strings.Repeat("═", boxWidth) + "╗" + Reset + "\n")
	title := "A Z U L"
	padding := (boxWidth - len(title)) / 2
	sb.WriteString(Bold + Cyan + "║" + Reset + strings.Repeat(" ", padding) + Bold + title + Reset + strings.Repeat(" ", boxWidth-padding-len(title)) + Bold + Cyan + "║" + Reset + "\n")
	roundStr := fmt.Sprintf("Round %d", g.Round)
	padding = (boxWidth - len(roundStr)) / 2
	sb.WriteString(Bold + Cyan + "║" + Reset + Dim + strings.Repeat(" ", padding) + roundStr + strings.Repeat(" ", boxWidth-padding-len(roundStr)) + Reset + Bold + Cyan + "║" + Reset + "\n")
	sb.WriteString(Bold + Cyan + "╚" + strings.Repeat("═", boxWidth) + "╝" + Reset + "\n")
	sb.WriteString("\n")

	// Factories
	sb.WriteString(RenderFactories(g))
	sb.WriteString("\n")

	// All player boards
	for i, player := range g.Players {
		name := fmt.Sprintf("Player %d", i+1)
		if i < len(playerNames) && playerNames[i] != "" {
			name = playerNames[i]
		}
		sb.WriteString(RenderPlayerBoard(player, i+1, i == g.CurrentPlayer, name))
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderGameOver displays final results
func RenderGameOver(g *game.Game, playerNames []string) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(Bold + Magenta + "╔" + strings.Repeat("═", boxWidth) + "╗" + Reset + "\n")
	title := "GAME OVER"
	padding := (boxWidth - len(title)) / 2
	sb.WriteString(Bold + Magenta + "║" + Reset + strings.Repeat(" ", padding) + Bold + title + Reset + strings.Repeat(" ", boxWidth-padding-len(title)) + Bold + Magenta + "║" + Reset + "\n")
	sb.WriteString(Bold + Magenta + "╚" + strings.Repeat("═", boxWidth) + "╝" + Reset + "\n")
	sb.WriteString("\n")

	sb.WriteString(Bold + "Final Scores:" + Reset + "\n\n")

	winner := g.GetWinner()
	for i, player := range g.Players {
		name := fmt.Sprintf("Player %d", i+1)
		if i < len(playerNames) && playerNames[i] != "" {
			name = playerNames[i]
		}

		marker := "  "
		if i == winner {
			marker = Green + "★ " + Reset
		}
		sb.WriteString(fmt.Sprintf("%s%-15s %s%3d%s points\n", marker, name, Bold, player.Score, Reset))
	}

	if winner >= 0 {
		winnerName := fmt.Sprintf("Player %d", winner+1)
		if winner < len(playerNames) && playerNames[winner] != "" {
			winnerName = playerNames[winner]
		}
		sb.WriteString(fmt.Sprintf("\n%s🎉 %s wins! 🎉%s\n", Bold+Green, winnerName, Reset))
	}

	return sb.String()
}

// ColorLegend shows what each color code means
func ColorLegend() string {
	var sb strings.Builder

	sb.WriteString("\n" + Dim + "Tiles: " + Reset)
	sb.WriteString(ColorTile(game.Blue) + "Blue ")
	sb.WriteString(ColorTile(game.Yellow) + "Yellow ")
	sb.WriteString(ColorTile(game.Red) + "Red ")
	sb.WriteString(ColorTile(game.Black) + "Black ")
	sb.WriteString(ColorTile(game.White) + "White ")
	sb.WriteString(ColorTile(game.FirstPlayerMarker) + "1st\n")

	return sb.String()
}

// SourceOption represents a pickable source
type SourceOption struct {
	Index  int
	Colors []game.TileColor
	Label  string
}

// GetAvailableSources returns sources that have tiles
func GetAvailableSources(g *game.Game) []SourceOption {
	sources := make([]SourceOption, 0)

	for i, f := range g.Factories {
		if !f.IsEmpty() {
			sources = append(sources, SourceOption{
				Index:  i,
				Colors: f.GetColors(),
				Label:  fmt.Sprintf("Factory %d", i+1),
			})
		}
	}

	if !g.Center.IsEmpty() || g.Center.HasFirstPlayerTile {
		sources = append(sources, SourceOption{
			Index:  -1,
			Colors: g.Center.GetColors(),
			Label:  "Center",
		})
	}

	return sources
}

// RenderSourceSelection shows available sources to pick from
func RenderSourceSelection(g *game.Game) string {
	var sb strings.Builder

	sb.WriteString("\n" + Bold + "Step 1: " + Reset + "Choose where to take tiles from:\n\n")

	sources := GetAvailableSources(g)
	for i, src := range sources {
		sb.WriteString(fmt.Sprintf("  %s[%d]%s %-12s", Cyan+Bold, i+1, Reset, src.Label))

		if src.Index == -1 {
			if g.Center.HasFirstPlayerTile {
				sb.WriteString(ColorTile(game.FirstPlayerMarker))
			}
			for _, t := range g.Center.Tiles {
				sb.WriteString(ColorTile(t))
			}
		} else {
			for _, t := range g.Factories[src.Index].Tiles {
				sb.WriteString(ColorTile(t))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderColorSelection shows available colors at a source
func RenderColorSelection(g *game.Game, sourceIdx int) string {
	var sb strings.Builder

	var colors []game.TileColor
	var sourceName string

	if sourceIdx == -1 {
		colors = g.Center.GetColors()
		sourceName = "Center"
	} else {
		colors = g.Factories[sourceIdx].GetColors()
		sourceName = fmt.Sprintf("Factory %d", sourceIdx+1)
	}

	sb.WriteString(fmt.Sprintf("\n"+Bold+"Step 2: "+Reset+"Choose which color to take from %s:\n\n", sourceName))

	for i, color := range colors {
		count := 0
		if sourceIdx == -1 {
			for _, t := range g.Center.Tiles {
				if t == color {
					count++
				}
			}
		} else {
			for _, t := range g.Factories[sourceIdx].Tiles {
				if t == color {
					count++
				}
			}
		}

		sb.WriteString(fmt.Sprintf("  %s[%d]%s ", Cyan+Bold, i+1, Reset))
		sb.WriteString(ColorTile(color))
		sb.WriteString(fmt.Sprintf(" × %d\n", count))
	}

	return sb.String()
}

// RenderLineSelection shows available pattern lines to place tiles with inline previews
func RenderLineSelection(pb *game.PlayerBoard, color game.TileColor, tileCount int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n"+Bold+"Step 3: "+Reset+"Place %d %s tile(s) on which line?\n\n", tileCount, color.FullName()))

	validLines := pb.GetValidPlacements(color)
	optNum := 1

	for _, lineIdx := range validLines {
		if lineIdx == -1 {
			continue
		}

		pl := pb.PatternLines[lineIdx]
		space := pl.Size - pl.Filled
		overflow := 0
		if tileCount > space {
			overflow = tileCount - space
		}
		tilesToPlace := min(tileCount, space)
		newFilled := pl.Filled + tilesToPlace
		willComplete := newFilled == pl.Size

		// Option number
		sb.WriteString(fmt.Sprintf("  %s[%d]%s ", Cyan+Bold, optNum, Reset))

		// Show pattern line AFTER placing tiles
		// Right-align: padding for smaller lines, each tile is 3 chars
		padding := 5 - pl.Size
		sb.WriteString(strings.Repeat("   ", padding))

		// Show tiles filling from LEFT to RIGHT (filled tiles on left, empty on right)
		// This makes it clear how many more you need to fill
		for i := 0; i < pl.Size; i++ {
			if i < newFilled {
				// Filled tiles (existing + new) - show with color background
				sb.WriteString(ColorTile(color))
			} else {
				// Still empty after placement - show empty box
				sb.WriteString(Gray + " □ " + Reset)
			}
		}

		// Arrow
		if willComplete {
			sb.WriteString(Green + Bold + " ▶ " + Reset)
		} else {
			sb.WriteString(Gray + " → " + Reset)
		}

		// Show wall row with the target slot highlighted (all 3-char tiles for alignment)
		sb.WriteString("│ ")
		for col := 0; col < 5; col++ {
			expectedColor := game.WallPattern[lineIdx][col]
			if pb.Wall[lineIdx][col] {
				// Already placed on wall
				sb.WriteString(ColorTile(expectedColor))
			} else if willComplete && expectedColor == color {
				// This slot WILL be filled - highlight it with full color
				sb.WriteString(ColorTile(expectedColor))
			} else {
				// Empty wall slot - dim but same width
				sb.WriteString(DimTile(expectedColor))
			}
		}

		// Completion and overflow indicators
		if willComplete {
			sb.WriteString(Green + " ✓" + Reset)
		}
		if overflow > 0 {
			sb.WriteString(fmt.Sprintf(Red+" +%d floor"+Reset, overflow))
		}

		sb.WriteString("\n")
		optNum++
	}

	// Floor option
	penalty := 0
	floorLen := len(pb.FloorLine)
	for i := 0; i < tileCount && floorLen+i < len(game.FloorPenalties); i++ {
		penalty += game.FloorPenalties[floorLen+i]
	}
	sb.WriteString(fmt.Sprintf("\n  %s[%d]%s Floor: all to floor "+Red+"(%d points)"+Reset+"\n", Red+Bold, optNum, Reset, penalty))

	return sb.String()
}

// RenderBoardPreview shows the player's board with a preview of tiles being placed
func RenderBoardPreview(pb *game.PlayerBoard, color game.TileColor, tileCount int, targetLine int) string {
	var sb strings.Builder

	sb.WriteString("\n" + Yellow + Bold + "  ┌─ PREVIEW ─────────────────────────────┐" + Reset + "\n")

	for row := 0; row < 5; row++ {
		pl := pb.PatternLines[row]

		// Determine if this row is being previewed
		isPreviewRow := row == targetLine
		previewFilled := pl.Filled
		previewColor := pl.Color
		overflow := 0

		if isPreviewRow {
			space := pl.Size - pl.Filled
			tilesToPlace := min(tileCount, space)
			if tileCount > space {
				overflow = tileCount - space
			}
			previewFilled = pl.Filled + tilesToPlace
			previewColor = color
		}

		// Row number - highlight if preview row
		if isPreviewRow {
			sb.WriteString(fmt.Sprintf("  %s│%s %s%d%s│", Yellow, Reset, Yellow+Bold, row+1, Reset))
		} else {
			sb.WriteString(fmt.Sprintf("  %s│%s %s%d%s│", Yellow, Reset, Cyan+Bold, row+1, Reset))
		}

		// Pattern line (right-aligned) - use 3-char tiles for preview row
		padding := 5 - pl.Size
		if isPreviewRow {
			// Use full colored tiles (3 chars each) for preview
			sb.WriteString(strings.Repeat("   ", padding))
			for i := 0; i < pl.Size; i++ {
				if i < pl.Size-previewFilled {
					sb.WriteString(Gray + " · " + Reset)
				} else if i < pl.Size-pl.Filled {
					// New tiles being placed - show with full background color
					sb.WriteString(ColorTile(previewColor))
				} else {
					sb.WriteString(ColorTile(previewColor))
				}
			}
		} else {
			// Normal rows use 2-char compact tiles
			sb.WriteString(strings.Repeat("  ", padding))
			sb.WriteString(" ") // Extra space to align with 3-char tiles
			for i := 0; i < pl.Size; i++ {
				if i < pl.Size-pl.Filled {
					sb.WriteString(White + "□ " + Reset)
				} else {
					sb.WriteString(ColorTileCompact(pl.Color))
				}
			}
		}

		// Arrow and wall
		willComplete := isPreviewRow && previewFilled == pl.Size
		if isPreviewRow {
			if willComplete {
				sb.WriteString(Green + Bold + " ▶ " + Reset)
			} else {
				sb.WriteString(Yellow + " → " + Reset)
			}
		} else {
			sb.WriteString(Gray + " → " + Reset)
		}

		sb.WriteString("│ ")

		// Wall row
		for col := 0; col < 5; col++ {
			expectedColor := game.WallPattern[row][col]
			if pb.Wall[row][col] {
				sb.WriteString(ColorTileCompact(expectedColor))
			} else {
				sb.WriteString(DimTileCompact(expectedColor))
			}
		}

		// Show overflow indicator on the preview row
		if isPreviewRow && overflow > 0 {
			sb.WriteString(fmt.Sprintf(Red + " +" + fmt.Sprint(overflow) + " floor" + Reset))
		}
		if willComplete {
			sb.WriteString(Green + " ✓" + Reset)
		}

		sb.WriteString(Yellow + " │" + Reset + "\n")
	}

	sb.WriteString(Yellow + Bold + "  └──────────────────────────────────────────┘" + Reset + "\n")

	return sb.String()
}
