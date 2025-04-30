package temp

import (
	"log"
	"text/template"
)

var Temp *template.Template

func InitTemplates() {
	var err error
	Temp, err = template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}
