package nlg

import (
	"fmt"
	"math/rand"
)

const (
	greeting1 = "Привет! Я кофеаппарат, который может сделать тебе кофе. " +
		"Просто попроси меня приготовить конкретный напиток!"
	greeting2 = "Привет! Я готов помочь насладиться идеальным кофе-экспериментом в любое время дня. Умею готовить капучино с бархатистой пенкой, латте с воздушным молочным шлейфом и горячий шоколад с богатым вкусом." +
		"Просто скажите мне, что вы хотите, и я с удовольствием приготовлю вам ваш любимый напиток."
	WhichCoffee = "Какой кофе вы хотите приготовить? В наличии латте, капучино, какао и другие."

	feedbackAddress = "При повторении ошибки напишите на рассылку alice-coffee@."
)

func RandomGreeting() string {
	var greetings = []string{
		greeting1,
		greeting2,
	}
	return greetings[rand.Intn(len(greetings))]
}

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

func MakingCoffeePhrase(coffeeType string, sugar uint) (string, string) {
	phrase := fmt.Sprintf("%s %s %s", randomOk(), randomMaking(), coffeeType)
	if sugar > 0 {
		phrase += "c сахаром."
	} else {
		phrase += "."
	}
	return phrase, phrase
}

func ErrorPhrase() (string, string) {
	var making = []string{
		"Ой-ой, произошла ошибка. ",
		"Кажется, я немного сломалась. ",
		"Упс, кажется произошла ошибка. ",
	}
	err := making[rand.Intn(len(making))] + feedbackAddress
	return err, err
}

func NoConnectionPhrase() (string, string) {
	var making = []string{
		"Не могу найти кнопки от вендинга. ",
		"Кажется, я немного сломалась. ",
		"Упс, кажется произошла ошибка. ",
	}
	err := making[rand.Intn(len(making))] + feedbackAddress
	return err, err
}

func IrrelevantPhrase() (string, string) {
	var irrels = []string{
		"Не смогла понять вас, извините.",
		"Упс, кажется, я еще так не умею.",
		"Простите, не совсем вас поняла.",
	}
	chosen := irrels[rand.Intn(len(irrels))]
	return chosen, chosen
}
