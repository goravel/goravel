package facades

import (
	"github.com/goravel/framework/contracts/ai"
)

func AI() ai.AI {
	return App().MakeAI()
}
