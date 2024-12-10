package util

import "os"

// EnvIf Looks up first non-empty environment variable from the arguments, or returns the last
// argument if all are empty
func EnvIf(args ...string) string {
	lastIndex := -1 + len(args)

	result := args[lastIndex]

	for _, envVarName := range args[0:lastIndex] {
		if newValue, ok := os.LookupEnv(envVarName); ok {
			return newValue
		}
	}

	return result
}
