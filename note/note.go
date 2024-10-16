package note

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vineshtv/yawn/config"
)

type State int

const (
	editingName State = iota
	editingBody
	onSave
	committingSave
	saveSuccess
)

type noteModel struct {
	// Name of the note.
	Name textinput.Model
	// Note body
	Body textarea.Model

	// Something to do with help. TODO: implement.
	help help.Model

	savingSpinner spinner.Model
	// State indicates the current state of this new Note.
	state State

	NoteFileName string
}

type NoteSaveModel struct {
	Name string
	Body string
}

const (
	accentColor    = lipgloss.Color("99")
	yellowColor    = lipgloss.Color("#ECFD66")
	whiteColor     = lipgloss.Color("255")
	grayColor      = lipgloss.Color("241")
	darkGrayColor  = lipgloss.Color("236")
	lightGrayColor = lipgloss.Color("247")
)

var (
	labelStyle              = lipgloss.NewStyle().Foreground(grayColor)
	textStyle               = lipgloss.NewStyle().Foreground(lightGrayColor)
	cursorStyle             = lipgloss.NewStyle().Foreground(whiteColor)
	placeholderStyle        = lipgloss.NewStyle().Foreground(darkGrayColor)
	activeTextStyle         = lipgloss.NewStyle().Foreground(whiteColor)
	activeLabelStyle        = lipgloss.NewStyle().Foreground(accentColor)
	saveButtonActiveStyle   = lipgloss.NewStyle().Background(accentColor).Foreground(yellowColor).Padding(0, 2)
	saveButtonInactiveStyle = lipgloss.NewStyle().Background(darkGrayColor).Foreground(lightGrayColor).Padding(0, 2)
	saveButtonStyle         = lipgloss.NewStyle().Background(darkGrayColor).Foreground(grayColor).Padding(0, 2)
	helpStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	yawnstyle               = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"})
)

func (m *noteModel) focusActiveInput() {
	// Reset the text input styles
	m.Name.PromptStyle = labelStyle
	m.Name.TextStyle = textStyle

	switch m.state {
	case editingName:
		m.Name.PromptStyle = activeLabelStyle
		m.Name.TextStyle = activeTextStyle
		m.Name.Focus()
		m.Name.CursorEnd()
	case editingBody:
		m.Body.Focus()
		m.Body.CursorEnd()
	}
}

func AddNewNoteModel() noteModel {
	name := textinput.New()
	name.Prompt = "Note name: "
	name.Placeholder = "My awesome note"
	name.PromptStyle = labelStyle
	name.TextStyle = textStyle
	name.Cursor.Style = cursorStyle
	name.PlaceholderStyle = placeholderStyle

	body := textarea.New()
	body.Placeholder = "# Note"
	body.ShowLineNumbers = false
	body.FocusedStyle.CursorLine = activeTextStyle
	body.FocusedStyle.Prompt = activeLabelStyle
	body.FocusedStyle.Text = activeTextStyle
	body.FocusedStyle.Placeholder = placeholderStyle
	body.BlurredStyle.CursorLine = textStyle
	body.BlurredStyle.Prompt = labelStyle
	body.BlurredStyle.Text = textStyle
	body.BlurredStyle.Placeholder = placeholderStyle
	body.Cursor.Style = cursorStyle
	body.CharLimit = 4000
	body.SetWidth(80)
	body.SetHeight(15)

	savingSpinner := spinner.New()
	savingSpinner.Style = activeLabelStyle
	savingSpinner.Spinner = spinner.Dot

	m := noteModel{
		Name:          name,
		Body:          body,
		help:          help.New(),
		savingSpinner: savingSpinner,
		state:         editingName,
		NoteFileName:  sanitizeFileName(name.Value()),
	}

	m.focusActiveInput()

	return m
}

func sanitizeFileName(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	sanitized := re.ReplaceAllString(input, "_")
	sanitized = strings.Trim(sanitized, "._-")

	if len(sanitized) == 0 {
		sanitized = "default_filename"
	}

	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}

type savingNoteSuccss struct{}

func (m *noteModel) saveNote() tea.Cmd {
	return func() tea.Msg {
		// So we need to save the note here
		// get the notes location from config
		dirname := config.Config.General.NoteLocation
		if _, err := os.Stat(dirname); os.IsNotExist(err) {
			// notes location does not exist. Create it
			err := os.MkdirAll(dirname, 0755)
			if err != nil {
				fmt.Println("Error creating notes directory - ", err)
				os.Exit(1)
			}
		}

		// Save the contents of the note
		noteName := m.Name.Value()
		noteBody := m.Body.Value()

		noteFileName := sanitizeFileName(noteName)
		noteFullPath := fmt.Sprintf("%s/%s", dirname, noteFileName)

		// if the file name for the note has been modified, then delete the old file first
		if m.NoteFileName != noteFileName {
			oldFilePath := fmt.Sprintf("%s/%s", dirname, m.NoteFileName)

			if _, err := os.Stat(oldFilePath); err == nil {
				err = os.Remove(oldFilePath)
				if err != nil {
					fmt.Println("Error deleting old file", err)
					os.Exit(1)

				}
			}
		}
		m.NoteFileName = noteFileName

		fd, err := os.Create(noteFullPath)
		if err != nil {
			fmt.Println("Error Writing Note: ", err)
			os.Exit(1)
		}
		defer fd.Close()

		m1 := NoteSaveModel{
			Name: noteName,
			Body: noteBody,
		}
		encoder := json.NewEncoder(fd)
		err = encoder.Encode(m1)
		if err != nil {
			fmt.Println("Error Encoding file: ", err)
		}
		return savingNoteSuccss{}
	}
}

// Init BubbleTea init method.
func (m noteModel) Init() tea.Cmd {
	return nil
}

func (m *noteModel) blurInputs() {
	m.Name.Blur()
	m.Body.Blur()
}

// Update BubbleTea update method.
func (m noteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case savingNoteSuccss:
		m.state = saveSuccess
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.blurInputs()
			switch m.state {
			case editingName:
				m.state = editingBody
				m.Body.Focus()
			case editingBody:
				m.state = onSave
			case onSave:
				m.state = editingName
				m.Name.Focus()
			}
			m.focusActiveInput()
		case "enter":
			if m.state == onSave {
				m.state = committingSave
				return m, tea.Batch(
					m.savingSpinner.Tick,
					m.saveNote(),
				)
			}
		}
	}

	// Update the actual fields with the new value
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.Name, cmd = m.Name.Update(msg)
	m.Body, cmd = m.Body.Update(msg)
	// _ = cmd

	switch m.state {
	case committingSave:
		m.savingSpinner, cmd = m.savingSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m noteModel) View() string {
	switch m.state {
	case committingSave:
		return m.savingSpinner.View() + "Saving Note"
	case saveSuccess:
		// return "Note saved... yawn ᶻ 𝗓 𐰁\n"
		return yawnstyle.Render("Note Saved. yawn ᶻ 𝗓 𐰁\n")
	}
	var s strings.Builder

	s.WriteString(m.Name.View())
	s.WriteString("\n\n")
	s.WriteString(m.Body.View())
	s.WriteString("\n\n")

	switch m.state {
	case onSave:
		s.WriteString(saveButtonActiveStyle.Render("Save Note"))
	case committingSave:
		s.WriteString(saveButtonActiveStyle.Render("Save Note"))
	default:
		s.WriteString(saveButtonInactiveStyle.Render("Save Note"))
	}
	s.WriteString(helpStyle.Render("\n\ntab: focus next • ctrl+c: exit\n"))

	return lipgloss.NewStyle().Padding(1).Render(s.String())
}
