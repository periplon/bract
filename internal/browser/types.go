package browser

// Tab represents a browser tab
type Tab struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Active  bool   `json:"active"`
	Index   int    `json:"index"`
	Favicon string `json:"favicon,omitempty"`
}

// Cookie represents a browser cookie
type Cookie struct {
	Name           string  `json:"name"`
	Value          string  `json:"value"`
	Domain         string  `json:"domain,omitempty"`
	Path           string  `json:"path,omitempty"`
	Secure         bool    `json:"secure,omitempty"`
	HTTPOnly       bool    `json:"httpOnly,omitempty"`
	SameSite       string  `json:"sameSite,omitempty"`
	ExpirationDate float64 `json:"expirationDate,omitempty"`
}

// Actionable represents an interactive element on a webpage
type Actionable struct {
	LabelNumber int    `json:"labelNumber"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Selector    string `json:"selector"`
}
