// START_MODULE_CONTRACT
//   PURPOSE: Open Dropbox consent in a chosen system browser.
//   SCOPE: Native browser commands and optional DROPBOX_BROWSER override.
//   DEPENDS: none
//   LINKS: M-INFRA-DROPBOX-AUTH, V-M-INFRA-DROPBOX-AUTH
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   browserCommand - build browser arguments without shell evaluation
//   openBrowser - run the browser launcher with cancellation
// END_MODULE_MAP

package dropbox

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// START_BLOCK_BROWSER_LAUNCH
func browserCommand(
	ctx context.Context,
	platform string,
	browser string,
	address string,
) *exec.Cmd {
	browser = strings.TrimSpace(browser)
	if platform == "darwin" {
		if browser != "" {
			return exec.CommandContext(
				ctx,
				"open",
				"-a",
				browser,
				address,
			)
		}
		return exec.CommandContext(
			ctx,
			"open",
			address,
		)
	}
	if browser != "" {
		return exec.CommandContext(
			ctx,
			browser,
			address,
		)
	}
	if platform == "windows" {
		return exec.CommandContext(
			ctx,
			"rundll32",
			"url.dll,FileProtocolHandler",
			address,
		)
	}
	return exec.CommandContext(
		ctx,
		"xdg-open",
		address,
	)
}

func openBrowser(
	ctx context.Context,
	authorizationURL string,
) error {
	return browserCommand(
		ctx,
		runtime.GOOS,
		os.Getenv("DROPBOX_BROWSER"),
		authorizationURL,
	).Run()
}

// END_BLOCK_BROWSER_LAUNCH
