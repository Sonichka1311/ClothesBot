package utils

import (
	"fmt"
	"strings"
)

func ToEng(text string) string {
	text = strings.ToLower(strings.Split(text, " ")[0])
	switch text {
	case "холодно":
		return "cold"
	case "нормально":
		return "normal"
	case "тепло":
		return "warm"
	case "жарко":
		return "hot"
	case "готово":
		return "done"
	case "верх":
		return "top"
	case "низ":
		return "bottom"
	case "комбинация":
		return "combo"
	case "верхняя":
		return "outer"
	case "обувь":
		return "shoes"
	case "светлое":
		return "light"
	case "тёмное":
		return "dark"
	case "любое":
		return "any"
	case "раздельный":
		return "sep"
	case "совместный":
		return "comb"
	case "грязное":
		return "dirty"
	case "чистое":
		return "clean"
	}
	panic(fmt.Sprintf("unknown word %s", text))
}

func ToRus(text string) string {
	text = strings.ToLower(text)
	switch text {
	case "cold":
		return "холодно"
	case "normal":
		return "нормально"
	case "warm":
		return "тепло"
	case "hot":
		return "жарко"
	case "top":
		return "верх"
	case "bottom":
		return "низ"
	case "combo":
		return "комбинация"
	case "outer":
		return "верхняя одежда"
	case "shoes":
		return "обувь"
	case "dirty":
		return "грязное"
	case "clean":
		return "чистое"
	}
	panic(fmt.Sprintf("unknown word %s", text))
}

func InvariantCondition(cond string) string {
	if cond == "чистое" {
		return "грязное"
	}
	return "чистое"
}