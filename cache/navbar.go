package cache

import (
	"bytes"
	"html/template"
	"sync"
	"andravirtual/models"
	"database/sql"
	"fmt"
)

var (
	navbarHTML string
	navbarOnce sync.Once
)

func GetNavbarHTML(db *sql.DB, tmpl *template.Template) (string, error) {
	var err error
	navbarOnce.Do(func() {
		menus, e := models.CarregarMenuComEditorias(db)
		if e != nil {
			err = e
			return
		}
		var buf bytes.Buffer
		e = tmpl.ExecuteTemplate(&buf, "navbar", map[string]interface{}{"Menus": menus})
		if e != nil {
			err = e
			return
		}
		fmt.Println("MENUS:", menus)
		navbarHTML = buf.String()
	})
	return navbarHTML, err
}