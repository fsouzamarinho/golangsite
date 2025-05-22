package models

import (
	"database/sql"
)

type Editoria struct {
	Codigo  int
	Nome    string
	CodMenu int
}

func CarregarEditoriasPorMenu(db *sql.DB, codMenu int) ([]Editoria, error) {
	rows, err := db.Query("SELECT Nome FROM editorias WHERE Cod_Menu = ? ORDER BY Nome", codMenu)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var editorias []Editoria
	for rows.Next() {
		var e Editoria
		rows.Scan(&e.Nome)
		editorias = append(editorias, e)
	}
	return editorias, nil
}