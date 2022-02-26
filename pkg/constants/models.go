package constants

import (
	"fmt"
	"strings"

	"bot/pkg/utils"
)

type Thing struct {
	UserID int    `db:"user_id" gorm:"primaryKey"`
	ID     int    `db:"id" gorm:"primaryKey"`
	Photo  string `db:"photo"`

	Name  string `db:"name"`
	Type  string `db:"type"`
	Color string `db:"color"`

	//Single    bool   `db:"single"`
	//Combo     bool   `db:"combo"`
	//ComboType string `db:"combo_type"`

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
	var curCond, toCond string
	if t.Purity == "dirty" {
		curCond = Dirty
		toCond = Clean
	} else {
		curCond = Clean
		toCond = Dirty
	}
	curCond = strings.ToLower(curCond)
	toCond = strings.ToLower(toCond)
	return fmt.Sprintf(
		"*%s* (%s)\n"+
			"Показать полную информацию: /thing\\_%d\n"+
			"Отметить как *%s*: /%s\\_%d\n",
		t.Name, curCond, t.ID, toCond, utils.ToEng(toCond), t.ID,
	)
}

func (t Thing) Caption() string {
	seasons := make([]string, 0)
	if t.Cold {
		seasons = append(seasons, "холодно")
	}
	if t.Normal {
		seasons = append(seasons, "нормально")
	}
	if t.Warm {
		seasons = append(seasons, "тепло")
	}
	if t.Hot {
		seasons = append(seasons, "жарко")
	}

	curCond := utils.ToRus(t.Purity)
	toCond := utils.InvariantCondition(curCond)

	return fmt.Sprintf(
		"*%s*\n"+
			"*Тип:* %s\n"+
			"*Сезон:* %v\n"+
			"*Цвет:* %s\n"+
			"*Состояние:* %s\n"+
			"Отметить как *%s*: /%s\\_%d\n"+
			"Удалить: /delete\\_%d\n",
		t.Name, t.Type, strings.Join(seasons, ", "),
		t.Color, curCond, toCond, utils.ToEng(toCond), t.ID, t.ID,
	)
}

func (t Thing) ShortCaption() string {
	return fmt.Sprintf(
		"*%s*: %s\nShow more: /thing\\_%d\n",
		strings.ToUpper(t.Type[0:1])+t.Type[1:], t.Name, t.ID,
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

type Look struct {
	BaseTop       *Thing
	AdditionalTop *Thing
	Bottom        *Thing
	Combination   *Thing
	Shoes         *Thing
	Outerwear     *Thing
	Caption       string
}
