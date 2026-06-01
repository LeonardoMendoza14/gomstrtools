# gomstrtools

Conjunto de herramientas en Go para reutilizar entre proyectos.

Por ahora incluye el paquete **`form`**: una API sencilla y estable para crear
**formularios interactivos en consola** (entradas de texto, selección con
flechas, confirmaciones, contraseñas…). Está pensada para **configurar
servicios desde la terminal, incluso con el servidor en marcha**.

Por debajo usa [`charmbracelet/huh`](https://github.com/charmbracelet/huh), pero
esa dependencia queda **totalmente encapsulada**: tus proyectos solo importan
`gomstrtools`, así que si algún día se cambia la librería subyacente tu código
no se entera.

---

## Instalación

```bash
go get github.com/LeonardoMendoza14/gomstrtools
```

Requiere **Go 1.21+** (el módulo está fijado en una versión reciente; mira
[`go.mod`](go.mod)).

---

## Uso rápido

```go
package main

import (
	"log"

	"github.com/LeonardoMendoza14/gomstrtools/form"
)

func main() {
	var puerto, entorno string
	var debug bool

	err := form.New(
		form.Input("Puerto del servidor", &puerto, form.Required()),
		form.Select("Entorno", &entorno, "dev", "staging", "prod"),
		form.Confirm("¿Activar modo debug?", &debug),
	).Run()
	if err != nil {
		log.Fatal(err)
	}

	// Tras Run(), las variables ya contienen lo que escribió el usuario.
	log.Printf("puerto=%s entorno=%s debug=%t", puerto, entorno, debug)
}
```

> **Cómo funciona:** a cada campo le pasas un **puntero** (`&puerto`). Cuando el
> usuario termina el formulario, Go habrá escrito los valores en esas variables.

Hay un ejemplo completo y ejecutable en [`examples/config`](examples/config/main.go):

```bash
go run ./examples/config
```

---

## Tipos de campo

Todos los constructores viven en el paquete `form`.

| Función | Para qué sirve | Guarda en |
|---|---|---|
| `Input(title, *string, ...Option)` | Texto de una línea | `*string` |
| `Password(title, *string, ...Option)` | Texto oculto (`••••`) | `*string` |
| `Text(title, *string, ...Option)` | Texto multilínea | `*string` |
| `Select(title, *string, choices...)` | Elegir **una** opción (flechas) | `*string` |
| `MultiSelect(title, *[]string, choices...)` | Elegir **varias** opciones | `*[]string` |
| `Confirm(title, *bool, ...Option)` | Pregunta sí/no | `*bool` |
| `Note(title, description)` | Texto informativo (sin entrada) | — |

### Ejemplos por campo

```go
var nombre, bio, region, clave string
var regiones []string
var aceptar bool

form.Input("Nombre", &nombre, form.Placeholder("ej. api-gateway"))

form.Password("Token", &clave, form.Required())

form.Text("Descripción", &bio, form.Description("Puede ocupar varias líneas"))

form.Select("Región principal", &region, "eu-west", "us-east", "sa-east")

form.MultiSelect("Réplicas en", &regiones, "eu-west", "us-east", "sa-east")

form.Confirm("¿Aplicar cambios?", &aceptar)

form.Note("Aviso", "Estos cambios se aplican en caliente.")
```

---

## Opciones (`...Option`)

Los campos de texto (`Input`, `Password`, `Text`) y `Confirm` aceptan opciones
al final:

| Opción | Efecto |
|---|---|
| `Required()` | El campo no puede quedar vacío |
| `Placeholder(texto)` | Texto guía dentro del campo vacío |
| `Description(texto)` | Línea descriptiva bajo el título |
| `Validate(func(string) error)` | Validación personalizada |

```go
import (
	"errors"
	"strconv"
)

form.Input("Puerto", &puerto,
	form.Required(),
	form.Placeholder("8080"),
	form.Validate(func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return errors.New("debe ser un número")
		}
		return nil
	}),
)
```

> `Select` y `MultiSelect` no usan `Option`: sus opciones son los `choices` que
> pasas como argumentos variádicos.

---

## Formularios de varios pasos

Para dividir la configuración en pantallas, usa `NewSteps` + `NewStep`:

```go
form.NewSteps(
	form.NewStep(
		form.Input("Host", &host, form.Required()),
		form.Input("Puerto", &puerto, form.Required()),
	),
	form.NewStep(
		form.Select("Entorno", &entorno, "dev", "prod"),
		form.Confirm("¿Guardar?", &guardar),
	),
).Run()
```

Cada `NewStep` es una página independiente; el usuario avanza de una a otra.

---

## Temas

Cambia el aspecto visual con `WithTheme`:

```go
form.New(/* ... */).WithTheme(form.ThemeDracula).Run()
```

Temas disponibles: `ThemeCharm` (por defecto), `ThemeDracula`, `ThemeBase16`,
`ThemeBase`.

---

## Manejo de errores y cancelación

`Run()` devuelve un `error`:

- `nil` → el usuario completó el formulario; los punteros ya tienen los valores.
- distinto de `nil` → fallo de ejecución **o** el usuario canceló (p. ej. `Ctrl+C`).

```go
if err := f.Run(); err != nil {
	log.Printf("formulario no completado: %v", err)
	return
}
```

---

## Estructura del módulo

```
gomstrtools/
├── go.mod
├── README.md
├── form/
│   ├── form.go        // Form, New, NewSteps, Run, WithTheme
│   ├── fields.go      // Input, Password, Text, Select, MultiSelect, Confirm, Note
│   ├── options.go     // Required, Placeholder, Description, Validate
│   ├── theme.go       // Temas
│   └── form_test.go
└── examples/
    └── config/
        └── main.go    // ejemplo ejecutable
```

`huh` solo se importa dentro del paquete `form`. El resto del mundo —incluidos
tus proyectos— solo conoce los tipos públicos de `gomstrtools/form`.

---

## Desarrollo

```bash
go test ./...      # ejecutar tests
go build ./...     # compilar todo
go run ./examples/config   # probar el ejemplo
```
