package app

import "context"

var (
	Ctx, Cancel = context.WithCancel(context.Background())
)
