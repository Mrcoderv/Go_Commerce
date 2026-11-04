package seller

import "go.uber.org/fx"

var Module = fx.Module(
	"seller",
	fx.Provide(
		NewSellerService,
	),
)
