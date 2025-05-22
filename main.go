package main

import (
	"andravirtual/config"
	"andravirtual/handlers"
	"andravirtual/middleware"
	"net/http"
)

func main() {
	// Conectar ao banco de dados
	config.Connect()

	// Rota para servir imagens e arquivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Rota de login e logout
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.Login(w, r)
		} else {
			handlers.LoginPage(w, r)
		}
	})
	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc("/", handlers.Home)

	// Rotas protegidas por login
	http.HandleFunc("/admin", handlers.ListarEditorias)

	http.HandleFunc("/noticias/nova", middleware.RequireLogin(handlers.FormNovaNoticia))
	http.HandleFunc("/noticias", middleware.RequireLogin(handlers.SalvarNoticia))

	http.HandleFunc("/editorias", middleware.RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.SalvarEditoria(w, r)
		} else {
			handlers.ListarEditorias(w, r)
		}
	}))
	http.HandleFunc("/editorias/nova", middleware.RequireLogin(handlers.NovaEditoria))
	http.HandleFunc("/editorias/editar", middleware.RequireLogin(handlers.EditarEditoria))
	http.HandleFunc("/editorias/excluir", middleware.RequireLogin(handlers.ExcluirEditoria))

	http.HandleFunc("/menus", middleware.RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.SalvarMenu(w, r)
		} else {
			handlers.ListarMenus(w, r)
		}
	}))
	http.HandleFunc("/menus/novo", middleware.RequireLogin(handlers.NovoMenu))
	http.HandleFunc("/menus/editar", middleware.RequireLogin(handlers.EditarMenu))
	http.HandleFunc("/menus/excluir", middleware.RequireLogin(handlers.ExcluirMenu))

	http.HandleFunc("/autores", middleware.RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.SalvarAutor(w, r)
		} else {
			handlers.ListarAutores(w, r)
		}
	}))
	http.HandleFunc("/autores/novo", middleware.RequireLogin(handlers.NovoAutor))
	http.HandleFunc("/autores/editar", middleware.RequireLogin(handlers.EditarAutor))
	http.HandleFunc("/autores/excluir", middleware.RequireLogin(handlers.ExcluirAutor))

	// Início do servidor
	http.ListenAndServe(":8080", nil)
}
