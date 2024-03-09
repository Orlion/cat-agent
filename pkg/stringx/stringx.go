package stringx

import (
	"fmt"
	"strconv"
)

func B2kbstr(b uint64) string {
	return strconv.Itoa(int(b / 1024))
}

func F642str(b float64) string {
	return fmt.Sprintf("%f", b)
}
