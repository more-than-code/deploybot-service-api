package model

const (
	TaskPending    = "Pending"
	TaskInProgress = "InProgress"
	TaskDone       = "Done"
	TaskCanceled   = "Canceled"
	TaskFailed     = "Failed"
	TaskTimedOut   = "TimedOut"
)

const (
	PipelineIdle = "Idle"
	PipelineBusy = "Busy"
)

const (
	EventBuild  = "Build"
	EventDeploy = "Deploy"
)

const (
	TaskTest    = "Test"
	TaskBuild   = "Build"
	TaskDeploy  = "Deploy"
	TaskRestart = "Restart"
)

const (
	CodeAuthenticationFailure = 100

	CodeInvalidParameters    = 1200
	CodeWrongEmailOrPassword = 1201
	CodeDeletedUser          = 1202

	CodeExpiredVerificationCode              = 1100
	CodeWrongVerificationCode                = 1101
	CodeTooFrequentlySendingVerificationCode = 1102
	CodeNeedingResendingVerificationCode     = 1103
	CodeMaximumAttemptsOnVerificationCode    = 1104
)

const (
	MsgAuthenticationFailure = "authentication failure"

	MsgInvalidParamters     = "invalid parameter(s)"
	MsgWrongEmailOrPassword = "wrong email or password"

	MsgExpiredVerificationCode              = "expired verification code"
	MsgWrongVerificationCode                = "wrong verification code"
	MsgDeletedUser                          = "deleted user"
	MsgTooFrequentlySendingVerificationCode = "too frequently sending code"
	MsgNeedingResendingVerificationCode     = "needing resending verification code"
	MsgMaximumAttemptsOnVerificationCode    = "maximum attempts on verification code"
)
