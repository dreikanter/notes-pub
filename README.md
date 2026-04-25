# npub

A static site builder for Markdown notes. Reads notes from a local directory, renders them to HTML with syntax highlighting, and outputs a complete static site with tag pages and an Atom feed.

## Prerequisites

- Go 1.25+
- Node.js (for Tailwind CSS)

## Install

```sh
go install github.com/dreikanter/npub/cmd/npub@latest
```

## Build

Install dependencies and build:

```sh
npm install
make build
```

`npm install` is only needed once (or when dependencies change). `make build` compiles the Tailwind CSS stylesheet and builds the `npub` binary.

## Configuration

Generate a starter `npub.yml` in the current directory (or pass a target path):

```sh
npub init
npub init path/to/project
```

The generated file lists every option set to its default and commented out — uncomment and edit what you need. Required fields are `site_root_url`, `site_name`, and `author_name`.

All values can be overridden with CLI flags:

| Config option | CLI flag | Default | Required |
|---|---|---|---|
| `notes_path` | `--notes` | `$NOTES_PATH` | |
| `assets_path` | `--assets` | `<notes_path>/images` | |
| `build_path` | `--out` | `./dist` | |
| `static_path` | `--static` | `<notes_path>/static` | |
| `site_root_url` | `--url` | | Yes |
| `site_name` | `--site-name` | | Yes |
| `author_name` | `--author` | | Yes |
| `license_name` | `--license-name` | CC BY 4.0 | |
| `license_url` | `--license-url` | https://creativecommons.org/licenses/by/4.0/ | |
| `intro` | | | |

Priority: CLI flags > YAML config.

Config file discovery order:

1. `--config` flag
2. `NPUB_CONFIG` env var
3. `npub.yml` inside `$NOTES_PATH` (or `--notes` value) if it exists
4. `npub.yml` in the current directory

See `npub.yml.sample` in the repo for the same starting template.

The optional `intro` field renders as a paragraph above the posts list on the index page. Leave it empty or unset to omit.

### Image assets

Downloaded images are cached in an assets directory, organized by note UID. By default this is the `images` subdirectory of `notes_path`. Override with `assets_path` in the config or `--assets` flag.

### Static files

Files in the `static` subdirectory of `notes_path` are copied as-is to the build output root. Use this for `CNAME`, `robots.txt`, `favicon.ico`, etc. Override with `static_path` in the config or `--static` flag.

## Usage

Build the site:

```sh
npub build
```

Serve locally:

```sh
npub serve
```

The `serve` command starts a local HTTP server on port 4000 (override with `--port`), serving from `dist` (override with `--dir`).

## Notes format

Notes are Markdown files managed by [notes-cli](https://github.com/dreikanter/notes-cli). A note becomes part of the published site when its frontmatter includes `public: true`.

## Versioning

`CHANGELOG.md` is the source of truth for the version. On PR merge, GitHub
Actions (`.github/workflows/tag.yml`) reads the topmost `## [X.Y.Z]` heading
from `CHANGELOG.md` and pushes `vX.Y.Z` as a git tag. Bump major/minor/patch
by writing the desired heading in the PR.
