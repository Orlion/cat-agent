package atomicx

import (
	"math"
	"sync/atomic"
	"unsafe"
)

func LoadFloat64(x *float64) float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(unsafe.Pointer(x))))
}

func StoreFloat64(addr *float64, val float64) {
	atomic.StoreUint64((*uint64)(unsafe.Pointer(addr)), math.Float64bits(val))
}
