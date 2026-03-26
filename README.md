# hartool

A command-line tool for inspecting HAR (HTTP Archive) files. List, search, and examine captured HTTP traffic with pretty-printed output and live request replay.

## Install

```
go install hartool@latest
```

## Usage

```
hartool <file.har> [command] [args] [flags]
```

Commands can appear in any position. With no command, `hartool file.har` defaults to `list`.

### Commands

| Command | Description |
|---|---|
| `list` | List all entries with index, method, URL, and status |
| `exchange` | Print full request and response (headers + bodies) |
| `headers` | Print request and response headers |
| `req-headers` | Print request headers only |
| `resp-headers` | Print response headers only |
| `req-body` | Print request body |
| `resp-body` | Print response body |
| `search` | Search across URLs, headers, and bodies |
| `replay` | Re-issue a captured request and print the live response |
| `mcp` | Start an MCP server over stdio |

### Entry selection

Commands that operate on a single entry accept either:

- **Index** &mdash; numeric position from `list` output (e.g. `0`, `5`)
- **URL substring** &mdash; case-sensitive match against the request URL

### Examples

```sh
# List all entries
hartool capture.har

# Show full exchange for entry 3
hartool capture.har 3

# Show response body, pretty-printed
hartool capture.har resp-body 3 -p

# Search for a term (regex, case-insensitive)
hartool capture.har search "session" -r -i

# Replay a request live
hartool capture.har replay 3 -p
```

### Flags

| Flag | Applies to | Description |
|---|---|---|
| `-p, --pretty` | `req-body`, `resp-body`, `exchange`, `replay` | Pretty-print content (JSON, XML, HTML, YAML, form data) |
| `-r, --regex` | `search` | Treat search term as a regular expression |
| `-i, --case-insensitive` | `search` | Case-insensitive matching |

## MCP server

`hartool mcp` starts a [Model Context Protocol](https://modelcontextprotocol.io/) server over stdio, exposing two tools:

- **har_search** &mdash; list or search entries in a HAR file
- **har_dump** &mdash; dump headers or bodies for a specific entry

This lets AI assistants query HAR files directly.

## License

This program is free software; you can redistribute it and/or modify it under the terms of the [GNU General Public License](LICENSE) as published by the Free Software Foundation; either version 2 of the License, or (at your option) any later version.
