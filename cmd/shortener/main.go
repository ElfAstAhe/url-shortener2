package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/ElfAstAhe/url-shortener2/internal/app/bootstrap"
)

func main() {
	// app instance
	app := bootstrap.NewApp()
	//	defer app.Close()
	logger := app.Log.GetLogger("main")
	//	defer _utl.CloseOnly(logger.(io.Closer))

	// app initialization
	logger.Info("app init")
	if err := app.Init(); err != nil {
		logger.Errorf("app initialization failed [%v]", err)

		os.Exit(1)
	}

	// app run
	logger.Info("app run")
	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("app run error [%v]", err)
	}

	app.WG.Wait()

	// app close
	logger.Info("app close")
	if err := app.Close(); err != nil {
		logger.Errorf("app close error [%v]", err)

		os.Exit(1)
	}

	logger.Info("app shutdown")
	os.Exit(0)
}
