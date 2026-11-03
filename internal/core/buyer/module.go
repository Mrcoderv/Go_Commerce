package buyer

import (
	"go.uber.org/fx"
)

var Module = fx.Module("buyer", fx.Provide(NewBuyerService))
