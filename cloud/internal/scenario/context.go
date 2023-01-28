package scenario

import (
	"github.com/ksusonic/alice-coffee/cloud/internal/queue"
)

type Context struct {
	MessageQueue *queue.MessageQueue
}
