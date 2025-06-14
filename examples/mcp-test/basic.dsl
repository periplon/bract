# Basic MCP Server Test
# This script tests basic connectivity and tool listing

# Connect to the MCP browser server
connect "./mcp-browser-server"

# List available tools
call list_tools -> tools
print "Available tools:"
print tools

# Verify essential tools exist
assert len(tools) > 0, "No tools available"
print "✓ Server has tools available"