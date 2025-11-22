package utils

func StatusCode1xxFamily(StatusCode int) bool {
	return StatusCode >= 100 && StatusCode < 200
}

func IsInformational(StatusCode int) bool {
	return StatusCode1xxFamily(StatusCode)
}

func StatusCode2xxFamily(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func IsSuccess(StatusCode int) bool {
	return StatusCode2xxFamily(StatusCode)
}

func StatusCode3xxFamily(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

func IsRedirection(StatusCode int) bool {
	return StatusCode3xxFamily(StatusCode)
}

func StatusCode4xxFamily(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func IsClientError(StatusCode int) bool {
	return StatusCode4xxFamily(StatusCode)
}

func StatusCode5xxFamily(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

func IsServerError(StatusCode int) bool {
	return StatusCode5xxFamily(StatusCode)
}
