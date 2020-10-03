package db

import (
	"bot/pkg/constants"
	"database/sql"
	"log"
)

type Database struct {
	DB *sql.DB
}

func (d Database) AddUser(id int) {
	_, dbErr := d.DB.Exec("INSERT INTO users (id) VALUES (?)", id)
	if dbErr != nil {
		log.Println("Error while add user to database users: ", dbErr.Error())
	}
}

func (d Database) GetState(id int) string {
	row := d.DB.QueryRow("SELECT state FROM users WHERE id = ? LIMIT 1", id)
	var state string
	dbErr := row.Scan(&state)
	if dbErr != nil {
		log.Println("Error while select state from database: ", dbErr.Error())
	}
	return state
}

func (d Database) UpdateState(id int, state string) {
	_, dbErr := d.DB.Exec("UPDATE users SET state = ? WHERE id = ?", state, id)
	if dbErr != nil {
		log.Println("Error while update state in database: ", dbErr.Error())
	}
}

func (d Database) GetRecent(id int) int {
	row := d.DB.QueryRow("SELECT recent FROM users WHERE id = ? LIMIT 1", id)
	var recent int
	dbErr := row.Scan(&recent)
	if dbErr != nil {
		log.Println("Error while select recent from database: ", dbErr.Error())
	}
	return recent
}

func (d Database) AddThing(id int, recent int) {
	_, dbErr := d.DB.Exec("INSERT INTO data (user, id) VALUES (?, ?)", id, recent)
	if dbErr != nil {
		log.Println("Error while add user to database data: ", dbErr.Error())
	}
}

func (d Database) SetPhoto(id int, recent int, file string) {
	_, dbErr := d.DB.Exec("UPDATE data SET photo = ? WHERE user = ? AND id = ?", file, id, recent)
	if dbErr != nil {
		log.Println("Error while set photo in database: ", dbErr.Error())
	}
}

func (d Database) SetName(id int, recent int, text string) {
	_, dbErr := d.DB.Exec("UPDATE data SET name = ? WHERE user = ? AND id = ?", text, id, recent)
	if dbErr != nil {
		log.Println("Error while set name in database: ", dbErr.Error())
	}
}

func (d Database) SetType(id int, recent int, s string) {
	_, dbErr := d.DB.Exec("UPDATE data SET type = ? WHERE user = ? AND id = ?", s, id, recent)
	if dbErr != nil {
		log.Println("Error while set type in database: ", dbErr.Error())
	}
}

func (d Database) SetCategory(id int, recent int, s string) {
	_, dbErr := d.DB.Exec("UPDATE data SET category = ? WHERE user = ? AND id = ?", s, id, recent)
	if dbErr != nil {
		log.Println("Error while set category in database: ", dbErr.Error())
	}
}

func (d Database) SetSeason(id int, recent int, s string) {
	_, dbErr := d.DB.Exec("UPDATE data SET "+s+" = ? WHERE user = ? AND id = ?", true, id, recent)
	if dbErr != nil {
		log.Println("Error while set season in database: ", dbErr.Error())
	}
}

func (d Database) UnsetSeason(id int, recent int, s string) {
	_, dbErr := d.DB.Exec("UPDATE data SET "+s+" = ? WHERE user = ? AND id = ?", false, id, recent)
	if dbErr != nil {
		log.Println("Error while set season in database: ", dbErr.Error())
	}
}

func (d Database) SetColor(id int, recent int, s string) {
	_, dbErr := d.DB.Exec("UPDATE data SET color = ? WHERE user = ? AND id = ?", s, id, recent)
	if dbErr != nil {
		log.Println("Error while set color in database: ", dbErr.Error())
	}
}

func (d Database) GetName(id int, recent int) string {
	row := d.DB.QueryRow("SELECT name FROM data WHERE user = ? AND id = ? LIMIT 1", id, recent)
	var name string
	dbErr := row.Scan(&name)
	if dbErr != nil {
		log.Println("Error while select name from database: ", dbErr.Error())
	}
	return name
}

func (d Database) SetRecent(id int, i int) {
	_, dbErr := d.DB.Exec("UPDATE users SET recent = ? WHERE id = ?", i, id)
	if dbErr != nil {
		log.Println("Error while update recent in database: ", dbErr.Error())
	}
}

func (d Database) GetAll(id int) *sql.Rows {
	rows, queryError := d.DB.Query("SELECT id, name, purity FROM data WHERE user = ?", id)
	if queryError != nil {
		log.Printf("Error while selecting all wardrobe from database: %s\n", queryError.Error())
	}
	return rows
}

