package states

import (
	"image"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
)

func GetTextIndex(thingType string, numberElements int) int {
	idx := 0
	switch thingType {
	case "bottom":
		idx = 1
	case "shoes":
		idx = 2
		if numberElements == 3 {
			idx = 1
		}
	case "outer":
		idx = 3
		if numberElements == 3 {
			idx = 2
		}
	}

	return idx
}

func GetThingIdByCaption(text string) int {
	stringsInText := strings.Split(text, "\n")

	if len(stringsInText) > 1 {
		num, err := strconv.Atoi(strings.Split(stringsInText[1], "_")[1])
		if err != nil {
			log.Println("Can't get id by caption: ", err.Error())
			return -1
		}

		return num
	}

	return -1
}

func ReGenerateThing(oldThingId int, thing **constants.Thing, things []*constants.Thing) {
	for oldThingId == (*thing).ID {
		log.Println("The same thing occurred.")
		*thing = things[rand.Intn(len(things))]
	}
}

func GetKeyboardWithoutButtonWithIdx(keyboard [][]tb.InlineButton, idx int) *tb.ReplyMarkup {
	return &tb.ReplyMarkup{
		InlineKeyboard: append(keyboard[:idx], keyboard[idx+1:]...),
	}
}

func GetButtonIndex(textIdx int, texts []string) int {
	buttonIdx := textIdx

	for idx := 0; idx < textIdx; idx++ {
		if len(strings.Split(texts[idx], "\n")) > 3 {
			buttonIdx--
		}
	}

	return buttonIdx
}

func ReplaceThingText(texts []string, thing *constants.Thing, textIdx int) string {
	return strings.Join(
		append(
			append(texts[:textIdx], thing.ShortCaption()),
			texts[textIdx+1:]...,
		),
		"\n",
	)
}

func GenerateSomething(things []*constants.Thing) *constants.Thing {
	if len(things) == 0 {
		return nil
	}
	return things[rand.Intn(len(things))]
}

func SplitBigMsg(text string) []string {
	texts := make([]string, 0)
	for len(text) > 0 {
		if len(text) < constants.MaxLength {
			texts = append(texts, text)
			break
		}

		ind := strings.LastIndex(text[:constants.MaxLength], "\n\n")
		texts = append(texts, text[:ind])
		text = text[ind+2:]
	}
	return texts
}

func ResizeImage(img image.Image, size int) *image.NRGBA {
	row := img.Bounds().Dy() / 4
	eps := img.Bounds().Dy() / 100

	first := 0
	last := img.Bounds().Dx()
	for y := row - eps; y < row + eps; y++ {
		firstOK := false
		lastOK := false
		for x := 0; x < img.Bounds().Dx(); x++ {
			a, b, c, d := img.At(x, row).RGBA()
			if a+b+c+d != 0 {
				if firstOK {
					last = x
					lastOK = true
				}

				if !firstOK {
					if first < x {
						first = x
					}
					firstOK = true
				}
			} else if lastOK {
				break
			}
		}
	}

	length := last - first

	return imaging.Resize(img, img.Bounds().Dx() * size / length, 0, imaging.NearestNeighbor)
}

func CropImage(img image.Image) *image.NRGBA {
	minx := 0
	miny := 0
	maxx := img.Bounds().Dx()
	maxy := img.Bounds().Dy()
	for x := 0; x < img.Bounds().Dx(); x++ {
		empty := true
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r + g + b + a != 0 {
				empty = false
			}
		}
		if empty {
			minx++
		} else {
			break
		}
	}
	for y := 0; y < img.Bounds().Dy(); y++ {
		empty := true
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r + g + b + a != 0 {
				empty = false
			}
		}
		if empty {
			miny++
		} else {
			break
		}
	}
	for x := img.Bounds().Dx() - 1; x >= 0; x-- {
		empty := true
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r + g + b + a != 0 {
				empty = false
			}
		}
		if empty {
			maxx--
		} else {
			break
		}
	}
	for y := img.Bounds().Dy() - 1; y >= 0; y-- {
		empty := true
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r + g + b + a != 0 {
				empty = false
			}
		}
		if empty {
			maxy--
		} else {
			break
		}
	}

	return imaging.Crop(img, image.Rect(minx, miny, maxx, maxy))
}
