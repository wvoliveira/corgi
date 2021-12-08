package link

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func checkLink(log log.Logger, link entity.Link) (err error) {
	// Validate domain.
	if err = validation.Validate(link.Domain,
		validation.Required,
		is.Domain,
	); err != nil {
		log.Warnf("user input a invalid domain: %s", err.Error())
		return e.ErrLinkInvalidDomain
	}

	// Validate keyword.
	if err = validation.Validate(link.Keyword,
		validation.Required,
		validation.Length(6, 15),
		is.Alphanumeric,
	); err != nil {
		log.Warnf("user input a invalid keyword: %s", err.Error())
		return e.ErrLinkInvalidKeyword
	}

	// Validate URL.
	if err = validation.Validate(link.URL,
		validation.Required,
		is.URL,
	); err != nil {
		log.Warnf("user input a invalid URL: %s", err.Error())
		return e.ErrLinkInvalidURL
	}
	return
}
