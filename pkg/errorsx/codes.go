package errorsx

// ErrorCode 定义应用程序错误代码
type ErrorCode int

// HTTP 状态码相关错误
const (
	// 2xx 成功
	CodeSuccess ErrorCode = 200

	// 4xx 客户端错误
	CodeBadRequest          ErrorCode = 400
	CodeUnauthorized        ErrorCode = 401
	CodeForbidden           ErrorCode = 403
	CodeNotFound            ErrorCode = 404
	CodeMethodNotAllowed    ErrorCode = 405
	CodeConflict            ErrorCode = 409
	CodeUnprocessableEntity ErrorCode = 422
	CodeTooManyRequests     ErrorCode = 429

	// 5xx 服务器错误
	CodeInternalServerError ErrorCode = 500
	CodeServiceUnavailable  ErrorCode = 503
)

// 业务逻辑错误代码 (1000-9999)
const (
	// 用户相关错误 (1000-1999)
	CodeUserNotFound        ErrorCode = 1001
	CodeUserAlreadyExists   ErrorCode = 1002
	CodeInvalidCredentials  ErrorCode = 1003
	CodeUserDisabled        ErrorCode = 1004
	CodePasswordTooWeak     ErrorCode = 1005
	CodeInvalidEmail        ErrorCode = 1006
	CodeInvalidPhone        ErrorCode = 1007

	// 认证授权错误 (2000-2999)
	CodeTokenExpired        ErrorCode = 2001
	CodeTokenInvalid        ErrorCode = 2002
	CodeTokenMissing        ErrorCode = 2003
	CodeInsufficientPermission ErrorCode = 2004

	// 数据验证错误 (3000-3999)
	CodeValidationFailed    ErrorCode = 3001
	CodeRequiredFieldMissing ErrorCode = 3002
	CodeInvalidFormat       ErrorCode = 3003
	CodeValueOutOfRange     ErrorCode = 3004

	// 外部服务错误 (4000-4999)
	CodeDatabaseError       ErrorCode = 4001
	CodeRedisError          ErrorCode = 4002
	CodeThirdPartyAPIError  ErrorCode = 4003

	// 系统错误 (5000-5999)
	CodeConfigError         ErrorCode = 5001
	CodeFileSystemError     ErrorCode = 5002
	CodeNetworkError        ErrorCode = 5003
)

// GetHTTPStatus 获取错误代码对应的HTTP状态码
func (code ErrorCode) GetHTTPStatus() int {
	switch {
	case code >= 200 && code < 300:
		return int(code)
	case code >= 400 && code < 600:
		return int(code)
	case code >= 1000 && code < 2000: // 用户相关错误
		switch code {
		case CodeUserNotFound:
			return 404
		case CodeUserAlreadyExists:
			return 409
		case CodeInvalidCredentials:
			return 401
		case CodeUserDisabled:
			return 403
		default:
			return 400
		}
	case code >= 2000 && code < 3000: // 认证授权错误
		return 401
	case code >= 3000 && code < 4000: // 数据验证错误
		return 400
	case code >= 4000 && code < 5000: // 外部服务错误
		return 500
	case code >= 5000 && code < 6000: // 系统错误
		return 500
	default:
		return 500
	}
}

// GetMessage 获取错误代码的默认消息
func (code ErrorCode) GetMessage() string {
	messages := map[ErrorCode]string{
		// HTTP 状态码
		CodeSuccess:             "Success",
		CodeBadRequest:          "Bad Request",
		CodeUnauthorized:        "Unauthorized",
		CodeForbidden:           "Forbidden",
		CodeNotFound:            "Not Found",
		CodeMethodNotAllowed:    "Method Not Allowed",
		CodeConflict:            "Conflict",
		CodeUnprocessableEntity: "Unprocessable Entity",
		CodeTooManyRequests:     "Too Many Requests",
		CodeInternalServerError: "Internal Server Error",
		CodeServiceUnavailable:  "Service Unavailable",

		// 用户相关错误
		CodeUserNotFound:      "User not found",
		CodeUserAlreadyExists: "User already exists",
		CodeInvalidCredentials: "Invalid credentials",
		CodeUserDisabled:      "User account is disabled",
		CodePasswordTooWeak:   "Password is too weak",
		CodeInvalidEmail:      "Invalid email format",
		CodeInvalidPhone:      "Invalid phone number format",

		// 认证授权错误
		CodeTokenExpired:           "Token has expired",
		CodeTokenInvalid:           "Invalid token",
		CodeTokenMissing:           "Authorization token is missing",
		CodeInsufficientPermission: "Insufficient permission",

		// 数据验证错误
		CodeValidationFailed:     "Validation failed",
		CodeRequiredFieldMissing: "Required field is missing",
		CodeInvalidFormat:        "Invalid format",
		CodeValueOutOfRange:      "Value is out of range",

		// 外部服务错误
		CodeDatabaseError:      "Database error",
		CodeRedisError:         "Redis error",
		CodeThirdPartyAPIError: "Third party API error",

		// 系统错误
		CodeConfigError:     "Configuration error",
		CodeFileSystemError: "File system error",
		CodeNetworkError:    "Network error",
	}

	if msg, exists := messages[code]; exists {
		return msg
	}
	return "Unknown error"
}