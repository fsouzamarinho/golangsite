package models

import (
	"database/sql"
)

type Menu struct {
	Codigo int
	Nome   string
	Link   string
	Ordem  int
	Editorias []Editoria
}

func CarregarMenuComEditorias(db *sql.DB) ([]Menu, error) {
	query := "SELECT Codigo, Nome, Link FROM menu ORDER BY Ordem"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []Menu

	for rows.Next() {
		var m Menu
		if err := rows.Scan(&m.Codigo, &m.Nome, &m.Link); err != nil {
			return nil, err
		}
		editorias, _ := CarregarEditoriasPorMenu(db, m.Codigo)
		m.Editorias = editorias
		menus = append(menus, m)
	}
	return menus, nil
}