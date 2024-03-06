// Package bilty API Request/Response Types
// Structs containing only fields which needed to perform specific task
package bilty

type CreateLinkRequest struct {
	Link string `json:"long_url"`
}

type CreateLinkResponse struct {
	ShortLink string `json:"link"`
}

type ErrorMessage struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}
