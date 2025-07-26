# t-switch
A simple terminal-based theme switcher that edits your configuration files to change between different color themes. It uses regex to find specific parts of configuration files to edit.

## Installation
1. Clone the repository:
```bash
git clone https://github.com/pj-hugo/t-switch.git
cd t-switch
```

2. Build and install:
```bash
make install
```

3. Create the configuration directory:
```bash
mkdir -p ~/.config/t-switch
```

4. Create the necessary configuration files:
```bash
touch ~/.config/t-switch/configs.yaml
touch ~/.config/t-switch/themes.yaml
```

## Configuration

t-switch uses two configuration files located in `~/.config/t-switch/`:

### 1. themes.yaml

Define your color themes here. Each theme is a collection of key-value pairs representing colors or other theme properties.

```yaml
# Example themes.yaml

kanagawa-wave:
  nvim: "kanagawa-wave"
  wezterm: 'Kanagawa (Gogh)'
  starship_dir: '#7E9CD8'
  starship_git: '#76946A'
  tmux-status-style: 'fg=#DCD7BA,bg=#16161D'

rosé-pine-dawn:
  nvim: "rose-pine-dawn"
  wezterm: 'rose-pine-dawn'
  starship_dir: '#907aa9'
  starship_git: '#56949f'
  tmux-status-style: 'fg=#575279,bg=#f2e9e1'
```

### 2. configs.yaml

Define which applications to theme and how to apply the themes.

```yaml
# Example configs.yaml

nvim:
  path: "~/.config/nvim/lua/active_colorscheme.lua"
  replacements:
    - key: nvim
      regex: 'return ".*"'
      replace: 'return "{}"'

wezterm:
  path: "~/.config/wezterm/wezterm.lua"
  replacements:
    - key: wezterm
      regex: "config.color_scheme = ['\"].*['\"]"
      replace: "config.color_scheme = '{}'"

starship:
  path: "~/.config/starship/starship.toml"
  replacements:
    - key: starship_dir
      regex: "dir = '#.*'"
      replace: "dir = '{}'"
    - key: starship_git
      regex: "git = '.*'"
      replace: "git = '{}'"

tmux:
  path: "~/.tmux.conf"
  cmd: "tmux source-file ~/.tmux.conf"
  replacements:
    - key: tmux-status-style
      regex: ".* # t-switch:status-style #"
      replace: "set -g status-style '{}' # t-switch:status-style #"
```

## Usage

Simply run the command:

```bash
t-switch
```

Use the arrow keys or `j`/`k` to navigate, press `Enter` to select a theme, or `q` to quit.

## Configuration Details

### Replacement Rules

Each replacement rule consists of:

- **key**: The theme property to use (must exist in themes.yaml)
- **regex**: A regular expression to find the text to replace
- **replace**: The replacement string. Use `{}` as a placeholder for the theme value

### Commands

You can specify a `cmd` field for any application to run a command after the theme is applied. This is useful for:
- Reloading configuration files
- Restarting services
- Sending signals to running applications

## Uninstallation

```bash
make uninstall
```

This removes the binary but keeps your configuration files. To remove everything also run:

```bash
rm -rf ~/.config/t-switch
```

## Acknowledgements
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - Used for making the TUI.
- [cultab/themr](https://github.com/cultab/themr) - What inspired me to create this application.
