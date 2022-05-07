package link

import (
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
)

func checkLink(link entity.Link) (err error) {
	// Validate domain.
	if err = validation.Validate(link.Domain,
		validation.Required,
		is.Domain,
	); err != nil {
		log.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidDomain
	}

	// Validate keyword.
	if err = validation.Validate(link.Keyword,
		validation.Required,
		validation.Length(6, 15),
		is.Alphanumeric,
	); err != nil {
		log.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidKeyword
	}

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
