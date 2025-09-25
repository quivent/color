package internal

import (
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// RGB represents RGB color values (0-255)
type RGB struct {
	R, G, B uint8
}

// HSV represents HSV color values
type HSV struct {
	H, S, V float64
}

// ColorManager handles terminal color operations
type ColorManager struct {
	rng         *rand.Rand
	persistence *PersistenceManager
}

// NewColorManager creates a new color manager
func NewColorManager() *ColorManager {
	return &ColorManager{
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		persistence: NewPersistenceManager(),
	}
}

// RGBToITerm converts RGB (0-255) to iTerm2 color values (0-65535)
func (c *ColorManager) RGBToITerm(r, g, b uint8) (uint16, uint16, uint16) {
	return uint16(r) * 257, uint16(g) * 257, uint16(b) * 257
}

// HSVToRGB converts HSV to RGB
func (c *ColorManager) HSVToRGB(h, s, v float64) RGB {
	// Ensure h is in [0, 1)
	for h >= 1.0 {
		h -= 1.0
	}
	for h < 0.0 {
		h += 1.0
	}

	// Clamp s and v to [0, 1]
	s = math.Max(0, math.Min(1, s))
	v = math.Max(0, math.Min(1, v))

	chroma := v * s
	x := chroma * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := v - chroma

	var rPrime, gPrime, bPrime float64

	switch {
	case h < 1.0/6:
		rPrime, gPrime, bPrime = chroma, x, 0
	case h < 2.0/6:
		rPrime, gPrime, bPrime = x, chroma, 0
	case h < 3.0/6:
		rPrime, gPrime, bPrime = 0, chroma, x
	case h < 4.0/6:
		rPrime, gPrime, bPrime = 0, x, chroma
	case h < 5.0/6:
		rPrime, gPrime, bPrime = x, 0, chroma
	default:
		rPrime, gPrime, bPrime = chroma, 0, x
	}

	return RGB{
		R: uint8((rPrime + m) * 255),
		G: uint8((gPrime + m) * 255),
		B: uint8((bPrime + m) * 255),
	}
}

// RGBToHSV converts RGB to HSV
func (c *ColorManager) RGBToHSV(rgb RGB) HSV {
	r := float64(rgb.R) / 255.0
	g := float64(rgb.G) / 255.0
	b := float64(rgb.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	diff := max - min

	var h, s, v float64

	// Value
	v = max

	// Saturation
	if max == 0 {
		s = 0
	} else {
		s = diff / max
	}

	// Hue
	if diff == 0 {
		h = 0
	} else {
		switch max {
		case r:
			h = (g-b)/diff + (func() float64 {
				if g < b {
					return 6
				}
				return 0
			})()
		case g:
			h = (b-r)/diff + 2
		case b:
			h = (r-g)/diff + 4
		}
		h /= 6
	}

	return HSV{H: h, S: s, V: v}
}

// GetCurrentColor gets current background color from iTerm2
func (c *ColorManager) GetCurrentColor() (RGB, error) {
	script := `
	tell application "iTerm2"
		tell current session of current tab of current window
			get background color
		end tell
	end tell
	`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		// Return default dark gray
		return RGB{30, 30, 30}, nil
	}

	// Parse output like "0, 0, 0"
	colorStr := strings.TrimSpace(string(output))
	parts := strings.Split(colorStr, ",")
	if len(parts) != 3 {
		return RGB{30, 30, 30}, nil
	}

	var values [3]uint16
	for i, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return RGB{30, 30, 30}, nil
		}
		values[i] = uint16(val)
	}

	// Convert from iTerm2 values (0-65535) to RGB (0-255)
	return RGB{
		R: uint8(values[0] / 257),
		G: uint8(values[1] / 257),
		B: uint8(values[2] / 257),
	}, nil
}

