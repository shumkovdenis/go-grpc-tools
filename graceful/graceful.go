package graceful

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

type Graceful interface {
	Start() error
	GracefulStop() error
}

func Run(service Graceful) {
	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// run blocking call in a separate goroutine, report errors via channel
	go func() {
		log.Info().Msg("starting server")
		if err := service.Start(); err != nil {
			errChan <- err
		}
	}()

	// terminate your environment gracefully before leaving main function
	defer func() {
		log.Info().Msg("stopping server")
		service.GracefulStop()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Fatal().Err(err).Msg("server failed")
	case <-stopChan:
	}
}
