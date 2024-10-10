package note

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	savingSpinner := spinner.New()
	savingSpinner.Style = activeLabelStyle
	savingSpinner.Spinner = spinner.Dot

	m := noteModel{
		Name:          name,
		Body:          body,
		help:          help.New(),
		savingSpinner: savingSpinner,
		state:         editingName,
	}

	m.focusActiveInput()

	return m
}

type savingNoteSuccss struct{}

func (m *noteModel) saveNote() tea.Cmd {
	return func() tea.Msg {
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

	m.focusActiveInput()

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
		return "\n " + m.savingSpinner.View() + "Saving Note"
	case saveSuccess:
		return "\n Note saved.\n"
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
