package facades

import (
	"github.com/goravel/framework/contracts/hash"
)

func Hash() hash.Hash {
	return App().MakeHash()
}
