package errs

import "net/http"

const (
	// ===== 通用错误 =====
	CodeOK               = 0    // 一切正常
	CodeInternalError    = 1000 // 未知内部错误
	CodeUnauthorized     = 1001 // 未登录或 token 无效
	CodeForbidden        = 1002 // 无权限
	CodeTooManyRequests  = 1003 // 请求过于频繁
	CodeInvalidSignature = 1004 // 签名错误
	CodeNotImplemented   = 1005 // 功能未实现

	// ===== 用户相关 =====
	CodeUserNotFound      = 2000 // 用户不存在
	CodeUserAlreadyExists = 2001 // 用户已存在
	CodeInvalidPassword   = 2002 // 密码错误
	CodeUserDisabled      = 2003 // 用户已禁用
	CodeTokenExpired      = 2004 // Token 过期

	// ===== 资源相关 =====
	CodeResourceNotFound      = 3000 // 资源不存在
	CodeResourceAlreadyExists = 3001 // 资源已存在
	CodeConflict              = 3002 // 操作冲突（如版本冲突）
	CodeResourceLocked        = 3003 // 资源被锁定或占用

	// ===== 参数与校验 =====
	CodeInvalidParam      = 4000 // 参数错误
	CodeMissingParam      = 4001 // 缺少必要参数
	CodeValidationFailed  = 4002 // 参数验证不通过
	CodeUnsupportedFormat = 4003 // 不支持的格式
	CodeOutOfRange        = 4004 // 数值或分页越界

	// ===== 系统与外部依赖 =====
	CodeDatabaseError    = 5000 // 数据库错误
	CodeCacheError       = 5001 // Redis 或缓存错误
	CodeNetworkError     = 5002 // 网络请求失败
	CodeExternalAPIError = 5003 // 调用外部接口失败
	CodeFileSystemError  = 5004 // 文件系统错误
)

func MapToHTTPStatus(code int) int {
	switch code {
	case CodeOK:
		return http.StatusOK
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInvalidSignature, CodeInvalidParam, CodeMissingParam, CodeValidationFailed:
		return http.StatusBadRequest
	case CodeResourceNotFound, CodeUserNotFound:
		return http.StatusNotFound
	case CodeConflict, CodeResourceAlreadyExists:
		return http.StatusConflict
	case CodeNotImplemented:
		return http.StatusNotImplemented
	case CodeDatabaseError, CodeCacheError, CodeNetworkError, CodeExternalAPIError, CodeFileSystemError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
