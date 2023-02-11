package link

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

var BLOCK_PREFIX_LIST = []string{"_next", "favicon", "search", "login", "register", "settings", "profile"}

func checkLink(link model.Link) (err error) {
	for _, prefix := range BLOCK_PREFIX_LIST {
		if strings.HasPrefix(link.Keyword, prefix) {
			return e.ErrLinkKeywordNotPermitted
		}
	}

	err = validation.Validate(link.URL,
		validation.Required,
		is.URL,
	)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidURL
	}

	domain_allowed := false
	domain_default := viper.GetString("domain_default")
	domain_alternatives := viper.GetStringSlice("domain_alternatives")

	log.Debug().Caller().Msg(fmt.Sprintf("Default domain: %s", domain_default))

	if link.Domain == domain_default {
		domain_allowed = true
		return
	}

	for _, domain := range domain_alternatives {
		if link.Domain == domain {
			domain_allowed = true
		}
	}

	if !domain_allowed {
		return e.ErrLinkInvalidDomain
	}

	return
}
