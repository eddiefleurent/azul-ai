package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eddie/azul-ai/ai"
	"github.com/eddie/azul-ai/display"
	"github.com/eddie/azul-ai/game"
)

func main() {
	// Command line flags
	numPlayers := flag.Int("players", 2, "Number of players (2-4)")
	aiDifficulty := flag.String("ai", "medium", "AI difficulty: easy, medium, hard")
	humanPlayer := flag.Int("human", 1, "Which player is human (1-4), 0 for AI vs AI")
	showHelp := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	// Parse AI difficulty
	var difficulty ai.Difficulty
	switch strings.ToLower(*aiDifficulty) {
	case "easy":
		difficulty = ai.Easy
	case "medium":
		difficulty = ai.Medium
	case "hard":
		difficulty = ai.Hard
	default:
		difficulty = ai.Medium
	}

	// Create game
	g := game.NewGame(*numPlayers)

	// Use the game's clamped player count (NewGame clamps to 2-4)
	numPlayersActual := g.NumPlayers

	// Player names
	playerNames := make([]string, numPlayersActual)
	aiPlayers := make(map[int]*ai.AIPlayer)

	for i := 0; i < numPlayersActual; i++ {
		if i+1 == *humanPlayer {
			playerNames[i] = "You"
		} else {
			aiPlayers[i] = ai.NewAIPlayer(difficulty, i)
			playerNames[i] = aiPlayers[i].Name()
		}
	}

	reader := bufio.NewReader(os.Stdin)

	// Main game loop
	for !g.GameOver {
		moves := g.GetValidMoves()
		if len(moves) == 0 {
			fmt.Print(display.RenderGame(g, playerNames))
			fmt.Println("No valid moves available!")
			break
		}

		var selectedMove game.Move

		if aiPlayer, isAI := aiPlayers[g.CurrentPlayer]; isAI {
			// AI's turn - show game state
			fmt.Print(display.RenderGame(g, playerNames))
			fmt.Printf("\n%s is thinking...\n", aiPlayer.Name())
			selectedMove = aiPlayer.ChooseMove(g, moves)
			fmt.Printf("%s chose: %s\n", aiPlayer.Name(), selectedMove.String())
			fmt.Println("\nPress Enter to continue...")
			reader.ReadString('\n')
		} else {
			// Human's turn - interactive selection (shows game state internally)
			selectedMove = getHumanMoveInteractive(reader, g, playerNames)
		}

		// Apply the move
		err := g.ApplyMove(selectedMove)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			reader.ReadString('\n')
		}
	}

	// Game over
	fmt.Print(display.RenderGameOver(g, playerNames))
}

