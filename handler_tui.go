package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mrjacz/gator/internal/database"
)

type tuiModel struct {
	posts    []database.Post
	cursor   int
	selected map[int]struct{}
	viewing  bool
	err      error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("235"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)
)

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.posts)-1 {
				m.cursor++
			}

		case "enter":
			if m.viewing {
				m.viewing = false
			} else {
				m.viewing = true
			}

		case "o":
			// Open in browser
			if len(m.posts) > 0 && m.cursor < len(m.posts) {
				openBrowser(m.posts[m.cursor].Url)
			}

		case "esc":
			if m.viewing {
				m.viewing = false
			}
		}
	}

	return m, nil
}

func (m tuiModel) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err))
	}

	if len(m.posts) == 0 {
		return "No posts found.\n\nPress q to quit."
	}

	if m.viewing {
		return m.renderDetailView()
	}

	return m.renderListView()
}

func (m tuiModel) renderListView() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("RSS Posts"))
	s.WriteString("\n\n")

	for i, post := range m.posts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %d. %s", cursor, i+1, post.Title)

		if m.cursor == i {
			line = selectedStyle.Render(line)
		}

		s.WriteString(line)
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/k up • ↓/j down • enter view • o open in browser • q quit"))
	s.WriteString("\n")

	return s.String()
}

func (m tuiModel) renderDetailView() string {
	post := m.posts[m.cursor]

	var content strings.Builder

	content.WriteString(titleStyle.Render(post.Title))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("Published: %s\n", post.PublishedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("URL: %s\n\n", post.Url))

	description := post.Description
	if len(description) > 500 {
		description = description[:500] + "..."
	}
	content.WriteString(description)
	content.WriteString("\n\n")

	content.WriteString(helpStyle.Render("enter/esc back to list • o open in browser • q quit"))

	return detailStyle.Render(content.String())
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func handlerTUI(s *state, cmd command, user database.User) error {
	limit := 20

	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s (no arguments)", cmd.Name)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	initialModel := tuiModel{
		posts:    posts,
		cursor:   0,
		selected: make(map[int]struct{}),
		viewing:  false,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
