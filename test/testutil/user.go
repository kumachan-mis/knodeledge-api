package testutil

import "fmt"

func UserId() string {
	return "auth0|65a3d656ca600978b0f9501b"
}

func ErrorUserId(i int) string {
	return fmt.Sprintf("error|%024d", i)
}

func UnknownUserId() string {
	return "unknown|000000000000000000000000"
}
