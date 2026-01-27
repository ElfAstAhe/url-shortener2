package dto

// UserShorten is user shorten url
type UserShorten struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewUserShorten(shortURL string, originalURL string) *UserShorten {
	return &UserShorten{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
