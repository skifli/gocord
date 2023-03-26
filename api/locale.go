package api

import (
	"github.com/jeandeaual/go-locale"
)

func mustGetSystemLocale() string {
	userLocales, err := locale.GetLocales()

	if err != nil {
		panic(err)
	}

	return userLocales[0]
}
