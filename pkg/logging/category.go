package logging

type Category string
type SubCategory string
type ExtraKey string

const (
	General         Category = "General"
	Internal        Category = "Internal"
	Postgres        Category = "Postgres"
	Redis           Category = "Redis"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
)

const (
	// General
	Initialized     SubCategory = "Initialized"
	Startup         SubCategory = "Startup"
	ExternalService SubCategory = "ExternalService"
	Down            SubCategory = "Down"
	//Redis
	Set SubCategory = "Set"
	Get SubCategory = "Get"
	// postgres
	Migration SubCategory = "Migration"
	Select    SubCategory = "Select"
	Rollback  SubCategory = "Rollback"
	Update    SubCategory = "Update"
	Delete    SubCategory = "Delete"
	Insert    SubCategory = "Insert"
	Verify    SubCategory = "Verify"

	// Internal
	Api          SubCategory = "Api"
	HashPassword SubCategory = "HashPassword"

	// Validation
	EmailValidation       SubCategory = "EmailValidation"
	PhonenumberValidation SubCategory = "PhonenumberValidation"
	UsernameValidation    SubCategory = "UsernameValidation"
)

const (
	AppName      ExtraKey = "AppName"
	LoggerName   ExtraKey = "Logger"
	ClientIp     ExtraKey = "ClientIp"
	HostIp       ExtraKey = "HostIp"
	Method       ExtraKey = "Method"
	StatusCode   ExtraKey = "StatusCode"
	BodySize     ExtraKey = "BodySize"
	Path         ExtraKey = "Path"
	Latency      ExtraKey = "Latency"
	RequestBody  ExtraKey = "RequestBody"
	ResponseBody ExtraKey = "ResponseBody"
	ErrorMessage ExtraKey = "ErrorMessage"
)
