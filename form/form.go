// Package form ofrece una API sencilla y estable para construir formularios
// interactivos en consola (entradas de texto, selección con flechas,
// confirmaciones, contraseñas, etc.).
//
// Está pensada para configurar servicios desde la terminal, incluso con el
// servidor en marcha. Por debajo usa github.com/charmbracelet/huh, pero esa
// dependencia queda totalmente encapsulada: tus proyectos solo importan este
// paquete, así que cambiar la librería subyacente no rompería tu código.
//
// Ejemplo mínimo:
//
//	var puerto, entorno string
//	var debug bool
//
//	err := form.New(
//	    form.Input("Puerto del servidor", &puerto, form.Required()),
//	    form.Select("Entorno", &entorno, "dev", "staging", "prod"),
//	    form.Confirm("¿Activar modo debug?", &debug),
//	).Run()
//	if err != nil {
//	    log.Fatal(err)
//	}
package form

import "github.com/charmbracelet/huh"

// Form es un formulario listo para ejecutarse.
type Form struct {
	groups []*huh.Group
	theme  *huh.Theme
}

// New crea un formulario de una sola página con los campos indicados,
// en el orden en que se pasan.
func New(fields ...Field) *Form {
	return &Form{groups: []*huh.Group{toGroup(fields)}}
}

// Step agrupa un conjunto de campos en una página/paso del formulario.
type Step struct {
	fields []Field
}

// NewStep crea un paso con los campos indicados. Se usa junto con NewSteps.
func NewStep(fields ...Field) Step {
	return Step{fields: fields}
}

// NewSteps crea un formulario de varias páginas. Cada Step se muestra como
// una pantalla independiente y el usuario avanza de una a otra.
//
//	form.NewSteps(
//	    form.NewStep(
//	        form.Input("Host", &host),
//	        form.Input("Puerto", &puerto),
//	    ),
//	    form.NewStep(
//	        form.Select("Entorno", &entorno, "dev", "prod"),
//	        form.Confirm("¿Guardar?", &guardar),
//	    ),
//	).Run()
func NewSteps(steps ...Step) *Form {
	groups := make([]*huh.Group, len(steps))
	for i, s := range steps {
		groups[i] = toGroup(s.fields)
	}
	return &Form{groups: groups}
}

// toGroup traduce los campos de la librería propia a un grupo de huh.
func toGroup(fields []Field) *huh.Group {
	built := make([]huh.Field, len(fields))
	for i, f := range fields {
		built[i] = f.build()
	}
	return huh.NewGroup(built...)
}

// WithTheme cambia el tema visual del formulario y devuelve el mismo Form
// para poder encadenar la llamada con Run.
//
//	form.New(...).WithTheme(form.ThemeDracula).Run()
func (f *Form) WithTheme(t Theme) *Form {
	f.theme = t.huh()
	return f
}

// Run muestra el formulario y bloquea hasta que el usuario lo completa o lo
// cancela. Devuelve un error si la ejecución falla o si el usuario aborta
// (por ejemplo con Ctrl+C). Cuando termina con éxito, los punteros que pasaste
// a cada campo contienen los valores introducidos.
func (f *Form) Run() error {
	hf := huh.NewForm(f.groups...)
	if f.theme != nil {
		hf = hf.WithTheme(f.theme)
	}
	return hf.Run()
}
