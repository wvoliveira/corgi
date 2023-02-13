package link

import (
	"sort"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

func checkLink(link model.Link) (err error) {
	blockedKeywords := constants.BLOCKED_KEYWORDS

	sort.Strings(blockedKeywords)
	index := sort.SearchStrings(blockedKeywords, link.Keyword)

	if index < len(blockedKeywords) && blockedKeywords[index] == link.Keyword {
		return e.ErrLinkKeywordNotPermitted
	}

	err = validation.Validate(link.URL,
		validation.Required,
		is.URL,
	)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidURL
	}

	domainDefault := viper.GetString("domain_default")

	if link.Domain == domainDefault {
		return nil
	}

	return e.ErrLinkInvalidDomain
}
