package config

import (
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vanng822/css"
)

// StyleSheet represents parsed CSS as a map of selector -> properties.
type StyleSheet map[string]map[string]string

// ParseCSS parses a CSS string into a StyleSheet using the vanng822/css parser.
func ParseCSS(data string) StyleSheet {
	sheet := make(StyleSheet)
	parsed := css.Parse(data)

	for _, rule := range parsed.GetCSSRuleList() {
		selector := rule.Style.Selector.Text()

		props := make(map[string]string)
		for _, decl := range rule.Style.Styles {
			props[decl.Property] = decl.Value.Text()
		}

		if len(props) > 0 {
			sheet[selector] = props
		}
	}

	return sheet
}

// ApplyStyle builds a lipgloss.Style from CSS properties.
func ApplyStyle(props map[string]string) lipgloss.Style {
	style := lipgloss.NewStyle()

	if val, ok := props["color"]; ok {
		style = style.Foreground(lipgloss.Color(val))
	}

	if val, ok := props["background-color"]; ok {
		style = style.Background(lipgloss.Color(val))
	}

	if val, ok := props["text-decoration"]; ok && val == "underline" {
		style = style.Underline(true)
	}

	if val, ok := props["padding"]; ok {
		top, right, bottom, left := parseBox(val)
		style = style.Padding(top, right, bottom, left)
	}

	if val, ok := props["margin"]; ok {
		top, right, bottom, left := parseBox(val)
		style = style.Margin(top, right, bottom, left)
	}

	if val, ok := props["border-style"]; ok {
		border := parseBorder(val)
		style = style.Border(border).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
	}

	if val, ok := props["border-color"]; ok {
		style = style.BorderForeground(lipgloss.Color(val))
	}

	return style
}

// parseBox parses CSS box values (1, 2, or 4 integers).
func parseBox(val string) (int, int, int, int) {
	parts := strings.Fields(val)

	nums := make([]int, 0, len(parts))
	for _, part := range parts {
		parsed, err := strconv.Atoi(part)
		if err != nil {
			log.Printf("config: could not parse box value %q: %v", part, err)

			continue
		}

		nums = append(nums, parsed)
	}

	switch len(nums) {
	case 1:
		return nums[0], nums[0], nums[0], nums[0]
	case 2:
		return nums[0], nums[1], nums[0], nums[1]
	case 3:
		return nums[0], nums[1], nums[2], nums[1]
	case 4:
		return nums[0], nums[1], nums[2], nums[3]
	}

	return 0, 0, 0, 0
}

func parseBorder(val string) lipgloss.Border {
	switch strings.ToLower(val) {
	case "rounded":
		return lipgloss.RoundedBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "hidden", "none":
		return lipgloss.HiddenBorder()
	default:
		return lipgloss.NormalBorder()
	}
}
