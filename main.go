package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/ini.v1"
)

type profileItem struct {
	name string
}

func (p profileItem) FilterValue() string { return p.name }
func (p profileItem) Title() string       { return p.name }
func (p profileItem) Description() string { return "" }

type errMsg error
type reloadMsg struct{}

type model struct {
	config        *ini.File
	allProfiles   []list.Item
	list          list.Model
	input         textinput.Model
	quitting      bool
	filteredItems []list.Item
	lastInput     string

	width  int
	height int

	home       string
	configPath string
}

var warnings []string = []string{}
var warningIcon string = "🚧"

var infos []string = []string{
	"ℹ️ (info) You will see only [profile ...] sections on the list",
	"ℹ️ (info) Profile [default] is updated with chosen profile values",
}

var infoIcon string = "ℹ️"

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(home, ".aws", "config")

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load AWS config file: %v\n", err)
	}

	profiles, warningsFound := loadProfiles(cfg)

	warnings = append(warnings, warningsFound...)

	delegate := minimalDelegate{}
	l := list.New(profiles, delegate, 40, 10)
	l.Title = "Select an AWS profile"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	ti := textinput.New()
	ti.Prompt = "🔍 "
	ti.Placeholder = "Type to filter..."
	ti.Focus()

	m := model{
		config:        cfg,
		allProfiles:   profiles,
		list:          l,
		input:         ti,
		filteredItems: profiles,
		home:          home,
		configPath:    configPath,
	}

	if _, err := bubbletea.NewProgram(m, bubbletea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() bubbletea.Cmd {
	return bubbletea.Batch(
		textinput.Blink,
		watchConfigFile(m.configPath),
	)
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var (
		cmd  bubbletea.Cmd
		cmds []bubbletea.Cmd
	)
	switch msg := msg.(type) {
	case bubbletea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case bubbletea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, bubbletea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem()
			if p, ok := selectedItem.(profileItem); ok {
				err := setDefaultProfile(m.config, p.name, m.configPath)
				if err != nil {
					fmt.Printf("Error setting default profile: %v\n", err)
				}
			}
			m.quitting = true
			return m, bubbletea.Quit
		}

	case errMsg:
		fmt.Printf("Error: %v\n", msg)
		return m, nil

	case reloadMsg:
		// Reload the config file
		err := m.reloadConfig()
		if err != nil {
			fmt.Printf("Error reloading config: %v\n", err)
		}
		// After reloading, start watching again
		cmds = append(cmds, watchConfigFile(m.configPath))
	}

	oldValue := m.input.Value()
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	newValue := m.input.Value()
	if newValue != oldValue {
		m.filterItems(newValue)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, bubbletea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	// If width is zero (not received), assume a default
	if m.width == 0 {
		m.width = 80
	}
	leftColumnWidth := m.width / 2
	rightColumnWidth := m.width - leftColumnWidth - 2

	// Instructions
	instructions := "↑/↓ to navigate • Enter to select • ESC or Ctrl+C to quit"
	instrStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	instructionsView := instrStyle.Render(instructions)

	// Render left column: input + list
	leftView := instructionsView + "\n\n" + m.input.View() + "\n\n" + m.list.View()
	leftStyle := lipgloss.NewStyle().Width(leftColumnWidth)
	leftColumn := leftStyle.Render(leftView)

	// Render right column: details
	details := m.renderDetails()
	rightStyle := lipgloss.NewStyle().Width(rightColumnWidth)
	rightColumn := rightStyle.Render(details)

	// Combine with a gap
	containerStyle := lipgloss.NewStyle().Width(m.width)
	fullView := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")) // Yellowish color
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#81D4FA")) // Blueish color

	// Infos
	if len(infos) != 0 {
		for _, info := range infos {
			infoMessage := infoStyle.Render(info)
			fullView += "\n" + infoMessage
		}
	}

	// Warnings
	if len(warnings) != 0 {
		fullView += "\n"
		for _, warning := range warnings {
			warningMessage := warnStyle.Render(warning)
			fullView += "\n" + warningMessage
		}
	}

	return containerStyle.Render(fullView)
}

func (m *model) filterItems(query string) {
	if query == "" {
		m.filteredItems = m.allProfiles
		m.list.SetItems(m.filteredItems)
		return
	}

	query = strings.ToLower(query)
	var filtered []list.Item
	for _, it := range m.allProfiles {
		p := it.(profileItem)
		if strings.Contains(strings.ToLower(p.name), query) {
			filtered = append(filtered, it)
		}
	}

	m.filteredItems = filtered
	m.list.SetItems(m.filteredItems)

	// Ensure valid selection after filtering
	if len(m.filteredItems) == 0 {
		m.list.Select(-1)
	} else {
		selected := m.list.Index()
		if selected < 0 || selected >= len(m.filteredItems) {
			m.list.Select(0)
		}
	}
}

func (m *model) renderDetails() string {
	si := m.list.SelectedItem()
	if si == nil {
		return "No profile selected."
	}
	p, ok := si.(profileItem)
	if !ok {
		return "No profile selected."
	}

	section := m.config.Section("profile " + p.name)
	if section == nil {
		return "No details available."
	}

	keys := section.Keys()
	if len(keys) == 0 {
		return "No keys in this profile."
	}

	var sb strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	styleBlue := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	styleGreen := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	instrStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	instructionsView := instrStyle.Render("Selected profile")
	sb.WriteString(instructionsView + "\n\n")
	title := titleStyle.Render("[" + section.Name() + "]")
	sb.WriteString(fmt.Sprintf("%s\n\n", title))
	for _, key := range keys {
		keyBlue := styleBlue.Render(key.Name())
		valueGreen := styleGreen.Render(key.Value())
		sb.WriteString(fmt.Sprintf("%s=%s\n", keyBlue, valueGreen))
	}
	return sb.String()
}

func (m *model) reloadConfig() error {
	cfg, err := ini.Load(m.configPath)
	if err != nil {
		return err
	}
	m.config = cfg

	profiles, warningsFound := loadProfiles(cfg)
	warnings = warningsFound
	m.allProfiles = profiles
	m.filteredItems = profiles
	m.list.SetItems(profiles)
	if len(profiles) > 0 {
		m.list.Select(0)
	}

	return nil
}

func loadProfiles(cfg *ini.File) ([]list.Item, []string) {
	var profiles []list.Item
	var w []string
	for _, section := range cfg.Sections() {
		name := section.Name()
		if name == "DEFAULT" || name == "default" {
			continue
		} else if len(name) > 8 && name[:8] == "profile " {
			name = name[8:]
			profiles = append(profiles, profileItem{name: name})
		} else if !(len(name) > 12 && name[:12] == "sso-session ") {
			warningMessage := fmt.Sprintf("%s (warn) Found sections with missing 'profile' or 'sso-session' prefix!\n\nRead more:\n%s", warningIcon, "https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html")
			w = append(w, warningMessage)
		}
	}
	return profiles, w
}

func setDefaultProfile(cfg *ini.File, chosen, configPath string) error {
	var chosenSection, defaultSection *ini.Section
	chosenSection = cfg.Section("profile " + chosen)

	if chosenSection == nil {
		return fmt.Errorf("chosen profile %q not found in config", chosen)
	}

	if !cfg.HasSection("default") {
		newDefaultSection, err := cfg.NewSection("default")
		if err != nil {
			return err
		}
		if newDefaultSection == nil {
			return fmt.Errorf("failed to create default section")
		}
		defaultSection = newDefaultSection
	} else {
		defaultSection = cfg.Section("default")
	}

	for _, key := range defaultSection.Keys() {
		defaultSection.DeleteKey(key.Name())
	}

	for _, key := range chosenSection.Keys() {
		fmt.Println("Writing key to ", defaultSection.Name(), key.Name(), key.Value())
		_, err := defaultSection.NewKey(key.Name(), key.Value())
		if err != nil {
			return err
		}
	}
	err := cfg.SaveTo(configPath)
	if err != nil {
		fmt.Println("Error saving config:", err)
	}
	time.Sleep(20 * time.Second)
	return cfg.SaveTo(configPath)
}

// A minimal delegate to reduce spacing
type minimalDelegate struct{}

func (d minimalDelegate) Height() int {
	return 1
}

func (d minimalDelegate) Spacing() int {
	return 0
}

func (d minimalDelegate) Update(msg bubbletea.Msg, m *list.Model) bubbletea.Cmd {
	return nil
}

func (d minimalDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	pi, ok := item.(profileItem)
	if !ok {
		return
	}

	indicator := "  "
	styledName := pi.name
	indexStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	if index == m.Index() {
		indicator = "> "
		styledName = indexStyle.Render(pi.name)
	}

	fmt.Fprintf(w, "%s%s", indicator, styledName)
}

// watchConfigFile sets up a fsnotify watcher that sends reloadMsg when the file changes
func watchConfigFile(path string) bubbletea.Cmd {
	return func() bubbletea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println("Error creating watcher:", err)
			return nil
		}
		defer watcher.Close()

		err = watcher.Add(path)
		if err != nil {
			log.Println("Error adding watch:", err)
			return nil
		}

		// Wait for events
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					return reloadMsg{}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
				log.Println("Watcher error:", err)
			}
		}
	}
}
