package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

var (
	allThemes   map[string]map[string]string
	configRules map[string]configRule
)

type configRule struct {
	Path         string            `yaml:"path"`
	Replacements []replacementRule `yaml:"replacements"`
	Cmd          string            `yaml:"cmd,omitempty"`
}

type replacementRule struct {
	Key     string `yaml:"key"`
	Regex   string `yaml:"regex"`
	Replace string `yaml:"replace"`
}

type model struct {
	cursor      int
	themes      []string
	chosenTheme *string
}

func initialModel() model {
	var themeNames []string
	for name := range allThemes {
		themeNames = append(themeNames, name)
	}
	sort.Strings(themeNames)

	return model{themes: themeNames}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("t-switch")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}
		case "enter":
			m.chosenTheme = &m.themes[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString("Choose a theme:\n\n")

	for i, theme := range m.themes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, theme))
	}

	s.WriteString("\nPress 'enter' to select, 'q' to quit.\n")
	return s.String()
}

func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return strings.Replace(path, "~", usr.HomeDir, 1), nil
}

func loadThemes() error {
	themesPath, err := expandPath("~/.config/t-switch/themes.yaml")
	if err != nil {
		return fmt.Errorf("could not determine themes path: %w", err)
	}

	themeFile, err := os.ReadFile(themesPath)
	if err != nil {
		return fmt.Errorf("error reading themes.yaml: %w", err)
	}

	if err := yaml.Unmarshal(themeFile, &allThemes); err != nil {
		return fmt.Errorf("error parsing themes.yaml: %w", err)
	}

	return nil
}

func loadConfig() error {
	configsPath, err := expandPath("~/.config/t-switch/configs.yaml")
	if err != nil {
		return fmt.Errorf("could not determine configs path: %w", err)
	}

	configFile, err := os.ReadFile(configsPath)
	if err != nil {
		return fmt.Errorf("error reading configs.yaml: %w", err)
	}

	if err := yaml.Unmarshal(configFile, &configRules); err != nil {
		return fmt.Errorf("error parsing configs.yaml: %w", err)
	}

	return nil
}

func applyTheme(themeName string) error {
	selectedThemeValues := allThemes[themeName]

	for appName, rule := range configRules {
		if err := applyRuleToApp(appName, rule, selectedThemeValues, themeName); err != nil {
			return err
		}
	}
	return nil
}

func applyRuleToApp(appName string, rule configRule, themeValues map[string]string, themeName string) error {
	filePath, err := expandPath(rule.Path)
	if err != nil {
		return fmt.Errorf("could not expand path for %s: %w", appName, err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file %s for app '%s': %w", filePath, appName, err)
	}

	modifiedContent := string(content)

	for _, rep := range rule.Replacements {
		value, ok := themeValues[rep.Key]
		if !ok {
			log.Printf("Warning: key '%s' not found in theme '%s'", rep.Key, themeName)
			continue
		}

		re, err := regexp.Compile(rep.Regex)
		if err != nil {
			return fmt.Errorf("invalid regex for key '%s' in app '%s': %w", rep.Key, appName, err)
		}

		replacementString := strings.Replace(rep.Replace, "{}", value, 1)
		modifiedContent = re.ReplaceAllString(modifiedContent, replacementString)
	}

	if err := os.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("could not write to file %s for app '%s': %w", filePath, appName, err)
	}

	if rule.Cmd != "" {
		runCommand(rule.Cmd, appName)
	}

	return nil
}

func runCommand(cmdStr, appName string) {
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Env = os.Environ()

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Warning: command '%s' for app '%s' failed: %v", cmdStr, appName, err)
		if stderr.Len() > 0 {
			log.Printf("Stderr: %s", stderr.String())
		}
	} else if stdout.Len() > 0 || stderr.Len() > 0 {
		if stdout.Len() > 0 {
			log.Printf("Stdout: %s", stdout.String())
		}
		if stderr.Len() > 0 {
			log.Printf("Stderr: %s", stderr.String())
		}
	}
}

func main() {
	if err := loadThemes(); err != nil {
		log.Fatalf("Failed to load themes: %v", err)
	}

	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok && m.chosenTheme != nil {
		if err := applyTheme(*m.chosenTheme); err != nil {
			log.Fatalf("Failed to apply theme '%s': %v", *m.chosenTheme, err)
		}
	} else {
		fmt.Println("No theme selected. Exiting.")
	}
}
