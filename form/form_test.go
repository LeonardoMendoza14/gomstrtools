package form

import (
	"errors"
	"testing"
)

// TestBuildFields comprueba que todos los tipos de campo se construyen sin
// errores ni panics. No ejecuta el formulario (eso requeriría una terminal real).
func TestBuildFields(t *testing.T) {
	var s string
	var b bool
	var ms []string

	fields := []Field{
		Input("Texto", &s, Required(), Placeholder("escribe aquí")),
		Password("Clave", &s),
		Text("Multilínea", &s, Description("varias líneas")),
		Select("Único", &s, "a", "b", "c"),
		MultiSelect("Varios", &ms, "x", "y", "z"),
		Confirm("¿Seguro?", &b),
		Note("Aviso", "solo informativo"),
	}

	for i, f := range fields {
		if f.build() == nil {
			t.Fatalf("el campo %d devolvió un huh.Field nil", i)
		}
	}
}

// TestNewSingleGroup verifica que New crea un formulario de un solo grupo.
func TestNewSingleGroup(t *testing.T) {
	var s string
	f := New(Input("a", &s), Input("b", &s))
	if len(f.groups) != 1 {
		t.Fatalf("esperaba 1 grupo, obtuve %d", len(f.groups))
	}
}

// TestNewSteps verifica que cada Step se convierte en un grupo independiente.
func TestNewSteps(t *testing.T) {
	var s string
	f := NewSteps(
		NewStep(Input("a", &s)),
		NewStep(Input("b", &s)),
		NewStep(Input("c", &s)),
	)
	if len(f.groups) != 3 {
		t.Fatalf("esperaba 3 grupos, obtuve %d", len(f.groups))
	}
}

// TestRequiredValidation comprueba la regla de campo obligatorio.
func TestRequiredValidation(t *testing.T) {
	v := composeValidate(newConfig([]Option{Required()}))
	if v("") == nil {
		t.Error("un campo obligatorio vacío debería devolver error")
	}
	if v("  ") == nil {
		t.Error("un campo obligatorio con solo espacios debería devolver error")
	}
	if err := v("valor"); err != nil {
		t.Errorf("un campo obligatorio con valor no debería fallar: %v", err)
	}
}

// TestCustomValidation comprueba que la validación personalizada se respeta.
func TestCustomValidation(t *testing.T) {
	soloNumeros := Validate(func(s string) error {
		if s != "123" {
			return errors.New("inválido")
		}
		return nil
	})
	v := composeValidate(newConfig([]Option{soloNumeros}))
	if err := v("123"); err != nil {
		t.Errorf("no debería fallar con valor válido: %v", err)
	}
	if v("abc") == nil {
		t.Error("debería fallar con valor inválido")
	}
}

// TestWithTheme verifica que WithTheme asigna un tema y permite encadenar.
func TestWithTheme(t *testing.T) {
	var s string
	f := New(Input("a", &s)).WithTheme(ThemeDracula)
	if f.theme == nil {
		t.Fatal("WithTheme debería haber asignado un tema")
	}
}
