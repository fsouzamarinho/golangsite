package handlers

import (
	"andravirtual/config"
	"andravirtual/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"andravirtual/utils"
)

func ListarAutores(w http.ResponseWriter, r *http.Request) {
	rs, err := config.DB.Query("SELECT id, Nome, Facebook, Link, Perfil, Foto, Twitter, Instagram, Youtube, Linkedin FROM author ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Erro ao consultar autores", 500)
		return
	}
	defer rs.Close()

	autores := []models.Author{}
	for rs.Next() {
		var a models.Author
		rs.Scan(&a.ID, &a.Nome, &a.Facebook, &a.Link, &a.Perfil, &a.Foto, &a.Twitter, &a.Instagram, &a.Youtube, &a.Linkedin)
		autores = append(autores, a)
	}

	tmpl := template.Must(template.ParseFiles("templates/autores/index.html"))
	tmpl.Execute(w, autores)
}

func NovoAutor(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/autores/form.html"))
	tmpl.Execute(w, models.Author{})
}

func EditarAutor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var a models.Author
	err := config.DB.QueryRow("SELECT id, Nome, Facebook, Link, Perfil, Foto, Twitter, Instagram, Youtube, Linkedin FROM author WHERE id = ?", id).
		Scan(&a.ID, &a.Nome, &a.Facebook, &a.Link, &a.Perfil, &a.Foto, &a.Twitter, &a.Instagram, &a.Youtube, &a.Linkedin)
	if err != nil {
		http.Error(w, "Autor não encontrado", 404)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/autores/form.html"))
	tmpl.Execute(w, a)
}

func SalvarAutor(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	foto := ""
	file, handler, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(handler.Filename)
		path := "static/uploads/autores/" + filename
		out, err := os.Create(path)
		if err == nil {
			defer out.Close()
			io.Copy(out, file)
			foto = "autores/" + filename
		}
	}

	autor := models.Author{
		ID:        utils.AtoiOrZero(r.FormValue("id")),
		Nome:      r.FormValue("nome"),
		Facebook:  r.FormValue("facebook"),
		Link:      r.FormValue("link"),
		Perfil:    r.FormValue("perfil"),
		Foto:      foto,
		Twitter:   r.FormValue("twitter"),
		Instagram: r.FormValue("instagram"),
		Youtube:   r.FormValue("youtube"),
		Linkedin:  r.FormValue("linkedin"),
	}

	if autor.ID == 0 {
		stmt, _ := config.DB.Prepare("INSERT INTO author (Nome, Facebook, Link, Perfil, Foto, Twitter, Instagram, Youtube, Linkedin) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(autor.Nome, autor.Facebook, autor.Link, autor.Perfil, autor.Foto, autor.Twitter, autor.Instagram, autor.Youtube, autor.Linkedin)
	} else {
		if autor.Foto == "" {
			autor.Foto = r.FormValue("foto_atual")
		}
		stmt, _ := config.DB.Prepare("UPDATE author SET Nome = ?, Facebook = ?, Link = ?, Perfil = ?, Foto = ?, Twitter = ?, Instagram = ?, Youtube = ?, Linkedin = ? WHERE id = ?")
		stmt.Exec(autor.Nome, autor.Facebook, autor.Link, autor.Perfil, autor.Foto, autor.Twitter, autor.Instagram, autor.Youtube, autor.Linkedin, autor.ID)
	}

	http.Redirect(w, r, "/autores", http.StatusSeeOther)
}

func ExcluirAutor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	stmt, _ := config.DB.Prepare("DELETE FROM author WHERE id = ?")
	stmt.Exec(id)
	http.Redirect(w, r, "/autores", http.StatusSeeOther)
}