func (d Database) GetByParams(id int, color string, types string, season string) *sql.Rows {
	var rows *sql.Rows
	var queryError error
	if color != "any" {
		rows, queryError = d.DB.Query(
			"SELECT id, name, photo FROM data WHERE user = ? AND type = ? AND color = ? AND purity = ? AND "+season+" = ?",
			id, types, color, "clean", true)
	} else {
		rows, queryError = d.DB.Query(
			"SELECT id, name, photo FROM data WHERE user = ? AND type = ? AND purity = ? AND "+season+" = ?",
			id, types, "clean", true)
	}
	if queryError != nil {
		log.Printf("Error while selecting by params from database: %s\n", queryError.Error())
	}
	return rows
}

func (d Database) GetThing(user int, id int) *constants.Thing {
	row := d.DB.QueryRow("SELECT * FROM data WHERE user = ? AND id = ? LIMIT 1", user, id)
	var thing constants.Thing
	dbErr := row.Scan(&user, &thing.Id, &thing.Photo,
		&thing.Name, &thing.Type, &thing.Cold,
		&thing.Normal, &thing.Warm, &thing.Hot,
		&thing.Color, &thing.Purity)
	if dbErr != nil {
		log.Println("Error while select thing from database: ", dbErr.Error())
	}
	return &thing
}

func (d Database) SetTopColor(id int, s string) {
	_, dbErr := d.DB.Exec("UPDATE users SET topcolor = ? WHERE id = ?", s, id)
	if dbErr != nil {
		log.Println("Error while update top color in database: ", dbErr.Error())
	}
}

func (d Database) SetBottomColor(id int, s string) {
	_, dbErr := d.DB.Exec("UPDATE users SET bottomcolor = ? WHERE id = ?", s, id)
	if dbErr != nil {
		log.Println("Error while update bottom color in database: ", dbErr.Error())
	}
}

func (d Database) GetTopColor(id int) string {
	row := d.DB.QueryRow("SELECT topcolor FROM users WHERE id = ? LIMIT 1", id)
	var topColor string
	dbErr := row.Scan(&topColor)
	if dbErr != nil {
		log.Println("Error while select top color from database: ", dbErr.Error())
	}
	return topColor
}

func (d Database) GetBottomColor(id int) string {
	row := d.DB.QueryRow("SELECT bottomcolor FROM users WHERE id = ? LIMIT 1", id)
	var bottomColor string
	dbErr := row.Scan(&bottomColor)
	if dbErr != nil {
		log.Println("Error while select bottom color from database: ", dbErr.Error())
	}
	return bottomColor
}

func (d Database) MakeDirty(user int, id int) {
	_, dbErr := d.DB.Exec("UPDATE data SET purity = ? WHERE user = ? AND id = ?", "dirty", user, id)
	if dbErr != nil {
		log.Println("Error while update purity in database: ", dbErr.Error())
	}
}

func (d Database) MakeClean(user int, id int) {
	_, dbErr := d.DB.Exec("UPDATE data SET purity = ? WHERE user = ? AND id = ?", "clean", user, id)
	if dbErr != nil {
		log.Println("Error while update purity in database: ", dbErr.Error())
	}
}

func (d Database) GetDirty(id int) *sql.Rows {
	rows, queryError := d.DB.Query("SELECT id, name FROM data WHERE user = ? AND purity = ?", id, "dirty")
	if queryError != nil {
		log.Printf("Error while selecting all dirty wardrobe from database: %s\n", queryError.Error())
	}
	return rows
}

func (d Database) GetByType(id int, types string) *sql.Rows {
	rows, queryError := d.DB.Query("SELECT id, name, purity, photo FROM data WHERE user = ? AND type = ?", id, types)
	if queryError != nil {
		log.Printf("Error while selecting by type from database: %s\n", queryError.Error())
	}
	return rows
}

func (d Database) DeleteThing(user int, id int) {
	_, dbErr := d.DB.Exec("DELETE FROM data WHERE user = ? AND id = ?", user, id)
	if dbErr != nil {
		log.Println("Error while delete from database: ", dbErr.Error())
	}
}

func (d Database) GetSeason(id int) string {
	row := d.DB.QueryRow("SELECT season FROM users WHERE id = ? LIMIT 1", id)
	var season string
	dbErr := row.Scan(&season)
	if dbErr != nil {
		log.Println("Error while select season from database: ", dbErr.Error())
	}
	return season
}

func (d Database) SetUserSeason(id int, season string) {
	_, dbErr := d.DB.Exec("UPDATE users SET season = ? WHERE id = ?", season, id)
	if dbErr != nil {
		log.Println("Error while season in database: ", dbErr.Error())
	}
}
