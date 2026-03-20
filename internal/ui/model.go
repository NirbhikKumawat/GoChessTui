package ui

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/NirbhikKumawat/GoChess/chess"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type EnterMode int

const (
	SelectMode EnterMode = iota
	MoveMode
)

var pieceMap = map[uint8]map[uint8]string{
	chess.White: {chess.Pawn: "♟", chess.Knight: "♞", chess.Bishop: "♝", chess.Rook: "♜", chess.Queen: "♛", chess.King: "♚", chess.Empty: "  "},
	chess.Black: {chess.Pawn: "♙", chess.Knight: "♘", chess.Bishop: "♗", chess.Rook: "♖", chess.Queen: "♕", chess.King: "♔", chess.Empty: "  "},
}

var (
	bgValidMove = lipgloss.Color("238") // Dark gray
	bgCapture   = lipgloss.Color("124") // Dark red
	bgCheck     = lipgloss.Color("196") // Bright red
	bgCursor    = lipgloss.Color("62")  // Blue/Purple
	bgSelected  = lipgloss.Color("220") // Yellow
	bgLightSq   = lipgloss.Color("252") // Light board square
	bgDarkSq    = lipgloss.Color("240") // Dark board square
)

type Model struct {
	board         *chess.Board
	cursorSq      uint8
	selectedSq    uint8
	validMoves    chess.MoveList
	statusMessage string
	enterMode     EnterMode
}

func InitialModel() Model {
	b, _ := chess.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return Model{
		board:         b,
		cursorSq:      12,
		selectedSq:    255,
		statusMessage: "White to move",
		validMoves:    b.GenerateLegalMoves(),
		enterMode:     SelectMode,
	}
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursorSq < 56 {
				m.cursorSq += 8
			}
		case "down":
			if m.cursorSq > 7 {
				m.cursorSq -= 8
			}
		case "left":
			if m.cursorSq%8 != 0 {
				m.cursorSq--
			}
		case "right":
			if m.cursorSq%8 != 7 {
				m.cursorSq++
			}
		case "c":
			m.statusMessage = "...thinking"
			move := m.board.SearchBestMove(5)
			if move != 0 {
				m.board.MakeMove(move)
				m.validMoves = m.board.GenerateLegalMoves()
				fromStr, _ := chess.ParseSquareI2S(move.From())
				toStr, _ := chess.ParseSquareI2S(move.To())
				m.statusMessage = fmt.Sprintf("AI played %s %s", fromStr, toStr)
			} else {
				m.statusMessage = "Game Over"
			}
			m.enterMode = SelectMode
			m.selectedSq = 255
		case "enter":
			color := m.board.GetColorType(m.cursorSq)
			//piece:= m.board.GetPieceType(m.cursorSq)
			if m.enterMode == SelectMode {
				if m.board.SideToMove == color {
					m.selectedSq = m.cursorSq
					m.statusMessage = "Piece Selected"
					m.enterMode = MoveMode
				} else if color == 2 {
					m.statusMessage = "Square is empty"
				} else {
					m.statusMessage = "invalid selection"
				}
			} else if m.enterMode == MoveMode {
				if m.cursorSq == m.selectedSq {
					m.selectedSq = 255
					m.statusMessage = "Piece DeSelected"
					m.enterMode = SelectMode
					break
				}
				moveFound := false
				for i := 0; i < m.validMoves.Count; i++ {
					move := m.validMoves.Moves[i]
					if move.From() == m.selectedSq && move.To() == m.cursorSq {
						flags := move.Flags()
						if (flags >= 8 && flags <= 10) || (flags >= 12 && flags <= 14) {
							continue
						}
						m.board.MakeMove(move)
						m.validMoves = m.board.GenerateLegalMoves()
						m.selectedSq = 255
						m.statusMessage = "Moved square"
						m.enterMode = SelectMode
						moveFound = true
						break
					}
				}
				if !moveFound {
					m.selectedSq = 255
					m.statusMessage = "invalid move"
					m.enterMode = SelectMode
				}
			}

		}
	}
	return m, nil
}
func (m Model) ViewChessBoard() string {
	squareStyles := make(map[uint8]lipgloss.Style)
	kingBits := m.board.Colors[m.board.SideToMove] & m.board.Pieces[chess.King]
	if kingBits != 0 {
		kingSq := uint8(bits.TrailingZeros64(kingBits))
		if m.board.IsSquareAttacked(kingSq, m.board.SideToMove^1) {
			squareStyles[kingSq] = lipgloss.NewStyle().Background(bgCheck)
		}
	}
	if m.selectedSq != 255 {
		for i := 0; i < m.validMoves.Count; i++ {
			move := m.validMoves.Moves[i]
			if move.From() == m.selectedSq {
				to := move.To()
				flags := move.Flags()
				isCapture := flags == 4 || flags == 5 || flags >= 12
				if isCapture {
					squareStyles[to] = lipgloss.NewStyle().Foreground(bgCapture)
				} else {
					squareStyles[to] = lipgloss.NewStyle().Foreground(bgValidMove)
				}
			}
		}
		squareStyles[m.selectedSq] = lipgloss.NewStyle().Background(bgSelected).Foreground(lipgloss.Color("0"))
	}
	squareStyles[m.cursorSq] = lipgloss.NewStyle().Background(bgCursor).Foreground(lipgloss.Color("255"))
	chessBoard := make([][]string, 8)
	for rank := 7; rank >= 0; rank-- {
		chessBoard[7-rank] = make([]string, 8)
		for file := 0; file < 8; file++ {
			sq := uint8(rank*8 + file)
			color := m.board.GetColorType(sq)
			ptype := m.board.GetPieceType(sq)
			chessBoard[7-rank][file] = fmt.Sprintf(" %s ", pieceMap[color][ptype])
		}
	}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(chessBoard...).
		StyleFunc(func(row, col int) lipgloss.Style {
			rank := 7 - row
			file := col
			sq := uint8(rank*8 + file)
			baseStyle := lipgloss.NewStyle().Width(3).Height(1).Align(lipgloss.Center)
			if customStyle, exists := squareStyles[sq]; exists {
				return baseStyle.Inherit(customStyle)
			}
			isLightSquare := (rank+file)%2 != 0
			if isLightSquare {
				return baseStyle.Background(bgLightSq).Foreground(lipgloss.Color("16"))
			}
			return baseStyle.Background(bgDarkSq).Foreground(lipgloss.Color("255"))
		})
	ranks := strings.Join([]string{"A", "B", "C", "D", "E", "F", "G", "H  "}, "   ")
	files := strings.Join([]string{" 8", "7", "6", "5", "4", "3", "2", "1 "}, "\n\n ")

	return lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(lipgloss.Center, files, t.Render()),
		ranks,
	)
}
func (m Model) View() string {
	var b strings.Builder
	b.WriteString("GoChess\n")
	b.WriteString(m.ViewChessBoard())
	b.WriteString("\n" + m.statusMessage + "\n")
	b.WriteString("Use arrow keys to move, Space/Enter to select. 'q' to quit.\n")
	return b.String()
}
