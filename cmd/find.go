package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vineshtv/yawn/config"
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

func runSearch(dirName string) (string, error) {
	// TODO: Move this to a utils folder and support other OSes
	cmd := exec.Command("fzf", "--ansi", "--layout=reverse", "--border", "--height=90%", "--cycle")
	// cmd := exec.Command("fzf --ansi --layout=reverse --border --height=90% --cycle")
	cmd.Dir = dirName
	cmd.Stdin = os.Stdin

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("fzf failed -", err)
		return "", err
	}

	// There probably is a better way to do this.
	filename := fmt.Sprintf("%s/%s", dirName, strings.TrimSuffix(string(output), "\n"))
	return filename, err
}

// findNote1 Finds an existing note based on just the note name.
func findNote(cmd *cobra.Command, args []string) {
	dirname := config.Config.General.NoteLocation
	noteName, err := runSearch(dirname)
	if err != nil {
		fmt.Println("Error searching note: ", err)
		os.Exit(1)
	}

	fd, err := os.Open(noteName)
	if err != nil {
		fmt.Println("Error reading note:", err)
		os.Exit(1)
	}
	defer fd.Close()

	m1 := note.NoteSaveModel{}
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&m1)
	if err != nil {
		fmt.Println("Error decoding note: ", err)
		os.Exit(1)
	}

	// Create a new Notemodel to viaualize
	foundNote := note.AddNewNoteModel()
	foundNote.Name.SetValue(m1.Name)
	foundNote.Body.SetValue(m1.Body)
	foundNote.NoteFileName = filepath.Base(noteName)

	// run the tea program
	p := tea.NewProgram(foundNote)
	m, err := p.Run()
	if err != nil {
		fmt.Println("Yawn broke: ", err)
		os.Exit(1)
	}
	_ = m
}

func init() {
	rootCmd.AddCommand(findCmd)
}
