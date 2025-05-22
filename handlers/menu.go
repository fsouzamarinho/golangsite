// handlers/menu.go
package handlers

import (
	"andravirtual/config"
	"andravirtual/models"
	"html/template"
	"net/http"
	"strconv"
)

func ListarMenus(w http.ResponseWriter, r *http.Request) {
	rs, err := config.DB.Query("SELECT Codigo, Nome, Link, Ordem FROM menu ORDER BY Ordem")
	if err != nil {
		http.Error(w, "Erro ao consultar menus", 500)
		return
	}
	defer rs.Close()

	menus := []models.Menu{}
	for rs.Next() {
		var m models.Menu
		rs.Scan(&m.Codigo, &m.Nome, &m.Link, &m.Ordem)
		menus = append(menus, m)
	}

	tmpl := template.Must(template.ParseFiles("templates/menus/index.html"))
	tmpl.Execute(w, menus)
}

func NovoMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/menus/form.html"))
	tmpl.Execute(w, models.Menu{})
}

func EditarMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("codigo"))
	var m models.Menu

	err := config.DB.QueryRow("SELECT Codigo, Nome, Link, Ordem FROM menu WHERE Codigo = ?", id).Scan(&m.Codigo, &m.Nome, &m.Link, &m.Ordem)
	if err != nil {
		http.Error(w, "Menu não encontrado", 404)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/menus/form.html"))
	tmpl.Execute(w, m)
}

func SalvarMenu(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	nome := r.FormValue("nome")
	link := r.FormValue("link")
	ordem, _ := strconv.Atoi(r.FormValue("ordem"))
	codigo := r.FormValue("codigo")

	if codigo == "" || codigo == "0" {
		stmt, _ := config.DB.Prepare("INSERT INTO menu (Nome, Link, Ordem) VALUES (?, ?, ?)")
		stmt.Exec(nome, link, ordem)
	} else {
		id, _ := strconv.Atoi(codigo)
		stmt, _ := config.DB.Prepare("UPDATE menu SET Nome = ?, Link = ?, Ordem = ? WHERE Codigo = ?")
		stmt.Exec(nome, link, ordem, id)
	}

	http.Redirect(w, r, "/menus", http.StatusSeeOther)
}

func ExcluirMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("codigo"))
	stmt, _ := config.DB.Prepare("DELETE FROM menu WHERE Codigo = ?")
	stmt.Exec(id)
	http.Redirect(w, r, "/menus", http.StatusSeeOther)
}