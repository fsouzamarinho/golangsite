package handlers

import (
	"andravirtual/config"
	"andravirtual/models"
	"html/template"
	"net/http"
	"strconv"
)

type EditoriaFormData struct {
	Editoria models.Editoria
	Menus    []models.Menu
}

func ListarEditorias(w http.ResponseWriter, r *http.Request) {
	rs, err := config.DB.Query("SELECT Codigo, Nome, Cod_menu FROM editorias ORDER BY Codigo ASC")
	if err != nil {
		http.Error(w, "Erro ao consultar editorias", 500)
		return
	}
	defer rs.Close()

	editorias := []models.Editoria{}
	for rs.Next() {
		var e models.Editoria
		rs.Scan(&e.Codigo, &e.Nome, &e.CodMenu)
		editorias = append(editorias, e)
	}

	tmpl := template.Must(template.ParseFiles("templates/editorias/index.html"))
	tmpl.Execute(w, editorias)
}

func NovaEditoria(w http.ResponseWriter, r *http.Request) {
	menus := getMenus()
	data := EditoriaFormData{Editoria: models.Editoria{}, Menus: menus}
	tmpl := template.Must(template.ParseFiles("templates/editorias/form.html"))
	tmpl.Execute(w, data)
}

func EditarEditoria(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("codigo"))
	var e models.Editoria

	err := config.DB.QueryRow("SELECT Codigo, Nome, Cod_menu FROM editorias WHERE Codigo = ?", id).Scan(&e.Codigo, &e.Nome, &e.CodMenu)
	if err != nil {
		http.Error(w, "Editoria não encontrada", 404)
		return
	}

	menus := getMenus()
	data := EditoriaFormData{Editoria: e, Menus: menus}
	tmpl := template.Must(template.ParseFiles("templates/editorias/form.html"))
	tmpl.Execute(w, data)
}

func SalvarEditoria(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	nome := r.FormValue("nome")
	codMenu, _ := strconv.Atoi(r.FormValue("codmenu"))
	codigo := r.FormValue("codigo")

	if codigo == "" || codigo == "0" {
		stmt, _ := config.DB.Prepare("INSERT INTO editorias (Nome, Cod_menu) VALUES (?, ?)")
		stmt.Exec(nome, codMenu)
	} else {
		id, _ := strconv.Atoi(codigo)
		stmt, _ := config.DB.Prepare("UPDATE editorias SET Nome = ?, Cod_menu = ? WHERE Codigo = ?")
		stmt.Exec(nome, codMenu, id)
	}

	http.Redirect(w, r, "/editorias", http.StatusSeeOther)
}

func ExcluirEditoria(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("codigo"))
	stmt, _ := config.DB.Prepare("DELETE FROM editorias WHERE Codigo = ?")
	stmt.Exec(id)
	http.Redirect(w, r, "/editorias", http.StatusSeeOther)
}

func getMenus() []models.Menu {
	rs, err := config.DB.Query("SELECT Codigo, Nome FROM menu ORDER BY Nome")
	if err != nil {
		return nil
	}
	defer rs.Close()

	menus := []models.Menu{}
	for rs.Next() {
		var m models.Menu
		rs.Scan(&m.Codigo, &m.Nome)
		menus = append(menus, m)
	}
	return menus
}
