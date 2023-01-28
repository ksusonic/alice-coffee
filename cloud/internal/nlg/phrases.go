package nlg

import (
	"fmt"
	"math/rand"
)

const (
	Greeting = "Привет! Я кофеаппарат, который может сделать тебе кофе. " +
		"Просто попроси меня приготовить конкретный напиток!"
	WhichCoffee = "Какой кофе вы хотите приготовить? В наличии латте, капучино, какао и другие."
)

func randomOk() string {
	var oks = []string{
		"Окей.",
		"Хорошо.",
		"Поняла.",
	}
	return oks[rand.Intn(len(oks))]
}
func randomMaking() string {
	var making = []string{
		"Готовлю",
		"Делаю",
		"Варю",
	}
	return making[rand.Intn(len(making))]
}

func randomError() string {
	askLater := "Спросите еще раз попозже, пожалуйста."
	var making = []string{
		"Ой-ой, произошла ошибка. " + askLater,
		"Кажется, я немного сломалась. " + askLater,
		"Упс, кажется произошла ошибка. " + askLater,
	}
	return making[rand.Intn(len(making))]
}

func MakingCoffeePhrase(coffeeType string, sugar uint32) (string, string) {
	phrase := fmt.Sprintf("%s %s %s", randomOk(), randomMaking(), coffeeType)
	if sugar > 0 {
		phrase += "c сахаром."
	} else {
		phrase += "."
	}
	return phrase, phrase
}

func ErrorPhrase() (string, string) {
	err := randomError()
	return err, err
}
