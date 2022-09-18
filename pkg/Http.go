package pkg

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func timeTrack(start time.Time) time.Duration {
	return time.Since(start)
}

func Ping(url string, timeout time.Duration, responseCode int) bool {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error in the request")
		return false
	}
	client := http.Client{Timeout: timeout * time.Second}

	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Error on execution of request")
		return false
	}
	log.Info().Int("status code", response.StatusCode).Int64("duration", timeTrack(start).Milliseconds()).Msgf("Execution of %s", url)
	if responseCode != response.StatusCode {
		return false
	}
	return true
}
