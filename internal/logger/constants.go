package logger

// 日志消息常量
const (
	// 服务器相关
	MsgServerStarting = "Server starting"
	MsgServerStarted  = "Server started successfully"
	MsgServerStopping = "Server stopping"
	MsgServerStopped  = "Server stopped"
	MsgServerError    = "Server error"

	// 数据库相关
	MsgDBConnecting    = "Connecting to database"
	MsgDBConnected     = "Database connected successfully"
	MsgDBDisconnected  = "Database disconnected"
	MsgDBError         = "Database operation error"
	MsgDBQueryStart    = "Database query started"
	MsgDBQueryComplete = "Database query completed"
	MsgDBTransaction   = "Database transaction"

	// 认证相关
	MsgAuthLogin       = "User login attempt"
	MsgAuthLoginSuccess = "User login successful"
	MsgAuthLoginFailed  = "User login failed"
	MsgAuthLogout      = "User logout"
	MsgAuthTokenGenerated = "JWT token generated"
	MsgAuthTokenValidated = "JWT token validated"
	MsgAuthTokenExpired   = "JWT token expired"
	MsgAuthUnauthorized   = "Unauthorized access attempt"

	// 用户相关
	MsgUserCreated    = "User created"
	MsgUserUpdated    = "User updated"
	MsgUserDeleted    = "User deleted"
	MsgUserNotFound   = "User not found"
	MsgUserValidation = "User validation"

	// API相关
	MsgAPIRequest     = "API request received"
	MsgAPIResponse    = "API response sent"
	MsgAPIError       = "API error occurred"
	MsgAPIValidation  = "API validation error"
	MsgAPIRateLimit   = "API rate limit exceeded"

	// 配置相关
	MsgConfigLoaded   = "Configuration loaded"
	MsgConfigError    = "Configuration error"
	MsgConfigMissing  = "Configuration missing"

	// 中间件相关
	MsgMiddlewareStart = "Middleware processing started"
	MsgMiddlewareEnd   = "Middleware processing completed"
	MsgMiddlewareError = "Middleware error"

	// 业务逻辑相关
	MsgBusinessLogic = "Business logic processing"
	MsgValidation    = "Data validation"
	MsgProcessing    = "Data processing"
)

// 字段名常量
const (
	FieldUserID     = "user_id"
	FieldRequestID  = "request_id"
	FieldTraceID    = "trace_id"
	FieldSessionID  = "session_id"
	FieldModule     = "module"
	FieldComponent  = "component"
	FieldOperation  = "operation"
	FieldMethod     = "method"
	FieldPath       = "path"
	FieldStatusCode = "status_code"
	FieldLatency    = "latency"
	FieldClientIP   = "client_ip"
	FieldUserAgent  = "user_agent"
	FieldBodySize   = "body_size"
	FieldError      = "error"
	FieldErrorCode  = "error_code"
	FieldMessage    = "message"
	FieldTimestamp  = "timestamp"
	FieldDuration   = "duration"
	FieldQuery      = "query"
	FieldParams     = "params"
	FieldHeaders    = "headers"
	FieldResponse   = "response"
	FieldRequest    = "request"
)

// 模块名常量
const (
	ModuleAuth       = "auth"
	ModuleUser       = "user"
	ModuleDatabase   = "database"
	ModuleMiddleware = "middleware"
	ModuleConfig     = "config"
	ModuleServer     = "server"
	ModuleAPI        = "api"
	ModuleService    = "service"
	ModuleController = "controller"
	ModuleRepository = "repository"
	ModuleValidator  = "validator"
	ModuleJWT        = "jwt"
	ModuleSecurity   = "security"
)

// 操作名常量
const (
	OpCreate   = "create"
	OpRead     = "read"
	OpUpdate   = "update"
	OpDelete   = "delete"
	OpLogin    = "login"
	OpLogout   = "logout"
	OpValidate = "validate"
	OpConnect  = "connect"
	OpQuery    = "query"
	OpExecute  = "execute"
	OpProcess  = "process"
	OpStart    = "start"
	OpStop     = "stop"
	OpRestart  = "restart"
)

// 错误代码常量
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeAuthentication = "AUTH_ERROR"
	ErrCodeAuthorization = "AUTHZ_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeDatabase      = "DATABASE_ERROR"
	ErrCodeNetwork       = "NETWORK_ERROR"
	ErrCodeTimeout       = "TIMEOUT_ERROR"
	ErrCodeRateLimit     = "RATE_LIMIT_ERROR"
)