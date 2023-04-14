package scene

import (
	"time"

	"github.com/ksusonic/alice-coffee/cloud/internal/ctx"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/models"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/nlg"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
	"go.uber.org/zap"
)

func MakeCoffee(
	_ *ctx.SceneCtx,
	_ *dialogs.Request,
	_ string,
	_ dialogs.Slots,
	resp *dialogs.Response) *dialogs.Response {
	return resp.Text(nlg.WhichCoffee)
}

func MakeCoffeeTyped(
	ctx *ctx.SceneCtx,
	_ *dialogs.Request,
	_ string,
	slots dialogs.Slots,
	resp *dialogs.Response) *dialogs.Response {

	coffeeTypeValue := slots.Slots["coffee_type"].Value
	coffeeType := models.ParseCoffee(coffeeTypeValue)
	if coffeeType == nil {
		ctx.Logger.Errorf("no coffee type found for %s", coffeeTypeValue)
		coffeeType = &models.UnknownCoffee
	}

	if ctx.GlobalCtx.Socket.CheckActive() == false {
		ctx.Logger.Error("no connection to vending")
		return resp.Text(nlg.NoConnectionPhrase())
	}

	var sugarAmount uint = 0 // TODO

	request := models.MakeCoffeeRequest{
		ID:    time.Now().String(), // todo request-id
		Type:  *coffeeType,
		Sugar: sugarAmount,
	}
	err := ctx.GlobalCtx.Socket.SendJSON(request)
	if err != nil {
		ctx.Logger.Error("could not send make_coffee request", zap.Error(err), zap.String("id", request.ID))
		return resp.Text(nlg.ErrorPhrase())
	}
	ctx.Logger.Info("sent make_coffee request", zap.Any("request", request))

	return resp.Text(nlg.MakingCoffeePhrase(coffeeType.HumanReadable, sugarAmount))
}
