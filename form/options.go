package form

// config almacena la configuración opcional de un campo.
// Es interna: el exterior nunca la manipula directamente, solo a través
// de las funciones Option de abajo.
type config struct {
	placeholder string
	description string
	validate    func(string) error
	required    bool
}

// Option configura un campo de forma opcional.
//
// Se usan como argumentos variádicos en los constructores de campos:
//
//	form.Input("Nombre", &nombre,
//	    form.Required(),
//	    form.Placeholder("ej. servidor-1"),
//	)
type Option func(*config)

// Placeholder muestra un texto de ayuda dentro del campo cuando está vacío.
func Placeholder(text string) Option {
	return func(c *config) { c.placeholder = text }
}

// Description muestra una línea descriptiva debajo del título del campo.
func Description(text string) Option {
	return func(c *config) { c.description = text }
}

// Validate registra una función que valida el valor introducido.
// Debe devolver un error (que se mostrará al usuario) si el valor no es válido,
// o nil si es correcto.
//
//	form.Input("Puerto", &puerto, form.Validate(func(s string) error {
//	    if _, err := strconv.Atoi(s); err != nil {
//	        return errors.New("debe ser un número")
//	    }
//	    return nil
//	}))
func Validate(fn func(string) error) Option {
	return func(c *config) { c.validate = fn }
}

// Required marca el campo como obligatorio: no permite continuar si está vacío.
func Required() Option {
	return func(c *config) { c.required = true }
}

// newConfig aplica todas las opciones sobre una config vacía.
func newConfig(opts []Option) config {
	var c config
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
