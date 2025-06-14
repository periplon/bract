package mcpclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Client wraps the MCP client for testing
type Client struct {
	mcpClient client.MCPClient
}

// Config holds configuration for the MCP client
type Config struct {
	ServerCommand string   // Command to start the server
	ServerArgs    []string // Arguments for the server command
}

// NewClient creates a new MCP client with the given configuration
func NewClient(cfg Config) (*Client, error) {
	return &Client{}, nil
}

// Connect establishes a connection to the MCP server
func (c *Client) Connect(ctx context.Context, serverCmd string, args ...string) error {
	// Create stdio MCP client
	mcpClient, err := client.NewStdioMCPClient(serverCmd, []string{}, args...)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}

	c.mcpClient = mcpClient

	// Initialize the client
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-test-client",
		Version: "1.0.0",
	}

	_, err = c.mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	return nil
}

// CallTool invokes a tool on the MCP server
func (c *Client) CallTool(ctx context.Context, name string, arguments interface{}) (*mcp.CallToolResult, error) {
	if c.mcpClient == nil {
		return nil, fmt.Errorf("client not connected")
	}

	// Convert arguments to JSON if needed
	var args map[string]interface{}
	switch v := arguments.(type) {
	case map[string]interface{}:
		args = v
	case nil:
		args = make(map[string]interface{})
	default:
		// Convert to JSON and back to get a map
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %w", err)
		}
		if err := json.Unmarshal(data, &args); err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}
	}

	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	request.Params.Name = name
	request.Params.Arguments = args

	return c.mcpClient.CallTool(ctx, request)
}

// ListTools returns the list of available tools from the server
func (c *Client) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	if c.mcpClient == nil {
		return nil, fmt.Errorf("client not connected")
	}

	request := mcp.ListToolsRequest{}
	result, err := c.mcpClient.ListTools(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	return result.Tools, nil
}

// Close gracefully shuts down the client
func (c *Client) Close() error {
	if c.mcpClient != nil {
		return c.mcpClient.Close()
	}
	return nil
}
