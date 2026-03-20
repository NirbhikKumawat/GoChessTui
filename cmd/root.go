package cmd

import (
	"GoChessTui/internal/ui"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gochess",
	Short: "A TUI chess",
	Long:  "A TUI chess",
	Run:   RootRun,
}

func RootRun(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
