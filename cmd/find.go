package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vineshtv/yawn/note"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a note.",
	Long: `Find a note
	Example usage:
	yawn find`,
	Run: findNote,
}

func findNote(cmd *cobra.Command, args []string) {
	foundNote := note.AddNewNoteModel()
	foundNote.Name.SetValue("This note")
	foundNote.Body.SetValue("This is the body of the note.")

	p := tea.NewProgram(foundNote)
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Yawn broke -- %v", err)
	}

	_ = m
}

func init() {
	rootCmd.AddCommand(findCmd)
}
