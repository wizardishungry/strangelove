package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"jonwillia.ms/strangelove/citi"
	"jonwillia.ms/strangelove/clock"
)

const lat, lon = 40.688265, -73.9184594

// https://github.com/charmbracelet/bubbletea/pull/181/files
var (
	defaultWidth = 20

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "└",
		BottomRight: "┘",
	}
	highlight       = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	highlightActive = lipgloss.AdaptiveColor{Light: "#A76BFD", Dark: "#9D76F4"}
	tabBorder       = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "└",
		BottomRight: "┘",
	}
	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1).Height(2)
	activeTab = tab.Copy().Border(activeTabBorder, true).BorderForeground(highlightActive)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type Model struct {
	Tabs       []string
	TabContent []string
	Reading    clock.Reading

	activatedTab int
	citi         <-chan []string
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.bikeShare,
		m.doTick(),
	)
}

func (m *Model) doTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return clock.Coords{Lat: lat, Lon: lon}.Time(t)
	})
}

type bikeMessage []string

func (m *Model) bikeShare() tea.Msg {
	return (bikeMessage)(<-m.citi)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//case tea.WindowSizeMsg:
	case bikeMessage:
		m.Tabs = msg
		return m, m.bikeShare
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		// Cycle through tabs to the right
		case "right":
			if m.activatedTab < len(m.Tabs)-1 {
				m.activatedTab++
			}
		// Cycle through tabs to the left
		case "left":
			if m.activatedTab > 0 {
				m.activatedTab--
			}
		}
	case clock.Reading:
		m.Reading = msg
		return m, m.doTick()
	}

	return m, nil
}

func (m *Model) View() string {
	doc := strings.Builder{}

	// Clock
	var cellClock string
	{
		cellClock = activeTab.Render(m.Reading.Render())
	}

	// Bikes
	var cellBike string
	{
		var renderedTabs []string

		var maxWidth int
		for _, t := range m.Tabs {
			maxWidth = max(maxWidth, lipgloss.Width(t))
		}
		maxWidth += tab.GetHorizontalPadding()
		activeTab := activeTab.Width(maxWidth)
		tab := tab.Width(maxWidth)

		// Activate the correct tab
		for i, t := range m.Tabs {
			if i == m.activatedTab {
				renderedTabs = append(renderedTabs, activeTab.Render(t))
			} else {
				renderedTabs = append(renderedTabs, tab.Render(t))
			}
		}

		cellBike = lipgloss.JoinVertical(
			lipgloss.Top,
			renderedTabs...,
		)
		// gap := tabGap.Render(strings.Repeat(" ", max(0, defaultWidth-lipgloss.Width(row)-2)))
		// row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	}
	d := lipgloss.JoinHorizontal(lipgloss.Left, cellClock, cellBike)
	doc.WriteString(d + "\n\n")

	doc.WriteString("whatever!!!" + fmt.Sprintf("%d", len(m.Tabs)))

	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func NewModel() *Model {
	tabs := []string{
		"Lip Gloss",
		"Blush",
		"Eye Shadow\nEast",
		"Mascara",
		"Foundation",
	}

	tabContent := []string{
		"tab1",
		"tab2",
		"tab3",
		"tab4",
		"tab5",
	}

	return &Model{
		citi:       citi.Citi(lat, lon),
		Tabs:       tabs,
		TabContent: tabContent,
	}
}
