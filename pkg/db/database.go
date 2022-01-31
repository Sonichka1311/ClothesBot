package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"bot/pkg/constants"

	"gorm.io/gorm"
)

type Database struct {
	DB   *sql.DB
	Gorm *gorm.DB
}

//////////////////////////// users /////////////////////////////////////////////

func (d Database) CreateUser(id int) {
	err := d.Gorm.Create(&constants.User{ID: id, State: "main"}).Error
	if err != nil {
		log.Println("Error while add user to database users: ", err.Error())
		return
	}
}

func (d Database) UpdateUser(user *constants.User) {
	err := d.Gorm.Model(&constants.User{ID: user.ID}).Updates(user).Error
	if err != nil {
		log.Println("Error while user in database: ", err.Error())
		return
	}
}

func (d Database) GetUser(id int) *constants.User {
	var user constants.User
	err := d.Gorm.Take(&user, "id=?", id).Error
	if err != nil {
		log.Println("Error while select user from database: ", err.Error())
		return nil
	}

	return &user
}

func (d Database) UpdateState(id int, state string) {
	user := constants.User{
		ID:    id,
		State: state,
	}

	d.UpdateUser(&user)
}

func (d Database) SetTopColor(id int, color string) {
	user := constants.User{
		ID:       id,
		TopColor: color,
	}

	d.UpdateUser(&user)
}

func (d Database) SetBottomColor(id int, color string) {
	user := constants.User{
		ID:          id,
		BottomColor: color,
	}

	d.UpdateUser(&user)
}

func (d Database) SetUserSeason(id int, season string) {
	user := constants.User{
		ID:     id,
		Season: season,
	}

	d.UpdateUser(&user)
}

func (d Database) SetRecent(id int, recent int) {
	user := constants.User{
		ID:         id,
		LastFileID: recent,
	}

	d.UpdateUser(&user)
}

////////////////////////////////// data ////////////////////////////////

func (d Database) AddThing(userID, thingID int) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Purity: "clean",
	}

	err := d.Gorm.Create(&thing).Error
	if err != nil {
		log.Println("Error while add thing to database: ", err.Error())
		return
	}
}

func (d Database) UpdateThing(thing *constants.Thing) {
	err := d.Gorm.Model(&constants.Thing{UserID: thing.UserID, ID: thing.ID}).Updates(thing).Error
	if err != nil {
		log.Println("Error while update thing in database: ", err.Error())
		return
	}
}

func (d Database) GetThing(userID, thingID int) *constants.Thing {
	var user constants.Thing
	err := d.Gorm.Take(&user, "user_id=? and id=?", userID, thingID).Error
	if err != nil {
		log.Println("Error while select thing from database: ", err.Error())
		return nil
	}

	return &user
}

func (d Database) SetPhoto(userID, thingID int, file string) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Photo:  file,
	}

	d.UpdateThing(&thing)
}

func (d Database) SetName(userID, thingID int, text string) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Name:   text,
	}

	d.UpdateThing(&thing)
}

func (d Database) SetType(userID, thingID int, text string) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Type:   text,
	}

	d.UpdateThing(&thing)
}

func (d Database) SetSeason(userID, thingID int, season string) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
	}

	switch season {
	case "cold":
		thing.Cold = true
	case "normal":
		thing.Normal = true
	case "warm":
		thing.Warm = true
	case "hot":
		thing.Hot = true
	}

	d.UpdateThing(&thing)
}

func (d Database) UnsetSeason(userID, thingID int, season string) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
	}

	switch season {
	case "cold":
		thing.Cold = false
	case "normal":
		thing.Normal = false
	case "warm":
		thing.Warm = false
	case "hot":
		thing.Hot = false
	}

	err := d.Gorm.Model(&constants.Thing{UserID: thing.UserID, ID: thing.ID}).Select(season).Updates(thing).Error
	if err != nil {
		log.Println("Error while update thing in database: ", err.Error())
		return
	}
}

func (d Database) SetColor(userID, thingID int, color string) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Color:  color,
	}

	d.UpdateThing(&thing)
}

func (d Database) MakeDirty(userID, thingID int) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Purity: "dirty",
	}

	d.UpdateThing(&thing)
}

func (d Database) MakeClean(userID, thingID int) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
		Purity: "clean",
	}

	d.UpdateThing(&thing)
}

func (d Database) GetThingsByUser(userID int) []*constants.Thing {
	var things []*constants.Thing
	err := d.Gorm.Where("user_id=?", userID).Find(&things).Error
	if err != nil {
		log.Printf("Error while selecting all wardrobe from database: %s\n", err.Error())
		return nil
	}
	return things
}

func (d Database) GetByParams(userID int, color string, types string, season string) []*constants.Thing {
	conditions := []string{
		"user_id = ?",
		"type = ?",
		"purity = ?",
		fmt.Sprintf("%s = ?", season),
	}
	args := []interface{}{userID, types, "clean", true}

	if color != "any" {
		conditions = append(conditions, "color = ?")
		args = append(args, color)
	}

	var things []*constants.Thing
	err := d.Gorm.Where(strings.Join(conditions, " and "), args...).Find(&things).Error
	if err != nil {
		log.Printf("Error while selecting by params from database: %s\n", err.Error())
		return nil
	}

	return things
}

func (d Database) ListDirty(userID int) []*constants.Thing {
	var things []*constants.Thing
	err := d.Gorm.Where("user_id=? and purity=?", userID, "dirty").Find(&things).Error
	if err != nil {
		log.Printf("Error while selecting all dirty wardrobe from database: %s\n", err.Error())
		return nil
	}
	return things
}

func (d Database) ListByType(userID int, types string) []*constants.Thing {
	var things []*constants.Thing
	err := d.Gorm.Where("user_id=? and type=?", userID, types).Find(&things).Error
	if err != nil {
		log.Printf("Error while selecting all dirty wardrobe from database: %s\n", err.Error())
		return nil
	}
	return things
}

func (d Database) DeleteThing(userID, thingID int) {
	thing := constants.Thing{
		UserID: userID,
		ID:     thingID,
	}

	d.Gorm.Delete(&thing)
}
