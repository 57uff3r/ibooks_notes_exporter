# ibooks-notes-exporter

People remember less than 20% of what they read if they are not taking notes.
That's why it's always a good idea to take notes while you are reading.
But if you are reading in iBooks and taking notes there, it's very hard to get them out of that app.
You have to manually copy-paste quotes and your thoughts. ibooks-notes-exporter  solves
that problem. it's a little program that gets all your notes and highlights from iBooks and
exports them into markdown.

Note! **ibooks-notes-exporter** can only extract notes from EPUB files.


## Installation

ibooks-notes-exporter is available on OS X (both Intel and M-series processors).
It's distributed via a [homebrew](https://brew.sh/) package manager.

Run these commands in your terminal 

```shell
brew tap 57uff3r/mac-apps
brew install 57uff3r/mac-apps/ibooks_notes_exporter
```


## Usage

### CLI

First of all, you have to get a list of all your books with notes and highlights.

```shell
❯ ibooks_notes_exporter books
+----------------------------------+-----------------+----------------------------------------------------------------------------------+
| SINGLEBOOK ID                    | NUMBER OF NOTES | TITLE AND AUTHOR                                                                 |
+----------------------------------+-----------------+----------------------------------------------------------------------------------+
| 4BAE5DA3C95788753173EAE8C63E6034 |               1 | Lorem impsum — John Doe                                                          |
| 7C3FA4F94689D97444BB4E0FD97D7197 |              54 | Hamlet — william shakespeare                                                     |
+----------------------------------+-----------------+----------------------------------------------------------------------------------+
```

And then you can export all your notes from the book into a markdown file.


```shell
ibooks_notes_exporter export --book_id 4BAE5DA3C95788753173EAE8C63E6034 > ./LoremImpsum.md
```

You also can export your notes partially by skipping  first X notes:

```shell
ibooks_notes_exporter export --book_id 4BAE5DA3C95788753173EAE8C63E6034 --skip_first_x_notes 20 > ./LoremImpsum.md 
```
In this example, the exporter will skip first 20 notes and export the rest.


### MCP Server

ibooks-notes-exporter can work as an [MCP](https://modelcontextprotocol.io/) (Model Context Protocol) server. This allows AI assistants to read your iBooks highlights and notes directly.

The MCP server uses stdio transport. You start it with:

```shell
ibooks_notes_exporter mcp
```

The server provides three tools:

| Tool | Description |
|------|-------------|
| `list_books` | List all books that have highlights or notes. Returns book IDs, titles, authors, and note counts. |
| `get_notes` | Get all highlights and notes from a specific book by its ID. Returns Markdown output. |
| `search_notes` | Search across all highlights and notes for a keyword or phrase. Case-insensitive. |

All tools are read-only. They don't modify your iBooks data.


#### Claude Code / Claude Desktop

Add this to your MCP config file:

- Claude Code: `.mcp.json` in your project directory
- Claude Desktop: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "ibooks": {
      "command": "ibooks_notes_exporter",
      "args": ["mcp"]
    }
  }
}
```

If you installed via Homebrew, the binary is already in your PATH. Otherwise, use the full path to the binary.


#### ChatGPT Desktop

Open Settings > Tools > Add MCP Tool, then select "Add manually" and enter:

- **Transport:** `stdio`
- **Command:** `ibooks_notes_exporter`
- **Arguments:** `mcp`

Or add it to `~/.config/chatgpt/mcp.json`:

```json
{
  "mcpServers": {
    "ibooks": {
      "command": "ibooks_notes_exporter",
      "args": ["mcp"]
    }
  }
}
```


#### Cursor

Open Settings > MCP Servers > Add Server, then enter:

- **Name:** `ibooks`
- **Type:** `command`
- **Command:** `ibooks_notes_exporter mcp`

Or add it manually to `.cursor/mcp.json` in your project:

```json
{
  "mcpServers": {
    "ibooks": {
      "command": "ibooks_notes_exporter",
      "args": ["mcp"]
    }
  }
}
```


#### Windsurf

Add to `~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "ibooks": {
      "command": "ibooks_notes_exporter",
      "args": ["mcp"]
    }
  }
}
```


#### VS Code (Copilot)

Add to your VS Code `settings.json`:

```json
{
  "mcp": {
    "servers": {
      "ibooks": {
        "command": "ibooks_notes_exporter",
        "args": ["mcp"]
      }
    }
  }
}
```


#### Any MCP-compatible client

The server uses **stdio transport** — it reads JSON-RPC messages from stdin and writes responses to stdout. Any MCP client that supports the stdio transport can connect by running:

```shell
ibooks_notes_exporter mcp
```

No network ports or authentication are needed. The server runs locally and reads only your iBooks SQLite databases.


## How it works

iBooks stores data in two SQLite databases inside `~/Library/Containers/com.apple.iBooksX/Data/Documents/`:

- **BKLibrary/** — book metadata (titles, authors, asset IDs)
- **AEAnnotation/** — highlights and notes

The tool opens the books database, attaches the annotations database, and runs SQL queries that join the two. No network requests, no Apple APIs — just local SQLite reads.


## Feedback and contribution

Your feedback and pull requests are much appreciated.
Feel free to send your comments and thoughts to [me@akorchak.software](mailto:me@akorchak.software)


## Changelog


**0.0.6**

* Updated Go from 1.19 to 1.23
* Migrated CLI framework from urfave/cli v2 to v3
* Switched SQLite driver from mattn/go-sqlite3 to modernc.org/sqlite (pure Go, no CGO needed)
* Updated go-pretty from v6.4.4 to v6.7.10
* Updated goreleaser config and GitHub Actions to current versions
* Added MCP server mode (`ibooks_notes_exporter mcp`)

**0.0.5**

* New CLI syntax
* **--skip_first_x_notes** flag
* Fix for [Issue #5](https://github.com/57uff3r/ibooks_notes_exporter/issues/5) ('some bug in chinese book')

**0.0.4**

Fix for long titles made by @[NSBum](https://github.com/NSBum)

**0.0.3**

Bug fix: worng  order of notes and highlights. Quick and dirty fix, better solution requires to implement a parser 
for ISO/IEC 23736-6:2020 standard (EPUB Canonical Fragment Identifier or epubcfi) and this will be done in next 
versions.

**0.0.2**

Markdown fix: missing line break


**0.0.1**

Initial release
