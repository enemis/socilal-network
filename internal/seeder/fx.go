package seeder

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewSeeder,
	),
	fx.Invoke(
		RunImport,
	),
)
