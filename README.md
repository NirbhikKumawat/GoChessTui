# GoChess TUI ♟️

A sleek, terminal-based chess interface built in Go using the `bubbletea` and `cobra` framework.

This TUI serves as the visual frontend and debugging client for [Knightmare](htttps://github.com/NirbhikKumawat/GoChess), a custom-built, blazingly fast 64-bit Bitboard chess engine.

## 🔗 The Engine Backend
This client is powered entirely by my own chess engine **[Knightmare](htttps://github.com/NirbhikKumawat/GoChess)** built from scratch. 

## 🎨 Features
* **Elm Architecture:** Built on Charmbracelet's `bubbletea` for a snappy, state-driven terminal experience.
* **Rich Highlighting:** Uses `lipgloss` to render a fully checkered board with dynamic highlighting for the cursor, selected pieces, valid moves, captures, and kings in check.
* **Keyboard Driven:** Navigate the board entirely using the arrow keys and the Space/Enter bar.
* **Strict Legality:** Enforces 100% strict chess rules (en passant, castling rights, absolute pins) by interfacing directly with the backend bitboard move generator.

## 🚀 Running the TUI
Ensure you have Go installed, then clone the repository and run:
`go run .`

**Controls:**
* `Arrow Keys`: Move cursor
* `Enter`: Select piece / Confirm move
* `Q` / `Ctrl+C`: Quit