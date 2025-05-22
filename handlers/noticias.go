package handlers

import (
	"andravirtual/config"
	"andravirtual/models"
	"html/template"
	"net/http"
	"andravirtual/utils"
	"time"
)

func FormNovaNoticia(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/noticias/nova.html"))
	tmpl.Execute(w, nil)
}

func SalvarNoticia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	noticia := models.Noticia{
		Topico:        r.FormValue("topico"),
		Resumo:        r.FormValue("resumo"),
		Classe:        utils.AtoiOrZero(r.FormValue("classe")),
		Classe1:       r.FormValue("classe2"),
		PalavrasChave: r.FormValue("palavras"),
		Autor:         r.FormValue("autor"),
		Link:          r.FormValue("legenda"),
		Noticia:       r.FormValue("noticia"),
		CodColunistas: utils.AtoiOrZero(r.FormValue("colunista")),
		Pasta:         "2025/",
		Data:          time.Now().Format("2006-01-02 15:04:05"),
	}

	stmt, err := config.DB.Prepare(`INSERT INTO noticias 
	(Topico, Resumo, Classe, Classe1, Palavras_chave, Autor, Link, Noticia, Foto, Foto1, Cod_colunitas, Pasta, Data) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, '', '', ?, ?, ?)`)
	if err != nil {
		http.Error(w, "Erro ao preparar statement", 500)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(noticia.Topico, noticia.Resumo, noticia.Classe, noticia.Classe1,
		noticia.PalavrasChave, noticia.Autor, noticia.Link, noticia.Noticia,
		noticia.CodColunistas, noticia.Pasta, noticia.Data)

	if err != nil {
		http.Error(w, "Erro ao salvar notícia", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

