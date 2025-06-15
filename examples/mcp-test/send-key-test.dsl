# Test script for browser_send_key functionality

# Connect to browser extension
wait_for_connection

# Create a test page with input fields
create_tab url:"data:text/html,<html><body><h1>Send Key Test</h1><input id='input1' placeholder='First input'><br><br><input id='input2' placeholder='Second input'><br><br><textarea id='textarea' placeholder='Text area'></textarea><br><br><button id='button'>Test Button</button></body></html>" active:true
wait 2000

# Test 1: Basic key sending - type in first input
click selector:"#input1"
wait 500

# Send individual keys
send_key key:"H"
send_key key:"e"
send_key key:"l"
send_key key:"l"
send_key key:"o"
wait 500

# Test 2: Tab key to move to next input
send_key key:"Tab"
wait 500

# Type in second input
send_key key:"W"
send_key key:"o"
send_key key:"r"
send_key key:"l"
send_key key:"d"
wait 500

# Test 3: Tab to textarea
send_key key:"Tab"
wait 500

# Test special keys
send_key key:"T"
send_key key:"e"
send_key key:"s"
send_key key:"t"
send_key key:"Space"
send_key key:"Enter"
send_key key:"L"
send_key key:"i"
send_key key:"n"
send_key key:"e"
send_key key:"Space"
send_key key:"2"
wait 500

# Test 4: Navigation keys
send_key key:"Home"
wait 200
send_key key:"End"
wait 200
send_key key:"ArrowLeft"
send_key key:"ArrowLeft"
wait 500

# Test 5: Modifier keys - Select all (Ctrl+A or Cmd+A)
send_key key:"a" modifiers:{"ctrl":true}
wait 500

# Delete selected text
send_key key:"Delete"
wait 500

# Type new text
send_key key:"N"
send_key key:"e"
send_key key:"w"
send_key key:"Space"
send_key key:"T"
send_key key:"e"
send_key key:"x"
send_key key:"t"
wait 1000

# Test 6: Tab to button and press Enter
send_key key:"Tab"
wait 500
send_key key:"Enter"
wait 500

# Test 7: Escape key (often used to close dialogs)
send_key key:"Escape"
wait 500

# Extract final values to verify
extract_text selector:"#input1" > input1_value
extract_text selector:"#input2" > input2_value
extract_text selector:"#textarea" > textarea_value

# Print results
print "Input 1 value: ${input1_value}"
print "Input 2 value: ${input2_value}"
print "Textarea value: ${textarea_value}"

# Test 8: Multi-tab test - create another tab and test sendKey with specific tabId
list_tabs > tabs
parse_json input:"${tabs}" > parsed_tabs
set current_tab_id "${parsed_tabs[0].id}"

create_tab url:"data:text/html,<html><body><h1>Tab 2</h1><input id='tab2-input' placeholder='Tab 2 input'></body></html>" active:false > new_tab
parse_json input:"${new_tab}" > parsed_new_tab
set new_tab_id "${parsed_new_tab.id}"

# Send keys to non-active tab
send_key key:"T" tabId:${new_tab_id}
send_key key:"a" tabId:${new_tab_id}
send_key key:"b" tabId:${new_tab_id}
send_key key:"Space" tabId:${new_tab_id}
send_key key:"2" tabId:${new_tab_id}

# Switch to the new tab and verify
activate_tab tabId:${new_tab_id}
wait 500
extract_text selector:"#tab2-input" > tab2_value
print "Tab 2 input value: ${tab2_value}"

# Close the test tabs
close_tab tabId:${new_tab_id}
close_tab tabId:${current_tab_id}

print "Send key tests completed!"