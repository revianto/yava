package exceptions

import "github.com/gofiber/fiber/v2"

// =============================================================================
// TYPES
// =============================================================================

// AppError represents structured error response
type AppError struct {
	Messages   map[string][]string `json:"messages"`
	Suggestion string              `json:"suggestion"`
	Ext        map[string]string   `json:"ext"`
	Error      string              `json:"error"`
	ErrorCode  int                 `json:"error_code"`
	StatusCode int                 `json:"status_code"`
}

// =============================================================================
// CORS HELPER (Legacy - prefer middleware)
// =============================================================================

// SetCorsHeaders sets CORS headers manually
// NOTE: For production, use fiber/middleware/cors instead
func SetCorsHeaders(c *fiber.Ctx) {
	// c.Set("Access-Control-Allow-Origin", "*")
	// c.Set("Access-Control-Allow-Credentials", "true")
	// c.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

// =============================================================================
// RESPONSE HANDLER
// =============================================================================

// ResponseErrorException sends error response to client
func ResponseErrorException(c *fiber.Ctx, data AppError) error {
	resp := fiber.Map{
		"error":      data.Error,
		"error_code": data.ErrorCode,
		"messages":   []string{},
	}

	if len(data.Messages) > 0 {
		resp["messages"] = data.Messages
	}
	if data.Suggestion != "" {
		resp["suggestion"] = data.Suggestion
	}
	if len(data.Ext) > 0 {
		resp["ext"] = data.Ext
	}

	return c.Status(data.StatusCode).JSON(resp)
}

// =============================================================================
// ERROR FACTORIES
// =============================================================================

// InternalErrorException creates 500 Internal Server Error
func InternalErrorException(c *fiber.Ctx, message string) AppError {
	return AppError{StatusCode: 500, ErrorCode: 500, Error: message}
}

// PageNotFoundErrorException creates 404 Not Found Error
func PageNotFoundErrorException(c *fiber.Ctx, message string) AppError {
	return AppError{StatusCode: 404, ErrorCode: 404, Error: message}
}

// ErrorException creates generic error with custom code
func ErrorException(c *fiber.Ctx, code int, message string) AppError {
	return AppError{StatusCode: 406, ErrorCode: code, Error: message}
}

// ErrorExceptionWithExt creates error with extended data
func ErrorExceptionWithExt(c *fiber.Ctx, code int, message string, ext map[string]string) AppError {
	return AppError{StatusCode: 406, ErrorCode: code, Error: message, Ext: ext}
}

// ValidateException creates validation error
func ValidateException(c *fiber.Ctx, code int, message string) AppError {
	return AppError{StatusCode: 406, ErrorCode: code, Error: message}
}

// ValidateMapException creates validation error with field-level messages
func ValidateMapException(c *fiber.Ctx, code int, messages map[string][]string) AppError {
	return AppError{
		StatusCode: 406,
		ErrorCode:  code,
		Error:      extractFirstError(messages),
		Messages:   messages,
	}
}

// AuthException creates 401 Unauthorized Error
func AuthException(c *fiber.Ctx, code int, message string) AppError {
	return AppError{StatusCode: 401, ErrorCode: code, Error: message}
}

// AuthExceptionWithExt creates 401 error with extended data
func AuthExceptionWithExt(c *fiber.Ctx, code int, message string, ext map[string]string) AppError {
	return AppError{StatusCode: 401, ErrorCode: code, Error: message, Ext: ext}
}

// ThrottleException creates 429 Too Many Requests Error
func ThrottleException(c *fiber.Ctx, code int, message string) AppError {
	return AppError{StatusCode: 429, ErrorCode: code, Error: message}
}

// ErrorWithSuggestionException creates error with suggestion
func ErrorWithSuggestionException(c *fiber.Ctx, code int, message, suggestion string) AppError {
	return AppError{StatusCode: 406, ErrorCode: code, Error: message, Suggestion: suggestion}
}

// =============================================================================
// HELPERS
// =============================================================================

// extractFirstError gets first error message from map
func extractFirstError(messages map[string][]string) string {
	if msgs, ok := messages["0"]; ok && len(msgs) > 0 {
		return msgs[0]
	}
	for _, msgs := range messages {
		if len(msgs) > 0 {
			return msgs[0]
		}
	}
	return ""
}
