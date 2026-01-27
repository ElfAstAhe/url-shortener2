package main

import (
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/ElfAstAhe/url-shortener2/internal/app/bootstrap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printOrNA("Build version: %s\n", buildVersion)
	printOrNA("Build date: %s\n", buildDate)
	printOrNA("Build commit: %s\n", buildCommit)

	// app instance
	app := bootstrap.NewApp()
	//	defer app.Close()
	logger := app.Log.GetLogger("main")
	//	defer _utl.CloseOnly(logger.(io.Closer))

	// app initialization
	logger.Info("app init")
	if err := app.Init(); err != nil {
		logger.Errorf("app initialization failed [%v]", err)
		defer app.Close()

		//		os.Exit(1)
		panic(errors.New("app initialization failed"))
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

		//		os.Exit(1)
		panic(errors.New("app close failed"))
	}

	logger.Info("app shutdown")
	//	os.Exit(0)
}

func printOrNA(template, val string) {
	if strings.TrimSpace(val) == "" {
		fmt.Printf(template, "N/A")

		return
	}

	fmt.Printf(template, val)
}
