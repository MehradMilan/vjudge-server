package judge

/*
#cgo CFLAGS: -I/home/mehrad/Projects/vjudge/vjudge-core
#cgo LDFLAGS: -L/home/mehrad/Projects/vjudge/vjudge-core /home/mehrad/Projects/vjudge/vjudge-core/vjudge.o /home/mehrad/University/TA/402-1/DSD/libvcd/libvcd.o
#include "vjudge.h"
#include "libvcd.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type TestcaseResult struct {
	Name   string
	Passed bool
}
type JudgeResult struct {
	Testcases []TestcaseResult `json:"testcases"`
	score     int
}

// This file can contain any functions needed to implement the judge logic.
func JudgeCode(srcDir string) *JudgeResult {
	println(srcDir)
	input := C.judge_input_t{
		test_dir_path: C.CString("/home/mehrad/Projects/vjudge/vjudge-core/test/testdir"),
		src_dir_path:  C.CString("/home/mehrad/Projects/vjudge/vjudge-core/test/srcdir/"), // Simplifying for one file
	}
	defer C.free(unsafe.Pointer(input.test_dir_path))
	defer C.free(unsafe.Pointer(input.src_dir_path))

	var result C.judge_result_t

	C.run_judge(&input, &result)

	fmt.Println(result.passed_tests_count)
	if result.passed {
		fmt.Println("All tests passed!")
	} else {
		fmt.Printf("%d out of %d tests passed\n", result.passed_tests_count, result.tests_count)
	}
	return &JudgeResult{}
}
