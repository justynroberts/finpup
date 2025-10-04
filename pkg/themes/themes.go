package themes

import "github.com/gdamore/tcell/v2"

type Theme struct {
	Name       string
	Background tcell.Color
	Foreground tcell.Color
	StatusBG   tcell.Color
	StatusFG   tcell.Color
	LineNumFG  tcell.Color
	SelectionBG tcell.Color
	SelectionFG tcell.Color
	KeywordFG   tcell.Color
	StringFG    tcell.Color
	CommentFG   tcell.Color
	NumberFG    tcell.Color
}

var (
	Dark = Theme{
		Name:       "dark",
		Background: tcell.NewRGBColor(20, 20, 20),
		Foreground: tcell.NewRGBColor(220, 220, 220),
		StatusBG:   tcell.NewRGBColor(40, 80, 140),
		StatusFG:   tcell.NewRGBColor(255, 255, 255),
		LineNumFG:  tcell.NewRGBColor(120, 120, 120),
		SelectionBG: tcell.NewRGBColor(70, 130, 180),
		SelectionFG: tcell.NewRGBColor(255, 255, 255),
		KeywordFG:   tcell.NewRGBColor(255, 215, 0),
		StringFG:    tcell.NewRGBColor(152, 251, 152),
		CommentFG:   tcell.NewRGBColor(140, 140, 140),
		NumberFG:    tcell.NewRGBColor(255, 105, 180),
	}

	Light = Theme{
		Name:       "light",
		Background: tcell.NewRGBColor(250, 250, 250),
		Foreground: tcell.NewRGBColor(40, 40, 40),
		StatusBG:   tcell.NewRGBColor(173, 216, 230),
		StatusFG:   tcell.NewRGBColor(0, 0, 0),
		LineNumFG:  tcell.NewRGBColor(150, 150, 150),
		SelectionBG: tcell.NewRGBColor(200, 230, 255),
		SelectionFG: tcell.NewRGBColor(0, 0, 0),
		KeywordFG:   tcell.NewRGBColor(0, 0, 200),
		StringFG:    tcell.NewRGBColor(0, 128, 0),
		CommentFG:   tcell.NewRGBColor(128, 128, 128),
		NumberFG:    tcell.NewRGBColor(148, 0, 211),
	}

	Monokai = Theme{
		Name:       "monokai",
		Background: tcell.NewRGBColor(39, 40, 34),
		Foreground: tcell.NewRGBColor(248, 248, 242),
		StatusBG:   tcell.NewRGBColor(73, 72, 62),
		StatusFG:   tcell.NewRGBColor(248, 248, 242),
		LineNumFG:  tcell.NewRGBColor(144, 144, 140),
		SelectionBG: tcell.NewRGBColor(73, 72, 62),
		SelectionFG: tcell.NewRGBColor(248, 248, 242),
		KeywordFG:   tcell.NewRGBColor(249, 38, 114),
		StringFG:    tcell.NewRGBColor(230, 219, 116),
		CommentFG:   tcell.NewRGBColor(117, 113, 94),
		NumberFG:    tcell.NewRGBColor(174, 129, 255),
	}

	Solarized = Theme{
		Name:       "solarized",
		Background: tcell.NewRGBColor(0, 43, 54),
		Foreground: tcell.NewRGBColor(131, 148, 150),
		StatusBG:   tcell.NewRGBColor(7, 54, 66),
		StatusFG:   tcell.NewRGBColor(147, 161, 161),
		LineNumFG:  tcell.NewRGBColor(88, 110, 117),
		SelectionBG: tcell.NewRGBColor(7, 54, 66),
		SelectionFG: tcell.NewRGBColor(147, 161, 161),
		KeywordFG:   tcell.NewRGBColor(38, 139, 210),
		StringFG:    tcell.NewRGBColor(42, 161, 152),
		CommentFG:   tcell.NewRGBColor(88, 110, 117),
		NumberFG:    tcell.NewRGBColor(211, 54, 130),
	}
)

var AllThemes = []Theme{Dark, Light, Monokai, Solarized}

func GetTheme(name string) Theme {
	for _, theme := range AllThemes {
		if theme.Name == name {
			return theme
		}
	}
	return Dark
}

func NextTheme(current string) Theme {
	for i, theme := range AllThemes {
		if theme.Name == current {
			return AllThemes[(i+1)%len(AllThemes)]
		}
	}
	return Dark
}
