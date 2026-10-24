// START_MODULE_CONTRACT
//   PURPOSE: Open Dropbox consent in the system default browser.
//   SCOPE: Native default-browser commands.
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
	"os/exec"
	"runtime"
)

// START_BLOCK_BROWSER_LAUNCH
func browserCommand(
	ctx context.Context,
	platform string,
	address string,
) *exec.Cmd {
	if platform == "darwin" {
		return exec.CommandContext(
			ctx,
			"open",
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
		authorizationURL,
	).Run()
}

// END_BLOCK_BROWSER_LAUNCH
