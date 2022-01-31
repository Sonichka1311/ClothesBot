package constants

import "fmt"

type Thing struct {
	UserID int    `db:"user_id" gorm:"primaryKey"`
	ID     int    `db:"id" gorm:"primaryKey"`
	Photo  string `db:"photo"`

	Name  string `db:"name"`
	Type  string `db:"type"`
	Color string `db:"color"`

	Cold   bool `db:"cold"`
	Normal bool `db:"normal"`
	Warm   bool `db:"warm"`
	Hot    bool `db:"hot"`

	Purity string `db:"purity"`
}

func (t Thing) TableName() string {
	return "data"
}

func (t Thing) ListCaption() string {
	var toCond string
	if t.Purity == "dirty" {
		toCond = Clean
	} else {
		toCond = Dirty
	}
	return fmt.Sprintf(
		"*%s* (%s)\n"+
			"See more: /thing\\_%d\n"+
			"Mark thing *%s*: /%s\\_%d\n",
		t.Name, t.Purity, t.ID, toCond, toCond, t.ID,
	)
}

type User struct {
	ID    int    `db:"id" gorm:"primaryKey"`
	State string `db:"state"`

	LastFileID int `db:"last_file_id"`

	TopColor    string `db:"top_color"`
	BottomColor string `db:"bottom_color"`
	Season      string `db:"season"`
}
