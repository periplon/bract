# Perix Extension Changes for browser_send_key

The following changes need to be made to the perix Chrome extension in `perix/background.js`:

1. Add to commandHandlers object (around line 11):
```javascript
'tabs.sendKey': handleSendKey,
```

2. Add the handleSendKey function (after handleActivateTab, around line 180):
```javascript
async function handleSendKey(params) {
  const { tabId, key, modifiers = {} } = params;
  
  // Build key event data
  const keyEventData = {
    key: key,
    code: getKeyCode(key),
    ctrlKey: modifiers.ctrl || false,
    shiftKey: modifiers.shift || false,
    altKey: modifiers.alt || false,
    metaKey: modifiers.meta || false,
    bubbles: true,
    cancelable: true
  };

  // Execute script to send key events
  const results = await chrome.scripting.executeScript({
    target: { tabId },
    func: (eventData) => {
      // Get the active element or document body
      const activeElement = document.activeElement || document.body;
      
      // Send keydown event
      const keydownEvent = new KeyboardEvent('keydown', eventData);
      activeElement.dispatchEvent(keydownEvent);
      
      // Send keypress event for printable characters
      if (eventData.key.length === 1) {
        const keypressEvent = new KeyboardEvent('keypress', eventData);
        activeElement.dispatchEvent(keypressEvent);
      }
      
      // Send keyup event
      const keyupEvent = new KeyboardEvent('keyup', eventData);
      activeElement.dispatchEvent(keyupEvent);
      
      return { success: true, activeElement: activeElement.tagName };
    },
    args: [keyEventData]
  });
  
  return results[0].result;
}

// Helper function to convert key names to key codes
function getKeyCode(key) {
  const keyCodes = {
    'Tab': 'Tab',
    'Enter': 'Enter',
    'Escape': 'Escape',
    'Space': 'Space',
    'ArrowUp': 'ArrowUp',
    'ArrowDown': 'ArrowDown',
    'ArrowLeft': 'ArrowLeft',
    'ArrowRight': 'ArrowRight',
    'Backspace': 'Backspace',
    'Delete': 'Delete',
    'Home': 'Home',
    'End': 'End',
    'PageUp': 'PageUp',
    'PageDown': 'PageDown'
  };
  
  // For single characters, use Key + uppercase letter
  if (key.length === 1) {
    return 'Key' + key.toUpperCase();
  }
  
  return keyCodes[key] || key;
}
```

These changes enable the Chrome extension to handle the new `tabs.sendKey` command from the MCP server.