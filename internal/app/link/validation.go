package link

import (
	"sort"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

func checkLink(domain, keyword, url string) (err error) {
	blockedKeywords := constants.BLOCKED_KEYWORDS

	sort.Strings(blockedKeywords)
	index := sort.SearchStrings(blockedKeywords, keyword)

	if index < len(blockedKeywords) && blockedKeywords[index] == keyword {
		return e.ErrLinkKeywordNotPermitted
	}

	err = validation.Validate(url,
		validation.Required,
		is.URL,
	)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return e.ErrLinkInvalidURL
	}

	domainDefault := viper.GetString("domain_default")

	if domain == domainDefault {
		return nil
	}

	return e.ErrLinkInvalidDomain
}
