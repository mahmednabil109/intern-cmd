## <Placeholder for a cool name, now it's just `Intern-Cmd`>

![](https://imgs.xkcd.com/comics/donald_knuth.png)

`Intern-Cmd` is a tool that mainly leverages browser's custom search engines, to make an extensible plugin system (using golang plugin builtin features) to define different shortcuts.

### Usage:
- build/run the server from `cmd/server/main.go`
- in your favorite browser define a custom search engine, referring to `http://localhost:300/?q=%s`.
- lastly, in your search bar type `i g blink html`

### Define custom plugins:
- write a golang plugin just like [test-plugin.go](./test-plugins/test-plugin.go).
- build/run go command from `cmd/plugin/add.go` and run it with plugin file.
- should work, I hope...
