package link

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

func checkLink(link model.Link) (err error) {
	// Validate URL.
	if err = validation.Validate(link.URL,
		validation.Required,
		is.URL,
	); err != nil {
		log.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidURL
	}
	return
}
