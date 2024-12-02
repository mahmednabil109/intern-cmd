package main

import (
	"fmt"

	"github.com/mahmednabil109/intern-cmd/pkg/core/loader"
)

var Preloaded = []*loader.Plugin{
	{
		Cmd: "g",
		Func: func(q string) (string, error) {
			return fmt.Sprintf("https://google.com/search?q=%s", q), nil
		},
	},
	{
		Cmd: "yt",
		Func: func(q string) (string, error) {
			if len(q) == 0 {
				return "https://www.youtube.com/", nil
			}
			return fmt.Sprintf("https://www.youtube.com/results?search_query=%s", q), nil
		},
	},
	{
		Cmd: "wiki",
		Func: func(q string) (string, error) {
			if len(q) == 0 {
				return "https://en.wikipedia.org/w/index.php", nil
			}
			return fmt.Sprintf("https://en.wikipedia.org/w/index.php?search=%s", q), nil
		},
	},
	{
		Cmd: "hn",
		Func: func(q string) (string, error) {
			if len(q) == 0 {
				return "https://hn.algolia.com/", nil
			}
			return fmt.Sprintf("https://hn.algolia.com/?q=%s", q), nil
		},
	},
	{
		Cmd: "note",
		Func: func(_ string) (string, error) {
			return "http://localhost:3001", nil
		},
	},
	{
		Cmd: "draw",
		Func: func(_ string) (string, error) {
			return "http://localhost:3002", nil
		},
	},
	{
		Cmd: "osdev",
		Func: func(q string) (string, error) {
			return fmt.Sprintf("https://wiki.osdev.org/index.php?search=%s", q), nil
		},
	},
	{
		Cmd: "godoc",
		Func: func(q string) (string, error) {
			return fmt.Sprintf("https://pkg.go.dev/search?q=%s", q), nil
		},
	},
	{
		Cmd: "sourcegraph",
		Func: func(q string) (string, error) {
			return fmt.Sprintf("https://sourcegraph.com/search?q=context:global %s", q), nil
		},
	},
	{
		Cmd: "cal",
		Func: func(_ string) (string, error) {
			return "https://calendar.google.com/calendar/", nil
		},
	},
}
