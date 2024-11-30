package runner

import (
	"fmt"
	"log/slog"
	"path"
	"strings"
	"sync"

	"github.com/mahmednabil109/intern-cmd/pkg/config"
	"github.com/mahmednabil109/intern-cmd/pkg/core/loader"
)

type Config struct {
	preloaded []*loader.Plugin
	configDir string
}

type Option interface {
	apply(*Config)
}

type WithPreLoaded []*loader.Plugin

func (o WithPreLoaded) apply(cfg *Config) {
	cfg.preloaded = ([]*loader.Plugin)(o)
}

type WithConfig string

func (o WithConfig) apply(cfg *Config) {
	cfg.configDir = string(o)
}

type Runner struct {
	cfg    *Config
	logger *slog.Logger
	loader loader.Loader

	// TODO: replace this map with a LOCK-FREE TRIE, might be faster.
	cmds   map[string]*loader.Plugin
	closed bool
	lock   sync.RWMutex
}

func New(log *slog.Logger, ld loader.Loader, opts ...Option) (*Runner, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	cmds := make((map[string]*loader.Plugin), len(cfg.preloaded))
	if len(cfg.preloaded) != 0 {
		for _, plugin := range cfg.preloaded {
			cmds[plugin.Cmd] = plugin
		}
	}

	// load previously loaded plugins
	if len(cfg.configDir) != 0 {
		configFile, err := config.Load(path.Join(cfg.configDir, "loaded-plugins.json"))
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
		for cmd, path := range configFile.Plugins {
			if plugin, err := ld.Get(path); err == nil {
				cmds[cmd] = plugin
			} else {
				log.Error("failed to load plugin", "cmd", cmd, "plugin", plugin)
				// return nil, fmt.Errorf("failed to get plugin %s: %w", path, err)
			}
		}
	}

	return &Runner{
		cfg:    cfg,
		logger: log,
		loader: ld,
		cmds:   cmds,
	}, nil
}

func (r *Runner) Exec(q string) (string, error) {
	cmd, rest, _ := strings.Cut(q, " ")

	r.lock.RLock()
	plugin, found := r.cmds[cmd]
	r.lock.RUnlock()
	if !found {
		return "", fmt.Errorf("no plugin found for command %s", cmd)
	}

	return plugin.Exec(rest)
}

func (r *Runner) Register(path string) error {
	plugin, err := r.loader.Get(path)
	if err != nil {
		return fmt.Errorf("failed to load plugin %s: %w", path, err)
	}

	r.lock.Lock()
	r.cmds[plugin.Cmd] = plugin
	r.lock.Unlock()

	return nil
}

func (r *Runner) Close() error {
	if r.closed {
		return nil
	}
	configFile := config.File{
		Plugins: make(map[string]string, 0),
	}
	for cmd, plugin := range r.cmds {
		if plugin.Path != "" {
			configFile.Plugins[cmd] = plugin.Path
		}
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	err := config.Save(&configFile, path.Join(r.cfg.configDir, "loaded-plugins.json"))
	r.closed = err != nil
	return err
}
