<h1 align="center">Trainings logger</h1>

<p align="center">
  <a href="https://github.com/danielsavinoff/trainings-logger/actions/workflows/main.yml">
    <img src="https://github.com/danielsavinoff/trainings-logger/actions/workflows/main.yml/badge.svg" alt="Release">
  </a>
</p>

<p align="center">
Track your workout progress in the CLI while waiting for your agent to finish in another tab.
</p>

## ✨ Features

- 💾 **Your data is yours.** Stored locally and under your control.
- 💻 **CLI first.** Log your workouts without leaving the terminal.
- ⌨️ **Vim keybindings.** Navigate with the keys you already know.
- ☁️ **External storage.** Back up and restore your workout data with Dropbox.
- 📊 **Analytics and advice.** Planned for a future release.

## 📦 Installation

### macOS

Install with [Homebrew](https://brew.sh):

```sh
brew install --cask danielsavinoff/tap/trainings-logger
trainings-logger
```

To update:

```sh
brew upgrade --cask trainings-logger
```

The Homebrew package is updated automatically with each stable release.

### Linux and Windows

Download the archive for your operating system and architecture from
[Releases](https://github.com/danielsavinoff/trainings-logger/releases), then extract it.

On Linux, run:

```sh
chmod +x trainings-logger
./trainings-logger
```

On Windows, open PowerShell in the extracted folder and run:

```powershell
.\trainings-logger.exe
```

To run it from any directory, move the executable into a folder on your `PATH`.
