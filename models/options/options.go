package options

import (
	"context"
	"net/http"
	"strconv"
)

type Options struct {
	Limit  int
	Offset int
}

const defaultLimit = 25
const defaultOffset = 0

// NewOptionsFromRequest returns object with all required options
func NewOptionsFromRequest(r *http.Request) Options {
	var options Options

	if value, err := strconv.Atoi(r.FormValue("limit")); err == nil {
		options.Limit = value
	} else {
		options.Limit = defaultLimit
	}

	if value, err := strconv.Atoi(r.FormValue("offset")); err == nil {
		options.Offset = value
	} else {
		options.Offset = defaultOffset
	}

	return options
}

// NewOptionsFromContext return options from context
func NewOptionsFromContext(ctx context.Context) *Options {
	options, ok := ctx.Value("options").(Options)

	if !ok {
		return &Options{}
	}

	return &options
}
