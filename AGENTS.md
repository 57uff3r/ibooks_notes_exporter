# ibooks-notes-exporter

## What is this

A command-line tool that exports your highlights and notes from Apple iBooks into Markdown files. It works on macOS only (Intel and Apple Silicon).

## Why this exists

Apple iBooks has no built-in way to export your notes. If you highlight passages or write notes while reading, the only way to get them out is manual copy-paste. This tool solves that problem by reading the iBooks database directly and printing everything as clean Markdown.

The target audience is anyone who reads books in iBooks and wants to keep their notes outside of the app — in Obsidian, Notion, a plain text archive, or anywhere else.

## How it works

iBooks stores its data in two SQLite databases inside `~/Library/Containers/com.apple.iBooksX/Data/Documents/`:
- `BKLibrary/` — contains book metadata (titles, authors, asset IDs)
- `AEAnnotation/` — contains highlights and notes

The tool opens the books database, attaches the annotations database, and runs SQL queries that join the two. No network requests, no Apple APIs — just local SQLite reads.

Only EPUB books are supported. PDFs use a different annotation system that this tool does not handle.

## Project structure

```
main.go              — CLI entry point, two commands: `books` and `export`
db/connection.go     — finds and opens the two SQLite databases
db/db_queries.go     — SQL queries (get all books, get book by ID, get notes)
db/structs.go        — data structs for books and highlights
.goreleaser.yml      — release config, builds macOS-only binary with CGO
.github/workflows/   — GitHub Actions runs goreleaser on version tags
Makefile             — shortcuts for goreleaser (test release, publish)
```

## Commands

- `ibooks_notes_exporter books` — lists all books that have highlights or notes, shown in a table with book ID, note count, and title
- `ibooks_notes_exporter export --book_id <ID>` — exports all highlights and notes from a specific book as Markdown to stdout
- `--skip_first_x_notes N` — optional flag to skip the first N notes (useful for partial re-exports)

Typical workflow: run `books` to find the ID, then run `export` and redirect stdout to a `.md` file.

## Build and release

- Language: Go 1.19
- CGO is required because of the `go-sqlite3` driver
- Distributed via Homebrew: `brew install 57uff3r/mac-apps/ibooks_notes_exporter`
- Releases are created by pushing a `v*` tag, which triggers goreleaser via GitHub Actions
- goreleaser builds the binary and publishes it to the `57uff3r/homebrew-mac-apps` tap

## Key dependencies

- `github.com/mattn/go-sqlite3` — SQLite driver (requires CGO)
- `github.com/urfave/cli/v2` — CLI framework
- `github.com/jedib0t/go-pretty/v6` — table rendering for the `books` command

## Things to know

- The tool reads local iBooks databases directly. If Apple changes the database schema or file locations in a future macOS update, queries will break.
- There are no tests in this project.
- The sort order for exported notes uses `ZPLLOCATIONRANGESTART` (position in book) as the primary sort key, with creation date as a tiebreaker. This was a pragmatic fix — proper sorting would require parsing the EPUB CFI standard (ISO/IEC 23736-6:2020).
- Author names are shortened to last names in the book list for compact display.

## Writing style

All documentation and comments in this project should be written in plain English at B2 level. Keep sentences short. Avoid jargon when a simple word works. The goal is for any developer to understand this project quickly, regardless of their native language.
