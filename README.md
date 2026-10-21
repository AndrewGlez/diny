<div align="center">

<img src="https://diny.run/diny-v2-200-light.png" alt="diny logo" width="200"/>

# diny

**AI git commit messages using a configurable OpenAI-compatible provider.**

</div>

Open-source CLI that turns staged diffs into clean commits in one command.

---

## Fork configuration

This fork replaces the hosted Diny backend with an OpenAI-compatible provider.
Commit messages, split plans, changelogs, and timeline analysis use this configuration
in `~/.config/diny/config.yaml`:

```yaml
provider:
  base_url: https://ollama.com/v1
  model: deepseek-v4.1-flash:1017
  api_key_env: OLLAMA_API_KEY
  json_mode: true
  timeout_seconds: 120
```

Export the named environment variable with your provider key. For a keyless local
server, set `api_key_env: ""`. Hosted feedback prompts are disabled by default.
Split plans still require each staged file to belong to exactly one group;
conflicting assignments are rejected, not silently discarded.

## Features

- Interactive TUI for commit, config, changelog, timeline, and yolo
- Reads staged changes with `git diff` and filters out noise (lockfiles, binaries, build artifacts)
- Generates 3 message variants; inline editing or open in `$EDITOR`
- Split staged changes into multiple commits grouped by concern
- File picker to stage/unstage without leaving diny
- Timeline analysis with date presets or custom ranges
- AI-powered changelog generation between tags or commits
- Three-tier config system (global, project-shared, project-private)
- 10+ dark and light themes
- Configurable provider endpoint, model, API-key environment variable, and timeout

## Install

Build this fork from source with Go:

```bash
git clone --branch nixos-provider https://github.com/AndrewGlez/diny.git
cd diny
go build -o diny .
```

Run `./diny`, or install the binary into a directory on your `PATH`.
Upstream Homebrew, Scoop, and release binaries do not include this fork's provider support.

## Usage

Stage your changes, then run:

```bash
git add .
diny commit
```

This launches the interactive TUI — generate, review, and commit without leaving the terminal.

Staged too much at once? Press `S` on the ready screen to split the diff into multiple commits grouped by concern — edit, reassign files, regenerate, then commit them in order.

## Commands

| Command | Description |
|---------|-------------|
| `diny commit` | Launch the interactive TUI |
| `diny yolo` | Stage all changes, generate a commit, and push |
| `diny changelog` | Generate an AI-powered changelog between tags or commits |
| `diny timeline` | Summarize and analyze your commit history |
| `diny config` | Interactive TUI config editor |
| `diny theme` | List available UI themes |
| `diny auto` | Set up a `git auto` alias |
| `diny link lazygit` | Integrate diny with LazyGit |
| `diny update` | Update diny to the latest version |

## Integrations

### Git alias (`git auto`)

```bash
diny auto          # install
diny auto remove   # uninstall
```

Then use `git auto` anywhere you'd use `git commit`.

### LazyGit

```bash
diny link lazygit
```

Adds a custom command to [LazyGit](https://github.com/jesseduffield/lazygit) so you can generate commit messages from its UI.

## Configuration

Run `diny config` to open a TUI for browsing and editing all settings.

diny supports a three-tier configuration system:

| Level | Location | Description | Committed |
|-------|----------|-------------|-----------|
| Global | `~/.config/diny/config.yaml` | Your defaults for all projects | No |
| Project (versioned) | `.diny.yaml` | Team-shared settings for this project | Yes |
| Project (local) | `<gitdir>/diny/config.yaml` | Your personal per-project overrides | No |

**Priority:** local > versioned > global. Project configs only need to specify the values they override.

Outside a git repository, only global config is available.

### Options

| Option | Description | Values |
|--------|-------------|--------|
| `theme` | UI color theme | see [Themes](#themes) |
| `commit.conventional` | Use conventional commit format | `true` / `false` |
| `commit.emoji` | Prefix commits with emoji | `true` / `false` |
| `commit.tone` | Message tone | `professional` / `casual` / `friendly` |
| `commit.length` | Message length | `short` / `normal` / `long` |
| `commit.custom_instructions` | Extra guidance for the AI | free text |
| `commit.hash_after_commit` | Show and copy commit hash after committing | `true` / `false` |

### Themes

- **Dark:** `catppuccin`, `tokyo`, `nord`, `dracula`, `gruvbox-dark`, `onedark`, `monokai`, `solarized-dark`, `everforest-dark`, `flexoki-dark`
- **Light:** `solarized-light`, `github-light`, `gruvbox-light`, `flexoki-light`

Run `diny theme` to preview them all.

### Example

```yaml
theme: tokyo

commit:
  conventional: true
  emoji: true
  tone: casual
  length: short
  custom_instructions: "Include JIRA ticket from branch name"
  hash_after_commit: true
```

## Update

```bash
diny update                                     # built-in updater
brew update && brew upgrade dinoDanic/tap/diny  # macOS / Linux
scoop update diny                               # Windows
```

## Why diny

I'm terrible at commit messages — I type fast, make mistakes, and my history ends up full of "fix stuff" and "ui thing." diny reads your changes, writes a decent message, and summarizes your day when you need to remember what you worked on. It's AI doing what it's actually good at: the repetitive stuff.

## License

[MIT](LICENSE)
