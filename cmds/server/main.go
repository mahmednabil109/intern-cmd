package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/mahmednabil109/intern-cmd/pkg/core/loader"
	"github.com/mahmednabil109/intern-cmd/pkg/core/runner"
)

var (
	userHomeDir, _ = os.UserHomeDir()
	port           = flag.Int("port", 3000, "port to run the daemon on.")
	configDir      = flag.String("config", path.Join(userHomeDir, ".config/intern-cmd"), "path to configuration directory.")
)

func main() {
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	defer func() {
		if err := recover(); err != nil {
			logger.Error("daemon panic-ed :()", "err", err)
		}
	}()

	if _, err := os.Stat(*configDir); err != nil {
		if err := os.MkdirAll(*configDir, 0777); err != nil {
			logger.Error("failed to create all directories in the path", "err", err, "path", *configDir)
		}
	}

	loader := &loader.GoPluginLoader{}
	runner, err := runner.New(logger, loader, runner.WithPreLoaded(Preloaded), runner.WithConfig(*configDir))
	if err != nil {
		logger.Error("failed to init command runner", "err", err)
		os.Exit(1)
	}
	defer runner.Close()

	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy :+1:!"))
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if len(query) == 0 {
			w.Write([]byte("what do you want?!"))
			return
		}
		target, err := runner.Exec(query)
		if err != nil {
			logger.Error("failed to execute command", "err", err)
		}
		logger.Info("executed command and redirecting", "target", target)
		http.Redirect(w, r, target, http.StatusSeeOther)
	})
	router.HandleFunc("/plugin", func(w http.ResponseWriter, r *http.Request) {
		// path to plugin for loading
		path := r.URL.Query().Get("path")

		if err := runner.Register(path); err != nil {
			logger.Error("failed to register new plugin", "err", err)
		} else {
			logger.Info("plugin registered successfully", "plugin", path)
		}
	})

	addr := fmt.Sprintf(":%d", *port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		signal := <-sigChan
		logger.Info("received signal, closing ...", "signal", signal)
		if err := runner.Close(); err != nil {
			logger.Error("error occurred, when closing runner", "err", err)
		}
		if err := server.Close(); err != nil {
			logger.Error("error occurred, when closing server", "err", err)
		}
	}()

	logger.Info("Start Listening, ...", "addr", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server ListenAndServe failed", "err", err)
	}

}
