package chromeapp

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"bitbucket.org/kardianos/osext"
	"github.com/googollee/go-socket.io"
	"github.com/mgutz/str"
)

// Options for starting chrome.
type Options struct {
	// Width of the window.
	Width int
	// Height of the window.
	Height int
	// BrowserDataDir is where Chrome stores the app's data and profile.
	BrowserDataDir string
	// ChromeExecutable is path to Chrome executable.
	ChromeExecutable string
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// NewOptions creates an Options structure with default values.
func NewOptions() *Options {
	filename, err := osext.Executable()
	if err != nil {
		panic(fmt.Errorf("Could not get path to executable"))
	}
	log.Println(filename)
	executableDir := filepath.Dir(filename)

	chrome, err := chromeExecutable()
	if err != nil {
		panic(err)
	}

	return &Options{
		Width:            1024,
		Height:           768,
		BrowserDataDir:   filepath.Join(executableDir, ".chromeapp"),
		ChromeExecutable: chrome,
	}

}

func Simple(listener func(*socketio.Server)) {
	Start(NewOptions(), listener)
}

// Start starts chrome, HTTP and websocket server.
func Start(options *Options, handler func(*socketio.Server)) {
	if !fileExists(options.ChromeExecutable) {
		panic(fmt.Errorf(`Invalid path to chrome executable "%s"`, options.ChromeExecutable))
	}

	// use a random port
	l, err := net.Listen("tcp", ":0")
	defer l.Close()
	addr := l.Addr().String()
	port := str.Between(addr, "]:", "")

	args := []string{
		"--active-on-launch",
		"--app=http://127.0.0.1:" + port + "/",
		"--app-window-size=" + strconv.Itoa(options.Width) + "," + strconv.Itoa(options.Height),
		"--enable-crxless-web-apps",
		"--force-app-mode",
		"--no-default-browser-check",
		"--no-first-run",
		"--user-data-dir=" + options.BrowserDataDir,
	}
	spawn(options.ChromeExecutable, args)

	server, err := socketio.NewServer(nil)
	if err != nil {
		panic(err)
	}
	handler(server)

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at " + l.Addr().String() + "...")
	log.Fatal(http.Serve(l, nil))
}

func chromeExecutable() (string, error) {
	executable := ""
	switch runtime.GOOS {
	default:
		return "", fmt.Errorf("Could not get chrome executable for GOOS=%s", runtime.GOOS)
	case "darwin":
		executable = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	case "linux":
		executable = "/usr/bin/google-chrome"
		if !fileExists(executable) {
			executable = "/usr/bin/chromium-browser"
		}
	}
	return executable, nil
}
