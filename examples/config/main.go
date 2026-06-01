// Ejemplo ejecutable: configurar un servidor desde la consola.
//
// Ejecútalo con:
//
//	go run ./examples/config
package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/LeonardoMendoza14/gomstrtools/form"
)

func main() {
	var (
		host      string
		puerto    string
		entorno   string
		token     string
		servicios []string
		debug     bool
	)

	err := form.New(
		form.Note("Configuración del servidor", "Rellena los campos para aplicar la configuración."),
		form.Input("Host", &host,
			form.Required(),
			form.Placeholder("0.0.0.0"),
		),
		form.Input("Puerto", &puerto,
			form.Required(),
			form.Validate(func(s string) error {
				if _, err := strconv.Atoi(s); err != nil {
					return errors.New("el puerto debe ser un número")
				}
				return nil
			}),
		),
		form.Select("Entorno", &entorno, "dev", "staging", "prod"),
		form.Password("Token de admin", &token, form.Required()),
		form.MultiSelect("Servicios a iniciar", &servicios, "api", "worker", "scheduler", "metrics"),
		form.Confirm("¿Activar modo debug?", &debug),
	).WithTheme(form.ThemeCharm).Run()

	if err != nil {
		fmt.Println("Formulario cancelado:", err)
		return
	}

	fmt.Println("\n--- Configuración aplicada ---")
	fmt.Printf("Host:      %s\n", host)
	fmt.Printf("Puerto:    %s\n", puerto)
	fmt.Printf("Entorno:   %s\n", entorno)
	fmt.Printf("Servicios: %v\n", servicios)
	fmt.Printf("Debug:     %t\n", debug)
}
