package link

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

func checkLink(ctx context.Context, link entity.Link) (err error) {
	l := log.Ctx(ctx)

	err = validation.Validate(link.URL, validation.Required, is.URL)
	if err != nil {
		l.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidURL
	}

	return
}
