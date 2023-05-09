package nlg

import (
	"fmt"
	"math/rand"
)

const (
	WhichCoffee     = "Какой кофе вы хотите приготовить? В наличии латте, капучино, какао и другие."
	WhatCanYouDo    = "Я могу приготовить вам кофе: латте, капучино, горячий шоколад. Данный аппарат экспериментальный - при любых проблемах или пожеланиях пишите на рассылку alice-coffee."
	feedbackAddress = "При повторении ошибки напишите на рассылку alice-coffee@."
)

func randomPhrase(phrase ...string) string {
	return phrase[rand.Intn(len(phrase))]
}

func randomPhraseWithPrefix(prefix string, phrase ...string) string {
	return prefix + " " + phrase[rand.Intn(len(phrase))]
}

func randomPhraseWithSuffix(suffix string, phrase ...string) string {
	return phrase[rand.Intn(len(phrase))] + " " + suffix
}

func RandomGreeting() string {
	return randomPhraseWithPrefix(
		"Привет! Я экспериментальный (умный) кофеаппарат.",
		"Я умею готовить капучино, латте, чай, а также другие напитки.",
		"Просто попросите меня, и я с удовольствием приготовлю вам любимый напиток.",
	)
}

func randomOk() string {
	return randomPhrase(
		"Окей.",
		"Хорошо.",
		"Поняла.",
	)
}
func randomMaking() string {
	return randomPhrase(
		"Готовлю",
		"Делаю",
		"Варю",
	)
}

func MakingCoffeePhrase(coffeeType string, sugar uint) string {
	phrase := fmt.Sprintf("%s %s %s", randomOk(), randomMaking(), coffeeType)
	if sugar > 0 {
		phrase += "c сахаром."
	} else {
		phrase += "."
	}
	return phrase
}

func ErrorPhrase() string {
	return randomPhraseWithSuffix(
		feedbackAddress,
		"Ой-ой, произошла ошибка. ",
		"Кажется, я немного сломалась. ",
		"Упс, произошла ошибка. ",
	)
}

func NoConnectionPhrase() string {
	return randomPhraseWithSuffix(
		feedbackAddress,
		"Потеряла соединение с автоматом. ",
		"Кажется, я немного сломалась. ",
		"Упс, произошла ошибка. ",
	)
}

func IrrelevantPhrase() string {
	return randomPhrase(
		"Не смогла понять вас, извините.",
		"Упс, кажется, я еще так не умею.",
		"Простите, не совсем вас поняла.",
	)
}
