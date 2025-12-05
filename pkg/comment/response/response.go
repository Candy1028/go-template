package response

const (
	// 通用成功码
	OK = 200

	// 通用错误码
	InternalError    = 500
	InvalidParams    = 400
	Unauthorized     = 401
	Forbidden        = 403
	NotFound         = 404
	MethodNotAllowed = 405
	Timeout          = 408
	Conflict         = 409
	TooManyRequests  = 429

	// 业务错误码
	BusinessError     = 600
	ValidationFailed  = 601
	DataNotFound      = 602
	DataConflict      = 603
	OperationFailed   = 604
	InvalidState      = 605
	RateLimitExceeded = 606
	NotInitialized    = 607
	AlreadyExists     = 608

	// 认证授权错误
	InvalidToken       = 701
	TokenExpired       = 702
	InvalidCredentials = 703
	AccessDenied       = 704
	NeedTwoFactorAuth  = 705

	// 文件处理错误
	FileUploadFailed = 801
	FileTooLarge     = 802
	InvalidFileType  = 803
	FileNotFound     = 804

	// 支付相关错误
	PaymentFailed      = 901
	InvalidPayment     = 902
	InsufficientFunds  = 903
	PaymentMethodError = 904

	//用户相关错误
	EmailOrPasswordIncorrect = 1001
	EmailVerifyError         = 1002
	CaptchaInvalid           = 1003
	UserFrozen               = 1004
)

var messages = map[int]string{
	// 通用成功
	OK: "操作成功",

	// 通用错误
	InternalError:    "服务器内部错误",
	InvalidParams:    "参数错误",
	Unauthorized:     "未授权",
	Forbidden:        "禁止访问",
	NotFound:         "资源不存在",
	MethodNotAllowed: "方法不允许",
	Timeout:          "请求超时",
	Conflict:         "资源冲突",
	TooManyRequests:  "请求过于频繁",

	// 业务错误
	BusinessError:     "业务处理失败",
	ValidationFailed:  "参数校验失败",
	DataNotFound:      "数据不存在",
	DataConflict:      "数据冲突",
	OperationFailed:   "操作失败",
	InvalidState:      "状态无效",
	RateLimitExceeded: "超出访问频率限制",
	NotInitialized:    "未初始化",
	AlreadyExists:     "已存在",

	// 认证授权
	InvalidToken:       "无效的令牌",
	TokenExpired:       "令牌已过期",
	InvalidCredentials: "凭证无效",
	AccessDenied:       "访问被拒绝",
	NeedTwoFactorAuth:  "需要双因素认证",

	// 文件处理
	FileUploadFailed: "文件上传失败",
	FileTooLarge:     "文件过大",
	InvalidFileType:  "无效的文件类型",
	FileNotFound:     "文件不存在",

	// 支付相关
	PaymentFailed:      "支付失败",
	InvalidPayment:     "无效的支付请求",
	InsufficientFunds:  "余额不足",
	PaymentMethodError: "支付方式错误",

	EmailOrPasswordIncorrect: "邮箱或密码错误",
	EmailVerifyError:         "邮箱验证失败",
	CaptchaInvalid:           "图形验证码无效",
	UserFrozen:               "用户已被冻结",
}

func GetMessage(code int) string {
	if msg, exists := messages[code]; exists {
		return msg
	}
	return "未知错误"
}
