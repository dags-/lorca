package lorca

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ChromeExecutable returns a string which points to the preferred Chrome
// executable file.
var ChromeExecutable = LocateChrome

// LocateChrome returns a path to the Chrome binary, or an empty string if
// Chrome installation is not found.
func LocateChrome() string {
	path := locateChrome()
	if path != "" {
		return path
	}

	paths := []string{
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Users/" + os.Getenv("USERNAME") + "/AppData/Local/Google/Chrome/Application/chrome.exe",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

// locateChrome uses cmd/shell commands to find chrome executable candidates
func locateChrome() string {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// try windows' "where" command to find possible chrome executables
		cmd = exec.Command("where", "chrome")
	} else {
		// try shell's "which" command to find possible chrome executables
		// on other OS's
		cmd = exec.Command("which", "-a", "chrome")
	}

	// write cmd's output into a buffer
	out := &bytes.Buffer{}
	cmd.Stdout = out
	e := cmd.Run()
	if e != nil {
		return ""
	}

	// drain buffer line by line for possible paths
	for out.Len() > 0 {
		path, e := out.ReadString('\n')

		// error if no new-line found
		if e != nil {
			return ""
		}

		// remove new-line from end of string
		if len(path) > 1 {
			path = path[:len(path)-2]
		}

		// check exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		return path
	}

	return ""
}

// PromptDownload asks user if he wants to download and install Chrome, and
// opens a download web page if the user agrees.
func PromptDownload() {
	title := "Chrome not found"
	text := "No Chrome/Chromium installation was found. Would you like to download and install it now?"

	// Ask user for confirmation
	if !messageBox(title, text) {
		return
	}

	// Open download page
	url := "https://www.google.com/chrome/"
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Run()
	case "darwin":
		exec.Command("open", url).Run()
	case "windows":
		r := strings.NewReplacer("&", "^&")
		exec.Command("cmd", "/c", "start", r.Replace(url)).Run()
	}
}
