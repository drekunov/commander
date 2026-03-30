package config

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vanng822/css"
)

// StyleSheet represents parsed CSS as a map of selector -> properties.
type StyleSheet map[string]map[string]string

// ParseCSS parses a CSS string into a StyleSheet using the vanng822/css parser.
func ParseCSS(data string) StyleSheet {
	ss := make(StyleSheet)
	parsed := css.Parse(data)

	for _, rule := range parsed.GetCSSRuleList() {
		selector := rule.Style.Selector.Text()
		props := make(map[string]string)
		for _, decl := range rule.Style.Styles {
			props[decl.Property] = decl.Value.Text()
		}
		if len(props) > 0 {
			ss[selector] = props
		}
	}

	return ss
}

// ApplyStyle builds a lipgloss.Style from CSS properties.
func ApplyStyle(props map[string]string) lipgloss.Style {
	s := lipgloss.NewStyle()

	if v, ok := props["color"]; ok {
		s = s.Foreground(lipgloss.Color(v))
	}
	if v, ok := props["background-color"]; ok {
		s = s.Background(lipgloss.Color(v))
	}
	if v, ok := props["text-decoration"]; ok && v == "underline" {
		s = s.Underline(true)
	}
	if v, ok := props["padding"]; ok {
		top, right, bottom, left := parseBox(v)
		s = s.Padding(top, right, bottom, left)
	}
	if v, ok := props["margin"]; ok {
		top, right, bottom, left := parseBox(v)
		s = s.Margin(top, right, bottom, left)
	}
	if v, ok := props["border-style"]; ok {
		border := parseBorder(v)
		s = s.Border(border).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
	}
	if v, ok := props["border-color"]; ok {
		s = s.BorderForeground(lipgloss.Color(v))
	}

	return s
}

// parseBox parses CSS box values (1, 2, or 4 integers).
func parseBox(v string) (top, right, bottom, left int) {
	parts := strings.Fields(v)
	nums := make([]int, len(parts))
	for i, p := range parts {
		n, _ := strconv.Atoi(p)
		nums[i] = n
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

func parseBorder(v string) lipgloss.Border {
	switch strings.ToLower(v) {
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
