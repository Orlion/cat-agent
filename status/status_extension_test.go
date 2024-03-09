package status

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestAgentRuntimeGcExtension(t *testing.T) {
	ext := newAgentRuntimeGcExtension()

	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	runtime.GC()
	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	fmt.Println(len(allocForTest()))
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	fmt.Println(len(allocForTest()))
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	fmt.Println(len(allocForTest()))
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(ext.GetProperties())
}

func TestAgentRuntimeMemExtension(t *testing.T) {
	ext := newAgentRuntimeMemExtension()

	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	runtime.GC()
	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	fmt.Println(len(allocForTest()))
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	fmt.Println(len(allocForTest()))
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(ext.GetProperties())
	fmt.Println(len(allocForTest()))
	fmt.Println(len(allocForTest()))
	time.Sleep(time.Second)
	runtime.GC()
	t.Log(ext.GetProperties())
	t.Log(ext.GetProperties())
}

func allocForTest() []byte {
	b := make([]byte, 1024)
	return b
}
