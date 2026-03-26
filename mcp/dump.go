package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hartool/har"
	"hartool/internal"
)

func registerDumpTool(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "har_dump",
		Description: "Dump data from a specific request/response in a HAR file.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file":  {"type": "string", "description": "path to the HAR file"},
				"entry": {"type": "string", "description": "entry index or URL substring"},
				"view":  {"type": "string", "enum": ["all","headers","req-headers","resp-headers","req-body","resp-body"], "description": "what to display"}
			},
			"required": ["file", "entry", "view"]
		}`),
	}, handleDump)
}

func handleDump(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		File  string `json:"file"`
		Entry string `json:"entry"`
		View  string `json:"view"`
	}
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return errResult("invalid arguments: %v", err), nil
	}

	h, err := har.Load(args.File)
	if err != nil {
		return errResult("load: %v", err), nil
	}
	entry, err := internal.ResolveEntry(h.Log.Entries, args.Entry)
	if err != nil {
		return errResult("%v", err), nil
	}

	var buf bytes.Buffer
	switch args.View {
	case "all":
		internal.PrintRequestHeaders(&buf, entry.Request)
		fmt.Fprintln(&buf)
		internal.WriteRequestBody(&buf, entry.Request.PostData)
		fmt.Fprintln(&buf)
		internal.PrintResponseHeaders(&buf, entry.Response)
		fmt.Fprintln(&buf)
		internal.WriteResponseBody(&buf, entry.Response.Content)
	case "headers":
		fmt.Fprintln(&buf, "=== Request ===")
		internal.PrintRequestHeaders(&buf, entry.Request)
		fmt.Fprintln(&buf, "\n=== Response ===")
		internal.PrintResponseHeaders(&buf, entry.Response)
	case "req-headers":
		internal.PrintRequestHeaders(&buf, entry.Request)
	case "resp-headers":
		internal.PrintResponseHeaders(&buf, entry.Response)
	case "req-body":
		internal.WriteRequestBody(&buf, entry.Request.PostData)
	case "resp-body":
		internal.WriteResponseBody(&buf, entry.Response.Content)
	default:
		return errResult("unknown view: %s", args.View), nil
	}

	return textResult(buf.String()), nil
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errResult(format string, args ...any) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(format, args...)}},
		IsError: true,
	}
}
