package models

type EmailQuote struct {
	Body    *string `json:"Body,omitempty"`
	Subject *string `json:"Subject,omitempty"`
	Url     *string `json:"url,omitempty"`
}
