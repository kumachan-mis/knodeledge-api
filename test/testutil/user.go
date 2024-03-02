package testutil

import "fmt"

func ReadOnlyUserId() string {
	return "auth0|65a3d656ca600978b0f9501b"
}

func ModifyOnlyUserId() string {
	return "auth0|65e28e5aafc0548859b07ef3"
}

func ErrorUserId(i int) string {
	return fmt.Sprintf("error|%024d", i)
}

func UnknownUserId() string {
	return "unknown|000000000000000000000000"
}
