package middleware

// All path watch patterns
const (
	PatternGetRoot             string = `/(?!.*/)`
	PatternPostRoot            string = `/$`
	PatternGetPing             string = `/ping$`
	PatternPostApiShorten      string = `/api/shorten$`
	PatternPostApiShortenBatch string = `/api/shorten/batch$`
	PatternGetApiUserUrls      string = `/api/user/urls$`
	PatternDeleteApiUserUrls   string = `/api/user/urls$`
)
