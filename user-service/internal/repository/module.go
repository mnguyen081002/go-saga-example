package repository

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewItemRepository,
)
