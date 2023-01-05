package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// GetOrCreateDataFolder create Corgi folder for settings, database, etc.
func GetOrCreateDataFolder() (folder string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	folder = filepath.Join(home, ".corgi")
	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}
	return
}

// PrintRoutes print routes with methods based mux Router object.
func PrintRoutes(rs []*mux.Router) {
	for _, r := range rs {
		_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			uri, err := route.GetPathTemplate()
			if err != nil {
				log.Error().Msg(fmt.Sprintf("with get path template: %s", err.Error()))
				return err
			}

			method, err := route.GetMethods()
			if err != nil {
				if errors.Is(err, mux.ErrMethodMismatch) {
					return err
				}
			}

			if uri != "" && len(method) != 0 {
				log.Debug().Caller().Msg(fmt.Sprintf("%s %s", uri, method))
			}
			return nil
		})
	}
}

// AddContextInfo add any information inside context.
func AddContextInfo(ctx context.Context, key interface{}, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// Contains check if contain a specific item in a list.
func Contains(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}

// SplitURL domain and keyword from shortened URL
func SplitURL(url string) (domain, keyword string) {
	if url != "" {
		splited := strings.Split(url, "/")

		if len(splited) == 2 {
			domain = splited[0]
			keyword = splited[1]
		}
	}
	return
}
