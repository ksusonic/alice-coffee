package scenario

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/models"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/nlg"

	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
)

func MakeCoffee(
	_ *Context,
	_ *dialogs.Request,
	_ string,
	_ dialogs.Slots,
	resp *dialogs.Response) *dialogs.Response {
	return resp.TextWithTTS(nlg.WhichCoffee, nlg.WhichCoffee)
}

func MakeCoffeeTyped(
	ctx *Context,
	_ *dialogs.Request,
	_ string,
	slots dialogs.Slots,
	resp *dialogs.Response) *dialogs.Response {

	coffeeTypeValue := slots.Slots[IntentMakeCoffeeTypedCoffeeTypeSlot].Value
	coffeeType := models.ParseCoffee(coffeeTypeValue)
	if coffeeType == nil {
		log.Println("no coffee type found for ", coffeeTypeValue)
		coffeeType = &models.UnknownCoffee
	}

	// TODO state check

	var sugarAmount uint = 0 // TODO

	request := models.MakeCoffeeRequest{
		ID:    time.Now().String(), // todo request-id
		Type:  *coffeeType,
		Sugar: sugarAmount,
	}

	data, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
		return resp.TextWithTTS(nlg.ErrorPhrase())
	}

	err = ctx.Socket.SendMessage(data)
	if err != nil {
		log.Println(err)
		return resp.TextWithTTS(nlg.ErrorPhrase())
	}

	log.Printf("sent: %s\n", string(data))

	return resp.TextWithTTS(nlg.MakingCoffeePhrase(coffeeType.HumanReadable, sugarAmount))
}
