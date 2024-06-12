package post

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewPostService,
		fx.Annotate(
			NewPostRepository,
			fx.As(new(PostRepository)),
		),
	),
)
