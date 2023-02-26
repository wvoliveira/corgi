package link

import (
	"sort"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

// TODO: put theses keywords in database, so, we can update in real time.
// Ref: https://www.mediavine.com/keyword-anti-targeting/
var blockedKeywords = []string{"crash", "attack", "terrorist", "suicide", "nazi", "killed", "porn", "explosion",
	"rape", "death", "isis", "shooting", "bomb", "dead", "murder", "terror", "kill", "sex", "massacre", "gun"}

func checkLink(domain, keyword, url string) (err error) {
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
