package cmd

import (
	// "bytes"
	"encoding/json"
	"fmt"

	// "io"
	"os"
	"os/exec"

	// "strings"

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

func runsearch() (string, error) {
	cmd := exec.Command("fzf")
	cmd.Stdin = os.Stdin
	cmd.Dir = "/Users/vineshtv/pancake/YAWN/yawn/Notes"

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("fzf failed - ", err)
		return "", err
	}

	return string(output), err
}

func findNote(cmd *cobra.Command, args []string) {
	foundNote := note.AddNewNoteModel()
	// foundNote.Name.SetValue("This note")
	// foundNote.Body.SetValue("This is the body of the note.")

	noteName, err := runsearch()
	if err != nil {
		// TODO: Print an error.
		return
	}

	dirname := "/Users/vineshtv/pancake/YAWN/yawn/Notes"
	if _, err = os.Stat(dirname); os.IsNotExist(err) {
		fmt.Println("error - ", err)
	} else {
		fmt.Println("Directory exists")
	}
	notefullPath := fmt.Sprintf("%s/%s", dirname, noteName)
	notefullPath = "/Users/vineshtv/pancake/YAWN/yawn/Notes/vineshnote"

	if _, err = os.Stat(notefullPath); err != nil {
		fmt.Println(err)
	}
	fmt.Println(os.Getwd())
	fmt.Println("Opening - ", notefullPath)
	file, err := os.Open(notefullPath)
	if err != nil {
		// TODO: print errork
		fmt.Println("Error opening note - ", err)
		return
	}
	defer file.Close()

	m1 := note.NoteSaveModel{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&m1)
	if err != nil {
		fmt.Println("Error decoding note - ", err)
		return
	}
	foundNote.Name.SetValue(m1.Name)
	foundNote.Body.SetValue(m1.Body)

	fmt.Println(file)
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
