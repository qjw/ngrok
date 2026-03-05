//go:build release
// +build release

package client

func useInsecureSkipVerify() bool {
	// release 不能跳过校验
	return false
}
