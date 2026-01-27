package middleware

// All path watch patterns
const (
	PatternGetRoot  string = `^/[^/]*$`
	PatternPostRoot string = `/$`
	//	PatternGetPing             string = `/ping$`
	PatternPostAPIShorten      string = `/api/shorten$`
	PatternPostAPIShortenBatch string = `/api/shorten/batch$`
	PatternGetAPIUserUrls      string = `/api/user/urls$`
	//	PatternDeleteAPIUserUrls   string = `/api/user/urls$`
)
