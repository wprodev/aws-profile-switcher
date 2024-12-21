package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	bubbletea "github.com/charmbracelet/bubbletea"
	"gopkg.in/ini.v1"
)

type profileItem struct {
	name string
}

func (p profileItem) FilterValue() string { return p.name }
func (p profileItem) Title() string       { return p.name }
func (p profileItem) Description() string { return "" }

type minimalDelegate struct{}

func (d minimalDelegate) Height() int {
	// Only one line per item
	return 1
}

func (d minimalDelegate) Spacing() int {
	// No extra blank lines between items
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

	// If it's the selected item, add a "> " indicator, else just spaces
	indicator := "  "
	if index == m.Index() {
		indicator = "> "
	}

	fmt.Fprintf(w, "%s%s\n", indicator, pi.name)
}

type errMsg error

type model struct {
	config        *ini.File
	allProfiles   []list.Item
	list          list.Model
	input         textinput.Model
	quitting      bool
	filteredItems []list.Item
	lastInput     string
}

var configPath string = ""

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configPath = filepath.Join(home, ".aws", "config")

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load AWS config file: %v\n", err)
	}

	// Collect profiles
	var profiles []list.Item
	for _, section := range cfg.Sections() {
		name := section.Name()
		if name == "DEFAULT" || name == "default" {
			// don't list default
			continue
		} else if len(name) > 8 && name[:8] == "profile " {
			name = name[8:]
		}
		if name != "" {
			profiles = append(profiles, profileItem{name: name})
		}
	}

	delegate := minimalDelegate{}
	l := list.New(profiles, delegate, 40, 10)
	l.Title = "Select an AWS profile"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()

	m := model{
		config:        cfg,
		allProfiles:   profiles,
		list:          l,
		input:         ti,
		filteredItems: profiles,
	}

	if _, err := bubbletea.NewProgram(m).Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() bubbletea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var (
		cmd  bubbletea.Cmd
		cmds []bubbletea.Cmd
	)

	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, bubbletea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem()
			if p, ok := selectedItem.(profileItem); ok {
				err := setDefaultProfile(m.config, p.name)
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
	return fmt.Sprintf("%s\n\n%s\n", m.input.View(), m.list.View())
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

func setDefaultProfile(cfg *ini.File, chosen string) error {
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
		_, err := defaultSection.NewKey(key.Name(), key.Value())
		if err != nil {
			return err
		}
	}

	return cfg.SaveTo(configPath)
}
