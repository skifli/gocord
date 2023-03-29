package api

import (
	"time"

	"github.com/valyala/fasthttp"
)

type genericMap = map[string]any

var (
	clientBuildNumber = mustGetLatestBuild()
	systemLocale      = mustGetSystemLocale()
	requestClient     = fasthttp.Client{
		ReadBufferSize:                8192,
		ReadTimeout:                   time.Duration(time.Second * 5),
		WriteTimeout:                  time.Duration(time.Second * 5),
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}
)

const (
	API_VERSION     = "9"
	BROWSER         = "Firefox"
	BROWSER_VERSION = "111.0"
	CAPABILITIES    = 4093
	DEVICE          = "" // Discord's official client sends an empty string.
	OS              = "Windows"
	OS_VERSION      = "10"
	STATUS          = "offline" // https://discord.com/developers/docs/topics/gateway-events#update-presence-status-types
	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/" + BROWSER_VERSION
)
