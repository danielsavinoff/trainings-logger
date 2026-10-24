# Trainings logger

Run the terminal application with `go run ./cmd/cli`.

On first start, setup asks for kg/lb and an optional Dropbox account. Enter skips
an answer; skipping everything saves nothing. Press `r` from records to reopen
setup without deleting existing settings.

## Dropbox

Create a scoped Dropbox app in the [App Console](https://www.dropbox.com/developers/apps).
App folder access is sufficient. In its Permissions tab enable and submit:

- `account_info.read`
- `files.content.read`
- `files.content.write`

Set the app's **App key** before launching:

```sh
export DROPBOX_APP_KEY="your_app_key"
go run ./cmd/cli
```

Sign-in opens the browser and asks you to paste the authorization code back into
the terminal. The application uses PKCE and needs no app secret or redirect URI.
Access and refresh tokens stay in the local SQLite database `data.db` in the operating
system's `trainingslogger` configuration directory. The application restricts database
files to the current user.

From records, `s` saves exercises, sets, and optional preferences to
`/trainings-logger/data.json` in the app's Dropbox folder. `l` restores that backup, replacing
local exercises, sets, and preferences. Identity credentials stay local. Saving
again overwrites the remote backup; this is manual backup and restore.

## Browser recovery

In the code-entry step, `ctrl+r` starts a fresh authorization attempt. `ctrl+o`
shows the full sign-in URL in a scrollable view; use `j`/`k` and Escape to return.
You can open that URL in another browser.

Authorization always opens the system default browser.

The URL `/1/oauth2/authorize_submit` is part of Dropbox's consent submission.
An error there occurs before the application receives an authorization code;
the URL alone does not identify its cause. The reported Safari failure has not
been reproduced. Record Safari's full error text when diagnosing it, without
sharing app keys, authorization codes, or tokens.

Dropbox documents the manual code and PKCE flows in its
[OAuth guide](https://docs.dropboxapi.com/dropbox-api/docs/oauth).

## Checks

```sh
go test ./...
go vet ./...
```
