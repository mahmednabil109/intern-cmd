package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	homeDir, _ = os.UserHomeDir()
	filePath   = flag.String("plugin", "", "path to plugin file (.go)")
	addr       = flag.String("addr", "http://localhost:3000", "addr of the daemon")
	configDir  = flag.String("config", path.Join(homeDir, ".config/intern-cmd/plugins/"), "path plugin directory")
)

func main() {
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if *filePath == "" {
		logger.Error("argument -plugin is required.")
		os.Exit(1)
	}
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		logger.Error("couldn't find file", "err", err)
		os.Exit(1)
	}
	if _, err := os.Stat(*configDir); err != nil {
		if err := os.MkdirAll(*configDir, 0777); err != nil {
			logger.Error("couldn't create all directories in the path", "err", err, "path", *configDir)
			os.Exit(1)
		}
	}
	pluginDir, pluginName := path.Split(*filePath)
	pluginName, _, _ = strings.Cut(pluginName, ".")
	outputFile := path.Join(*configDir, fmt.Sprintf("%s-plugin.so", pluginName))

	cmd := exec.Command("go", "build", "-buildmode", "plugin", "-o", outputFile)
	cmd.Dir = pluginDir
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		logger.Error("failed to run go build command", "err", err)
		os.Exit(1)
	}

	daemonURL, err := url.Parse(*addr)
	if err != nil {
		logger.Error("failed to parse url", "url", addr, "err", err)
	}
	daemonURL = daemonURL.JoinPath("plugin")
	query := daemonURL.Query()
	query.Add("path", outputFile)
	daemonURL.RawQuery = query.Encode()
	logger.Info("sending registration msg", "addr", daemonURL.String())
	resp, err := http.Get(daemonURL.String())
	if err != nil {
		logger.Error("failed to register plugin", "err", err)
	}
	defer resp.Body.Close()
	buffer := make([]byte, 1024)
	n, _ := resp.Body.Read(buffer)
	if resp.StatusCode != http.StatusOK {
		logger.Error("request to register plugin failed", "resp", string(buffer[:n]))
	}
	logger.Info("request succeeded", "resp", string(buffer[:n]))
}
