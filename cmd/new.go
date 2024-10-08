package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vineshtv/yawn/note"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note.",
	Long: `Create a new note
	Example usage:
	yawn new`,
	Run: addNote,
}

func addNote(cmd *cobra.Command, args []string) {
	fmt.Println("Adding a new note")
	p := tea.NewProgram(note.AddNewNoteModel())

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Yawn broke -- %v", err)
	}

	_ = m
}

func init() {
	rootCmd.AddCommand(newCmd)
}
