package handlers

import (
	"andravirtual/cache"
	"andravirtual/config"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	var err error
	templates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Erro ao carregar templates:", err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()

	navbarHTML, err := cache.GetNavbarHTML(db, templates)
	if err != nil {
		http.Error(w, "Erro ao gerar navbar", http.StatusInternalServerError)
		log.Println("Erro ao gerar navbar:", err)
		return
	}

	data := map[string]interface{}{
		"Navbar": template.HTML(navbarHTML),
	}

	err = templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
		log.Println("Erro ao renderizar template:", err)
	}
}
