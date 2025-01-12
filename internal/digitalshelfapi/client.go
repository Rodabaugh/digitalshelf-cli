package digitalshelfapi

import (
	"net/http"
	"time"
)

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		HttpClient: http.Client{
			Timeout: timeout,
		},
	}
}
