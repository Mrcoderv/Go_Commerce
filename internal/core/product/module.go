package product

import "go.uber.org/fx"

var Module = fx.Module("product", fx.Provide(NewProductService))
