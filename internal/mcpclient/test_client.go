package mcpclient

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// TestClient wraps the MCP client with testing utilities
type TestClient struct {
	*Client
	t *testing.T
}

// NewTestClient creates a new test client
func NewTestClient(t *testing.T) *TestClient {
	client, err := NewClient(Config{})
	require.NoError(t, err)
	
	return &TestClient{
		Client: client,
		t:      t,
	}
}

// ConnectWithTimeout connects to the server with a timeout
func (tc *TestClient) ConnectWithTimeout(serverCmd string, args []string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := tc.Connect(ctx, serverCmd, args...)
	require.NoError(tc.t, err, "failed to connect to server")
}

// MustCallTool calls a tool and fails the test if there's an error
func (tc *TestClient) MustCallTool(ctx context.Context, name string, arguments interface{}) json.RawMessage {
	result, err := tc.CallTool(ctx, name, arguments)
	require.NoError(tc.t, err, "tool call failed for %s", name)
	require.NotNil(tc.t, result, "tool result is nil")
	
	// Extract the content from the result
	if len(result.Content) == 0 {
		return json.RawMessage("{}")
	}
	
	// If there's only one content item, return it directly
	if len(result.Content) == 1 {
		// Check if it's a text content
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			data, err := json.Marshal(textContent.Text)
			require.NoError(tc.t, err)
			return data
		}
		// Otherwise marshal the whole content
		data, err := json.Marshal(result.Content[0])
		require.NoError(tc.t, err)
		return data
	}
	
	// Otherwise return the full content array
	data, err := json.Marshal(result.Content)
	require.NoError(tc.t, err)
	return data
}

// MustListTools lists tools and fails the test if there's an error
func (tc *TestClient) MustListTools(ctx context.Context) []string {
	tools, err := tc.ListTools(ctx)
	require.NoError(tc.t, err, "failed to list tools")
	
	var names []string
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	
	return names
}

// AssertToolExists checks if a tool exists in the server
func (tc *TestClient) AssertToolExists(ctx context.Context, toolName string) {
	tools := tc.MustListTools(ctx)
	require.Contains(tc.t, tools, toolName, "tool %s not found", toolName)
}

// AssertToolResult checks if a tool call returns the expected result
func (tc *TestClient) AssertToolResult(ctx context.Context, toolName string, args interface{}, expected interface{}) {
	result := tc.MustCallTool(ctx, toolName, args)
	
	expectedJSON, err := json.Marshal(expected)
	require.NoError(tc.t, err)
	
	require.JSONEq(tc.t, string(expectedJSON), string(result), 
		"tool %s returned unexpected result", toolName)
}

// Cleanup closes the client and cleans up resources
func (tc *TestClient) Cleanup() {
	if err := tc.Close(); err != nil {
		tc.t.Logf("warning: failed to close test client: %v", err)
	}
}

// TestHarness provides a complete test harness for MCP servers
type TestHarness struct {
	t         *testing.T
	client    *TestClient
	serverCmd string
	serverArgs []string
}

// NewTestHarness creates a new test harness
func NewTestHarness(t *testing.T, serverCmd string, serverArgs ...string) *TestHarness {
	return &TestHarness{
		t:          t,
		serverCmd:  serverCmd,
		serverArgs: serverArgs,
	}
}

// Start starts the test harness
func (th *TestHarness) Start() *TestClient {
	th.client = NewTestClient(th.t)
	th.client.ConnectWithTimeout(th.serverCmd, th.serverArgs, 10*time.Second)
	return th.client
}

// Stop stops the test harness
func (th *TestHarness) Stop() {
	if th.client != nil {
		th.client.Cleanup()
	}
}

// RunTest runs a test function with the test client
func (th *TestHarness) RunTest(testFunc func(*TestClient)) {
	client := th.Start()
	defer th.Stop()
	testFunc(client)
}

// ToolTestCase represents a test case for a tool
type ToolTestCase struct {
	Name     string
	Tool     string
	Args     interface{}
	Expected interface{}
	Error    string
}

// RunToolTests runs a series of tool test cases
func (tc *TestClient) RunToolTests(ctx context.Context, tests []ToolTestCase) {
	for _, test := range tests {
		tc.t.Run(test.Name, func(t *testing.T) {
			if test.Error != "" {
				_, err := tc.CallTool(ctx, test.Tool, test.Args)
				require.Error(t, err)
				require.Contains(t, err.Error(), test.Error)
			} else {
				tc.AssertToolResult(ctx, test.Tool, test.Args, test.Expected)
			}
		})
	}
}

// ExpectError calls a tool expecting an error
func (tc *TestClient) ExpectError(ctx context.Context, toolName string, args interface{}, expectedError string) {
	_, err := tc.CallTool(ctx, toolName, args)
	require.Error(tc.t, err, "expected error for tool %s", toolName)
	require.Contains(tc.t, err.Error(), expectedError, 
		"error message doesn't contain expected text")
}

// WaitForCondition waits for a condition to be true
func (tc *TestClient) WaitForCondition(ctx context.Context, condition func() bool, timeout time.Duration, message string) {
	deadline := time.Now().Add(timeout)
	for !condition() {
		if time.Now().After(deadline) {
			tc.t.Fatalf("timeout waiting for condition: %s", message)
		}
		time.Sleep(100 * time.Millisecond)
	}
}