func getHumanMoveInteractive(reader *bufio.Reader, g *game.Game, playerNames []string) game.Move {
	player := g.Players[g.CurrentPlayer]

	// Step 1: Choose source - show full game state first
	var sourceIdx int
	for {
		// Show full game state
		fmt.Print(display.RenderGame(g, playerNames))
		fmt.Print(display.RenderSourceSelection(g))
		fmt.Print("\n  Enter number (or 'q' to quit, 'h' for help): ")

		input := readInput(reader)
		if handleSpecialInput(input) {
			continue
		}

		sources := display.GetAvailableSources(g)
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(sources) {
			fmt.Printf("\n  %sInvalid choice. Enter 1-%d%s\n", display.Red, len(sources), display.Reset)
			waitForEnter(reader)
			continue
		}

		sourceIdx = sources[num-1].Index
		break
	}

	// Step 2: Choose color - show full game state again
	var selectedColor game.TileColor
	var tileCount int
	for {
		// Show full game state
		fmt.Print(display.RenderGame(g, playerNames))
		fmt.Print(display.RenderColorSelection(g, sourceIdx))
		fmt.Print("\n  Enter number (or 'b' to go back): ")

		input := readInput(reader)
		if input == "b" || input == "back" {
			return getHumanMoveInteractive(reader, g, playerNames) // Start over
		}
		if handleSpecialInput(input) {
			continue
		}

		var colors []game.TileColor
		if sourceIdx == -1 {
			colors = g.Center.GetColors()
		} else {
			colors = g.Factories[sourceIdx].GetColors()
		}

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(colors) {
			fmt.Printf("\n  %sInvalid choice. Enter 1-%d%s\n", display.Red, len(colors), display.Reset)
			waitForEnter(reader)
			continue
		}

		selectedColor = colors[num-1]

		// Count tiles
		if sourceIdx == -1 {
			for _, t := range g.Center.Tiles {
				if t == selectedColor {
					tileCount++
				}
			}
		} else {
			for _, t := range g.Factories[sourceIdx].Tiles {
				if t == selectedColor {
					tileCount++
				}
			}
		}
		break
	}

	// Step 3: Choose destination line - previews are shown inline
	var lineIdx int
	validLines := player.GetValidPlacements(selectedColor)

	for {
		// Show full game state
		fmt.Print(display.RenderGame(g, playerNames))

		// Show line options with inline previews
		fmt.Print(display.RenderLineSelection(player, selectedColor, tileCount))

		fmt.Print("\n  Enter number (or 'b' to go back): ")

		input := readInput(reader)
		if input == "b" || input == "back" {
			return getHumanMoveInteractive(reader, g, playerNames) // Start over
		}
		if handleSpecialInput(input) {
			continue
		}

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(validLines) {
			fmt.Printf("\n  %sInvalid choice. Enter 1-%d%s\n", display.Red, len(validLines), display.Reset)
			waitForEnter(reader)
			continue
		}

		lineIdx = validLines[num-1]
		break
	}

	return game.Move{
		FactoryIdx: sourceIdx,
		Color:      selectedColor,
		LineIdx:    lineIdx,
	}
}

func waitForEnter(reader *bufio.Reader) {
	fmt.Print("  Press Enter to continue...")
	reader.ReadString('\n')
}

func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func handleSpecialInput(input string) bool {
	switch input {
	case "q", "quit":
		fmt.Println("\nThanks for playing!")
		os.Exit(0)
	case "h", "help":
		printHelp()
		return true
	}
	return false
}

func printHelp() {
	help := `
╔═══════════════════════════════════════════════════════════════════════════════╗
║                              AZUL - HELP                                      ║
╚═══════════════════════════════════════════════════════════════════════════════╝

` + display.Bold + `OBJECTIVE:` + display.Reset + `
  Score the most points by strategically placing tiles on your wall.

` + display.Bold + `HOW TO PLAY:` + display.Reset + `
  Each turn has 3 steps:
  1. Pick a source (factory or center)
  2. Pick a color from that source (you take ALL tiles of that color)
  3. Pick a pattern line to place them on

` + display.Bold + `PATTERN LINES:` + display.Reset + `
  - The 5 rows on the left (sizes 1-5)
  - Each line can only hold ONE color
  - When a line is full, one tile moves to the wall at round end

` + display.Bold + `WALL:` + display.Reset + `
  - The 5x5 grid on the right
  - Each row can only have each color once (shown as [b][y][r][k][w])
  - Tiles score points when placed based on adjacent tiles

` + display.Bold + `FLOOR:` + display.Reset + `
  - Overflow tiles and first-player marker go here
  - Costs: -1, -1, -2, -2, -2, -3, -3 points

` + display.Bold + `SCORING:` + display.Reset + `
  - Tile placement: 1 point + adjacent tiles in row/column
  - Complete horizontal line: +2 bonus
  - Complete vertical line: +7 bonus
  - All 5 of one color: +10 bonus

` + display.Bold + `GAME END:` + display.Reset + `
  The game ends when any player completes a horizontal wall row.

` + display.Bold + `CONTROLS:` + display.Reset + `
  - Enter a number to select
  - 'b' to go back a step
  - 'h' for this help
  - 'q' to quit

` + display.Bold + `COMMAND LINE OPTIONS:` + display.Reset + `
  -players N    Number of players (2-4), default 2
  -ai LEVEL     AI difficulty: easy, medium, hard (default medium)
  -human N      Which player is human (1-4), 0 for AI vs AI
  -help         Show this help

`
	fmt.Println(help)
}
