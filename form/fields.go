package form

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
)

// Field representa un campo de un formulario.
//
// Es una interfaz a propósito: oculta por completo la librería de terceros
// (huh) detrás del método privado build(). Tus proyectos solo conocen este
// tipo y los constructores (Input, Select, ...), nunca a huh directamente.
type Field interface {
	// build traduce el campo a su equivalente en la librería subyacente.
	// Es privado (minúscula) para que huh.Field no se filtre fuera del paquete.
	build() huh.Field
}

// composeValidate combina la regla "obligatorio" con la validación personalizada.
func composeValidate(c config) func(string) error {
	return func(s string) error {
		if c.required && strings.TrimSpace(s) == "" {
			return errors.New("este campo es obligatorio")
		}
		if c.validate != nil {
			return c.validate(s)
		}
		return nil
	}
}

// ---------------------------------------------------------------------------
// Input: texto de una sola línea.
// ---------------------------------------------------------------------------

type inputField struct {
	title string
	value *string
	cfg   config
}

// Input crea un campo de texto de una línea. El valor escrito se guarda en value.
func Input(title string, value *string, opts ...Option) Field {
	return &inputField{title: title, value: value, cfg: newConfig(opts)}
}

func (f *inputField) build() huh.Field {
	in := huh.NewInput().Title(f.title).Value(f.value).Validate(composeValidate(f.cfg))
	if f.cfg.placeholder != "" {
		in = in.Placeholder(f.cfg.placeholder)
	}
	if f.cfg.description != "" {
		in = in.Description(f.cfg.description)
	}
	return in
}

// ---------------------------------------------------------------------------
// Password: texto oculto (********).
// ---------------------------------------------------------------------------

type passwordField struct {
	title string
	value *string
	cfg   config
}

// Password crea un campo de texto cuyo contenido se oculta al escribir.
func Password(title string, value *string, opts ...Option) Field {
	return &passwordField{title: title, value: value, cfg: newConfig(opts)}
}

func (f *passwordField) build() huh.Field {
	in := huh.NewInput().
		Title(f.title).
		Value(f.value).
		EchoMode(huh.EchoModePassword).
		Validate(composeValidate(f.cfg))
	if f.cfg.placeholder != "" {
		in = in.Placeholder(f.cfg.placeholder)
	}
	if f.cfg.description != "" {
		in = in.Description(f.cfg.description)
	}
	return in
}

// ---------------------------------------------------------------------------
// Text: texto multilínea.
// ---------------------------------------------------------------------------

type textField struct {
	title string
	value *string
	cfg   config
}

// Text crea un campo de texto multilínea (útil para descripciones largas).
func Text(title string, value *string, opts ...Option) Field {
	return &textField{title: title, value: value, cfg: newConfig(opts)}
}

func (f *textField) build() huh.Field {
	t := huh.NewText().Title(f.title).Value(f.value).Validate(composeValidate(f.cfg))
	if f.cfg.placeholder != "" {
		t = t.Placeholder(f.cfg.placeholder)
	}
	if f.cfg.description != "" {
		t = t.Description(f.cfg.description)
	}
	return t
}

// ---------------------------------------------------------------------------
// Select: elegir una opción de una lista (navegación con flechas).
// ---------------------------------------------------------------------------

type selectField struct {
	title   string
	value   *string
	choices []string
}

// Select crea un campo de selección única. El usuario elige una de las
// opciones con las flechas y se guarda en value.
func Select(title string, value *string, choices ...string) Field {
	return &selectField{title: title, value: value, choices: choices}
}

func (f *selectField) build() huh.Field {
	return huh.NewSelect[string]().
		Title(f.title).
		Options(huh.NewOptions(f.choices...)...).
		Value(f.value)
}

// ---------------------------------------------------------------------------
// MultiSelect: elegir varias opciones (espacio para marcar, enter para confirmar).
// ---------------------------------------------------------------------------

type multiSelectField struct {
	title   string
	value   *[]string
	choices []string
}

// MultiSelect crea un campo de selección múltiple. Las opciones marcadas se
// guardan en value.
func MultiSelect(title string, value *[]string, choices ...string) Field {
	return &multiSelectField{title: title, value: value, choices: choices}
}

func (f *multiSelectField) build() huh.Field {
	return huh.NewMultiSelect[string]().
		Title(f.title).
		Options(huh.NewOptions(f.choices...)...).
		Value(f.value)
}

// ---------------------------------------------------------------------------
// Confirm: pregunta sí/no.
// ---------------------------------------------------------------------------

type confirmField struct {
	title string
	value *bool
	cfg   config
}

// Confirm crea una pregunta de confirmación (sí/no). La respuesta se guarda en value.
func Confirm(title string, value *bool, opts ...Option) Field {
	return &confirmField{title: title, value: value, cfg: newConfig(opts)}
}

func (f *confirmField) build() huh.Field {
	c := huh.NewConfirm().Title(f.title).Value(f.value)
	if f.cfg.description != "" {
		c = c.Description(f.cfg.description)
	}
	return c
}

// ---------------------------------------------------------------------------
// Note: solo informativo, no pide entrada.
// ---------------------------------------------------------------------------

type noteField struct {
	title       string
	description string
}

// Note muestra un bloque de texto informativo (sin entrada del usuario).
func Note(title, description string) Field {
	return &noteField{title: title, description: description}
}

func (f *noteField) build() huh.Field {
	return huh.NewNote().Title(f.title).Description(f.description)
}
