package x

import (
	"fmt"
	"runtime"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func GetPackageName() string {
	fmt.Println(MaxInt)
	pc, _, _, _ := runtime.Caller(MaxInt)
	return runtime.FuncForPC(pc).Name()
}
