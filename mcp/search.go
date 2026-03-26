package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hartool/har"
	"hartool/internal"
)

func registerSearchTool(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "har_search",
		Description: "List or search entries in a HAR file. Omit term to list all.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file":             {"type": "string", "description": "path to the HAR file"},
				"term":             {"type": "string", "description": "search term (omit to list all entries)"},
				"regex":            {"type": "boolean", "description": "treat term as regex"},
				"case_insensitive": {"type": "boolean", "description": "case-insensitive matching"}
			},
			"required": ["file"]
		}`),
	}, handleSearch)
}

func handleSearch(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		File            string `json:"file"`
		Term            string `json:"term"`
		Regex           bool   `json:"regex"`
		CaseInsensitive bool   `json:"case_insensitive"`
	}
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return errResult("invalid arguments: %v", err), nil
	}

	h, err := har.Load(args.File)
	if err != nil {
		return errResult("load: %v", err), nil
	}

	// No term: list all entries.
	if args.Term == "" {
		var sb strings.Builder
		for i, e := range h.Log.Entries {
			fmt.Fprintf(&sb, "[%d] %s %s (%d)\n", i, e.Request.Method, e.Request.URL, e.Response.Status)
		}
		return textResult(sb.String()), nil
	}

	// Search with term.
	match, err := internal.BuildMatcher(args.Term, args.Regex, args.CaseInsensitive)
	if err != nil {
		return errResult("%v", err), nil
	}

	var sb strings.Builder
	for i, entry := range h.Log.Entries {
		if fields := internal.SearchEntry(entry, match); len(fields) > 0 {
			fmt.Fprintf(&sb, "[%d] %s (%s)\n", i, entry.Request.URL, strings.Join(fields, ","))
		}
	}
	if sb.Len() == 0 {
		return textResult("no matches found"), nil
	}
	return textResult(sb.String()), nil
}
