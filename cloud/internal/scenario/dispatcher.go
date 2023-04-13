package scenario

import (
	"context"
	"strings"

	"github.com/ksusonic/alice-coffee/cloud/internal/ctx"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/nlg"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/scene"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx *ctx.SceneCtx, req *dialogs.Request, intent string, slots dialogs.Slots, resp *dialogs.Response) *dialogs.Response

type IntentDispatcher struct {
	IntentToFunc map[string]HandlerFunc
	globalCtx    *ctx.GlobalCtx
	logger       *zap.SugaredLogger
}

var IntentMapping = map[string]HandlerFunc{
	IntentMakeCoffee:      scene.MakeCoffee,
	IntentMakeCoffeeTyped: scene.MakeCoffeeTyped,
}

func NewIntentDispatcher(globalCtx *ctx.GlobalCtx, logger *zap.SugaredLogger) *IntentDispatcher {
	return &IntentDispatcher{
		IntentToFunc: IntentMapping,
		globalCtx:    globalCtx,
		logger:       logger,
	}
}

func (d *IntentDispatcher) Handler(c context.Context, k dialogs.Kit) (response *dialogs.Response) {
	req, resp := k.Init()
	d.logger = d.logger.With(zap.String("req-id", req.ReqId()))
	if req.IsNewSession() {
		return resp.Text(nlg.RandomGreeting())
	}
	sceneContext := &ctx.SceneCtx{
		Ctx: c,
		GlobalCtx: &ctx.GlobalCtx{
			Socket: d.globalCtx.Socket,
		},
		Logger: d.logger.Named("scene"),
		ReqId:  req.ReqId(),
	}
	defer func() {
		if r := recover(); r != nil {
			d.logger.Errorf("recovered scene error: %v", r)
			response = resp.TextWithTTS(nlg.ErrorPhrase())
		}
	}()

	result := d.TryResponse(sceneContext, req, resp)
	if result != nil {
		return result
	}
	d.logger.Warn("Could not make relevant response")
	return resp.TextWithTTS(nlg.IrrelevantPhrase())
}

func (d *IntentDispatcher) TryResponse(sceneContext *ctx.SceneCtx, req *dialogs.Request, resp *dialogs.Response) *dialogs.Response {
	if len(req.Request.NLU.Intents) == 0 {
		d.logger.Debug("no intents found")
		return nil
	}

	for intent, slots := range req.Request.NLU.Intents {
		f, ok := d.IntentToFunc[intent]
		if ok {
			d.logger.Debugf("matched intent %s", intent)
			return f(sceneContext, req, intent, slots, resp)
		}
	}

	d.logger.Warn("intents are unknown: ", strings.Join(req.Request.NLU.Intents.Names(), ", "))
	return nil
}
