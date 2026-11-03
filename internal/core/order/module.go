package order

import "go.uber.org/fx"

var Module = fx.Module("order", fx.Provide(NewOrderService))
