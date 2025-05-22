package handlers

import (
	"andravirtual/config"
	"andravirtual/models"
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RecaptchaResponse struct {
	Success bool `json:"success"`
}

func VerificarRecaptcha(token string, remoteIP string) bool {
	secret := "6LfDjMEqAAAAAEDvRToFg6Kz8BrwUaSNmOeVhNjz"
	data := "secret=" + secret + "&response=" + token + "&remoteip=" + remoteIP
	resp, err := http.Post("https://www.google.com/recaptcha/api/siteverify",
		"application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var resultado RecaptchaResponse
	err = json.NewDecoder(resp.Body).Decode(&resultado)
	return err == nil && resultado.Success
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	erro := r.URL.Query().Get("erro")
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, struct{ Erro string }{Erro: erro})
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	login := r.FormValue("login")
	senha := r.FormValue("senha")
	token := r.FormValue("g-recaptcha-response")
	ip := r.RemoteAddr

	if !VerificarRecaptcha(token, ip) {
		http.Redirect(w, r, "/login?erro=recaptcha", http.StatusSeeOther)
		return
	}

	var user models.Usuario
	err := config.DB.QueryRow("SELECT id_usuario, login_usuario, nome_usuario, senha_usuario, email_usuario FROM usuarios WHERE login_usuario = ?", login).
		Scan(&user.ID, &user.Login, &user.Nome, &user.Senha, &user.Email)
	if err != nil {
		http.Redirect(w, r, "/login?erro=usuario", http.StatusSeeOther)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(senha))
	if err != nil {
		http.Redirect(w, r, "/login?erro=senha", http.StatusSeeOther)
		return
	}

	cookie := &http.Cookie{
		Name:    "usuario_id",
		Value:   strconv.Itoa(user.ID),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "usuario_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func UsuarioLogado(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("usuario_id")
	if err != nil {
		return 0, false
	}
	id, _ := strconv.Atoi(cookie.Value)
	return id, true
}