// SetITermColor sets iTerm2 background color
func (c *ColorManager) SetITermColor(rgb RGB) error {
	iR, iG, iB := c.RGBToITerm(rgb.R, rgb.G, rgb.B)

	script := fmt.Sprintf(`
	tell application "iTerm2"
		tell current session of current tab of current window
			set background color to {%d, %d, %d}
		end tell
	end tell
	`, iR, iG, iB)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// GenerateClaudeTheme generates Claude-specific color theme
func (c *ColorManager) GenerateClaudeTheme() RGB {
	// Check if we have a recent Claude color stored
	if c.persistence != nil && c.persistence.IsEnabled() {
		if color, found := c.persistence.GetLastClaudeColor(); found {
			return color
		}
	}

	// Generate new Claude color
	baseHues := []float64{0.6, 0.75, 0.85} // Blue to purple range
	hue := baseHues[c.rng.Intn(len(baseHues))]
	saturation := 0.4 + c.rng.Float64()*0.4 // 0.4-0.8 (more saturated)
	value := 0.25 + c.rng.Float64()*0.15    // 0.25-0.4 (brighter for visibility)

	color := c.HSVToRGB(hue, saturation, value)

	// Store the new color
	if c.persistence != nil && c.persistence.IsEnabled() {
		c.persistence.SetLastClaudeColor(color)
	}

	return color
}

// GenerateDirectoryTheme generates consistent color for directory based on path hash
func (c *ColorManager) GenerateDirectoryTheme(directoryPath string) RGB {
	if directoryPath == "" {
		var err error
		directoryPath, err = os.Getwd()
		if err != nil {
			directoryPath = "/tmp"
		}
	}

	// Check if we have this directory color stored
	if c.persistence != nil && c.persistence.IsEnabled() {
		if color, found := c.persistence.GetDirectoryColor(directoryPath); found {
			return color
		}
	}

	// Generate new directory color based on path hash
	hash := md5.Sum([]byte(directoryPath))
	hashStr := fmt.Sprintf("%x", hash)

	// Convert hash to hue (0-1)
	hueHex := hashStr[:8]
	hueInt, _ := strconv.ParseUint(hueHex, 16, 64)
	hue := float64(hueInt) / float64(0xFFFFFFFF)

	// Use consistent saturation and value for better visibility
	satHex := hashStr[8:10]
	satInt, _ := strconv.ParseUint(satHex, 16, 8)
	saturation := 0.5 + (float64(satInt)/255.0)*0.3 // 0.5-0.8 (more saturated)

	valHex := hashStr[10:12]
	valInt, _ := strconv.ParseUint(valHex, 16, 8)
	value := 0.25 + (float64(valInt)/255.0)*0.2 // 0.25-0.45 (brighter)

	color := c.HSVToRGB(hue, saturation, value)

	// Store the new directory color
	if c.persistence != nil && c.persistence.IsEnabled() {
		c.persistence.SetDirectoryColor(directoryPath, color)
	}

	return color
}

// GenerateVariant generates color variant based on current color
func (c *ColorManager) GenerateVariant(baseColor RGB, mode string) RGB {
	hsv := c.RGBToHSV(baseColor)

	switch mode {
	case "hue_shift":
		hsv.H += 0.15 + c.rng.Float64()*0.2 // +0.15 to +0.35 (more dramatic)
		if hsv.H >= 1.0 {
			hsv.H -= 1.0
		}
	case "brightness":
		change := -0.4 + c.rng.Float64()*0.8 // -0.4 to +0.4 (more dramatic)
		hsv.V = math.Max(0.2, math.Min(0.8, hsv.V+change))
	case "saturation":
		change := -0.4 + c.rng.Float64()*0.8 // -0.4 to +0.4 (more dramatic)
		hsv.S = math.Max(0.2, math.Min(0.9, hsv.S+change))
	case "complement":
		hsv.H += 0.5
		if hsv.H >= 1.0 {
			hsv.H -= 1.0
		}
		// Also boost saturation and value for complement
		hsv.S = math.Min(0.9, hsv.S+0.2)
		hsv.V = math.Min(0.7, hsv.V+0.1)
	default:
		// Random mode selection if mode is unknown
		modes := []string{"hue_shift", "brightness", "saturation", "complement"}
		return c.GenerateVariant(baseColor, modes[c.rng.Intn(len(modes))])
	}

	return c.HSVToRGB(hsv.H, hsv.S, hsv.V)
}

// GetPersistenceStatus returns the status of the persistence system
func (c *ColorManager) GetPersistenceStatus() string {
	if c.persistence == nil {
		return "❌ Persistence system not initialized"
	}
	return c.persistence.GetConnectionStatus()
}

// GetColorHistory returns recent color history
func (c *ColorManager) GetColorHistory(limit int) ([]ColorEntry, error) {
	if c.persistence == nil {
		return []ColorEntry{}, nil
	}
	return c.persistence.GetColorHistory(limit)
}

// ClearColorCache clears all stored colors
func (c *ColorManager) ClearColorCache() error {
	if c.persistence == nil {
		return nil
	}
	return c.persistence.ClearColorCache()
}

// Close closes persistence connections
func (c *ColorManager) Close() error {
	if c.persistence != nil {
		return c.persistence.Close()
	}
	return nil
}