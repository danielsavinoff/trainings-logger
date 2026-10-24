# Trainings logger

Run the terminal application with `go run ./cmd/cli`.

On first start, setup asks for kg/lb and an optional Dropbox account. Enter skips
an answer; skipping everything saves nothing. Press `r` from records to reopen
setup without deleting existing settings.

## Records and sessions

A session is a complete workout. Add reusable catalog exercises to that session,
then record sets under each performance.

Sessions `<1>` occupy the left fifth of records. Performances `<2>` above sets `<3>`
share the remaining four fifths, with narrow gaps between panes. Press `1`–`3` to
focus them; Enter advances from a session to its exercises and then its sets.

Create a workout immediately with `1`, then `n`; no form is needed. In Performances,
`n` opens the exercise catalog full screen. Enter adds/selects the highlighted catalog exercise and returns to
Performances; repeated selection uses the existing performance. Esc returns without adding
a performance and preserves the records selection and errors.

Inside the catalog, `n/e/d` create, edit, and delete reusable definitions. Saving
or cancelling a form returns to the catalog; press Enter separately to select.
Saved catalog edits persist when you return, and referenced definitions cannot be
deleted. Catalog navigation and shortcuts do not trigger records actions.

Use `n/d` to create or delete sessions (`e` has no action). Use `n/e/d` for sets.
In Performances, `d` removes only that workout's performance and its sets; performances have no independent edit form.
Deleting a workout cascades its performances and sets. Deleting the final set or performance
preserves its parent. Catalog definitions remain available for other workouts.

Sessions have no title. Session rows show only their index, assigned as the largest
existing index plus one, starting at 1. Deleting a session does not renumber others.
The performance and catalog tables show names and body weight in percent; catalog edits refresh those values in workout performances. Set rows show
reps, weight, and rest without date separators. The records footer shows only
`? - help` as a shortcut hint, with the same outer padding as tables. The catalog
has no shortcut footer; actual errors remain visible. Every view shares zero
vertical and one-cell horizontal padding, with compact gaps in help, setup, and
forms. Help places Navigation and Actions on the left and Help on the right,
stacking them on compact screens. Body-weight share is editable in the catalog form. The schema is supplied by a single initial migration for fresh
databases; older databases need a separate migration strategy or a new database.
The application does not reset them automatically. Backup format version 4 includes
workouts and performances; earlier versions cannot be restored.

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

From records, `s` saves catalog exercises, sessions, performances, sets, and optional preferences to
`/trainings-logger/data.json` in the app's Dropbox folder. `l` restores that backup, replacing
local catalog exercises, sessions, performances, sets, and preferences. Identity credentials stay local. Saving
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

## Migrations

Create new timestamped migration files with golang-migrate:

```sh
migrate create -ext sql -dir migrations/sqlite -tz Europe/Moscow change_name
```

The current `init` migration creates all seven application tables and their indexes.
Its down migration removes those tables in dependency-safe order.

## Checks

```sh
go test ./...
go vet ./...
```

## Domain layout

`internal/domain/session` contains workout sessions. Its `exercise` package holds
reusable catalog definitions; the sibling `exercise/set` and `exercise/performance`
packages contain sets and performances. A performance links a session and a catalog
exercise, owns its sets, and is stored in the `performances` table. Sets reference
`performance_id`. Backup version 4 uses the same names; earlier backups are rejected.
