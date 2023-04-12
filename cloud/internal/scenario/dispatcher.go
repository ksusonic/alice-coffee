package scenario

import (
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/nlg"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
)

type HandlerFunc func(ctx *Context, req *dialogs.Request, intent string, slots dialogs.Slots, resp *dialogs.Response) *dialogs.Response

type IntentDispatcher struct {
	IntentToFunc map[string]HandlerFunc
	Context      *Context
}

var IntentMapping = map[string]HandlerFunc{
	IntentMakeCoffee:      MakeCoffee,
	IntentMakeCoffeeTyped: MakeCoffeeTyped,
}

func NewIntentDispatcher(ctx *Context) *IntentDispatcher {
	return &IntentDispatcher{
		IntentToFunc: IntentMapping,
		Context:      ctx,
	}
}

func (d *IntentDispatcher) TryResponse(req *dialogs.Request, resp *dialogs.Response) *dialogs.Response {
	if len(req.Request.NLU.Intents) == 0 {
		return nil
	}

	for intent, slots := range req.Request.NLU.Intents {
		f, ok := d.IntentToFunc[intent]
		if ok {
			return f(d.Context, req, intent, slots, resp)
		}
	}
	return resp
}

func (d *IntentDispatcher) IrrelevantResponse(_ *dialogs.Request, resp *dialogs.Response) *dialogs.Response {
	return resp.TextWithTTS(nlg.IrrelevantPhrase())
}
