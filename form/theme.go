package form

import "github.com/charmbracelet/huh"

// Theme representa un tema visual del formulario. Se expone como enum propio
// para no obligar a tus proyectos a importar huh solo por el tema.
type Theme int

const (
	// ThemeCharm es el tema por defecto: colorido y moderno.
	ThemeCharm Theme = iota
	// ThemeDracula usa la conocida paleta Dracula.
	ThemeDracula
	// ThemeBase16 es un tema de contraste alto basado en Base16.
	ThemeBase16
	// ThemeBase es un tema sobrio, casi sin color (útil en terminales limitadas).
	ThemeBase
)

// huh traduce el tema propio al *huh.Theme correspondiente.
func (t Theme) huh() *huh.Theme {
	switch t {
	case ThemeDracula:
		return huh.ThemeDracula()
	case ThemeBase16:
		return huh.ThemeBase16()
	case ThemeBase:
		return huh.ThemeBase()
	default:
		return huh.ThemeCharm()
	}
}
