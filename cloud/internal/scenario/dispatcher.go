package scenario

import (
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
)

type HandlerFunc func(ctx *Context, req *dialogs.Request, intent string, slots dialogs.Slots, resp *dialogs.Response) *dialogs.Response

type IntentDispatcher struct {
	IntentToFunc map[string]HandlerFunc
	Context      *Context
}

func NewIntentDispatcher(
	ctx *Context,
) *IntentDispatcher {
	return &IntentDispatcher{
		IntentToFunc: map[string]HandlerFunc{
			IntentMakeCoffee:      MakeCoffee,
			IntentMakeCoffeeTyped: MakeCoffeeTyped,
		},
		Context: ctx,
	}
}

func (d *IntentDispatcher) TryResponse(req *dialogs.Request, resp *dialogs.Response) *dialogs.Response {
	if len(req.Request.NLU.Intents) == 0 {
		return nil
	}
	intent, slots := firstIntent(req.Request.NLU.Intents)

	f, ok := d.IntentToFunc[intent]
	if !ok {
		return nil
	}

	return f(d.Context, req, intent, slots, resp)
}

func firstIntent(i dialogs.Intents) (string, dialogs.Slots) {
	for intent, slots := range i {
		return intent, slots
	}
	return "", dialogs.Slots{}
}
