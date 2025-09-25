# Color CLI

A fast and flexible terminal color management tool built with Go and Cobra.

## Features

- **Automatic Directory Colors**: Each directory gets a unique, consistent color based on its path hash
- **Claude Code Session Themes**: Blue/purple themes optimized for Claude Code sessions  
- **Color Cycling**: Generate variations based on current terminal color
- **Command Wrapping**: Wrap commands with automatic color management
- **Redis Persistence**: Colors persist across terminal sessions and reboots
- **Cross-Platform**: Works on any system with Redis support
- **Fast & Reliable**: Built with Go for speed and cross-platform compatibility

## Installation

### Quick Install
```bash
make install
```

### Manual Build
```bash
go build -o color .
mv color ~/.local/bin/
```

### Development Install
```bash
make dev-install  # Creates symlink for development
```

## Usage

### Basic Commands

```bash
# Cycle through color variations (default behavior)
color

# Apply Claude Code session theme
color claude

# Apply directory-based theme
color directory
color directory /path/to/directory

# Cycle through specific color modes
color cycle hue_shift
color cycle brightness
color cycle saturation
color cycle complement

# Reset to default dark theme
color reset

# Wrap a command with automatic color management
color wrap claude --help
color wrap claude code --session mysession

# Check persistence status and color history
color status

# Clear all stored colors
color clear
```

### Integration with Shell

Add to your `.zshrc` or `.bashrc`:

```bash
# Automatic directory color changes
if [[ -n "$ITERM_SESSION_ID" ]]; then
    chpwd() {
        if [[ -z "$CLAUDE_SESSION_ACTIVE" ]]; then
            color directory "$PWD"
        fi
    }
    
    # Claude wrapper function
    claude() {
        export CLAUDE_SESSION_ACTIVE=true
        color claude
        command claude "$@"
        local exit_code=$?
        unset CLAUDE_SESSION_ACTIVE
        color directory "$PWD"
        return $exit_code
    }
fi
```

### Available Color Modes

- **`hue_shift`**: Shift the hue while keeping saturation/value (default)
- **`brightness`**: Adjust brightness/value
- **`saturation`**: Adjust color saturation  
- **`complement`**: Use complementary color
- **`random`**: Random mode selection

## How It Works

### Directory Colors
Each directory path is hashed using MD5, and the hash is used to generate consistent HSV color values:
- **Hue**: Based on first 8 hex characters of hash
- **Saturation**: 0.4-0.7 range based on next 2 hex characters  
- **Value**: 0.15-0.3 range based on next 2 hex characters

### Claude Themes
Blue/purple color palette optimized for terminal readability:
- **Hue**: 0.6, 0.75, or 0.85 (blue to purple range)
- **Saturation**: 0.3-0.7 for good contrast
- **Value**: 0.15-0.25 (kept dark for terminal use)

### iTerm2 Integration
Uses AppleScript to communicate with iTerm2:
- Gets current background color via AppleScript
- Sets new background color via AppleScript
- Converts between RGB (0-255) and iTerm2 values (0-65535)

## Development

### Build System
```bash
make help          # Show available targets
make build         # Build binary
make test          # Run tests
make fmt           # Format code
make clean         # Clean build artifacts
```

### Project Structure
```
├── cmd/           # Cobra command definitions
│   ├── root.go    # Root command and CLI setup
│   ├── claude.go  # Claude theme command
│   ├── directory.go # Directory theme command
│   ├── cycle.go   # Color cycling command
│   ├── reset.go   # Reset command
│   └── wrapper.go # Command wrapper
├── internal/      # Internal packages
│   └── color.go   # Color management logic
├── main.go        # Application entry point
├── Makefile       # Build system
└── README.md      # This file
```

## Requirements

- **macOS** with iTerm2 (currently iTerm2 specific)
- **Go 1.21+** for building from source

## Migrating from Python Version

The Go CLI is a drop-in replacement for the Python version:

```bash
# Old Python usage
python3 ~/iterm_color_variant.py --mode=claude_theme

# New Go CLI usage  
color claude

# Old Python directory theme
python3 ~/iterm_color_variant.py --mode=directory_theme --path="$PWD"

# New Go CLI directory theme
color directory
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt` 
6. Submit a pull request

## License

MIT License - see LICENSE file for details.