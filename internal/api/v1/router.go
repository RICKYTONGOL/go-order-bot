// Package v1 implements API server version 1.0.
package v1

import (
	"fmt"
	"net/http"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
	"github.com/rs/cors"
)

const (
	ApiV1UrlPattern = "/api/v1"
)

func SetupHandlers(mux *http.ServeMux, cfg *config.Config, log logger.Logger, corsHandler *cors.Cors) {
	setupTradeHistoryHandlers(mux, log, corsHandler)
}

func urlPatternStr(subURL string) string {
	return fmt.Sprintf("%s%s", ApiV1UrlPattern, subURL)
}
