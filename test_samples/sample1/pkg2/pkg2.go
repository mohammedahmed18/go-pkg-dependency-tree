package pkg2

import "runtime"

//

func Subtract(x, y int) int {
	return x - y
}

func ShowPlatform() string {
	return runtime.GOOS
}
