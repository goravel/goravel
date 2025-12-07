package facades

import (
	"github.com/goravel/framework/contracts/queue"
)

func Queue() queue.Queue {
	return App().MakeQueue()
}
