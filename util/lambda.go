package util

func IsRunningOnLambda() bool {
	return EnvIf("_LAMBDA_SERVER_PORT", "AWS_LAMBDA_RUNTIME_API", "") != ""
}
