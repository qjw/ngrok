//go:build !release
// +build !release

package client

func useInsecureSkipVerify() bool {
	return true
}
