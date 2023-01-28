package scenario

import (
	"github.com/golang/protobuf/proto"
	"github.com/ksusonic/alice-coffee/cloud/internal/nlg"
	"log"

	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"

	"github.com/ksusonic/alice-coffee/protos/request"
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
	req *dialogs.Request,
	_ string,
	slots dialogs.Slots,
	resp *dialogs.Response) *dialogs.Response {

	coffeeType := slots.Slots[IntentMakeCoffeeTypedCoffeeTypeSlot].Value

	// TODO state check

	var sugarAmount uint32 = 0 // TODO

	request := &requestpb.Request{
		ID: int32(req.MessageID()),
		Action: &requestpb.Request_MakeCoffee{MakeCoffee: &requestpb.MakeCoffeeAction{
			CoffeeType:  coffeeTypeToEnum(coffeeType),
			SugarAmount: sugarAmount,
		}},
	}
	data, err := proto.Marshal(request)
	if err != nil {
		log.Println(err)
		return resp.TextWithTTS(nlg.ErrorPhrase())
	}

	messageId, err := ctx.MessageQueue.SendMessage(data)
	if err != nil {
		log.Println(err)
		return resp.TextWithTTS(nlg.ErrorPhrase())
	}

	log.Printf("sent to queue message_id=%s, body=%s\n", *messageId, request)

	return resp.TextWithTTS(nlg.MakingCoffeePhrase(coffeeTypeToHumanReadable(coffeeType), sugarAmount))
}
