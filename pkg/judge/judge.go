package judge

/*
#cgo LDFLAGS: -lvcd -lvjudge
#include <vjudge.h>
#include <vcd.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
	"vjudge/pkg/util"
)

type TestcaseResult struct {
	Name   string
	Passed bool
}
type JudgeResult struct {
	Status           ErrorStatus      `json:"status"`
	Passed           bool             `json:"passed"`
	TestsCount       int              `json:"testscount"`
	PassedTestsCount int              `json:"passedtestscount"`
	Testcases        []TestcaseResult `json:"testcases"`
	Score            float64          `json:"score"`
}

type ErrorStatus struct {
	Code    int
	Message string
}

const (
	VJUDGE_NO_ERROR = iota
	VJUDGE_ERROR_OPENING_VCD_FILE
	VJUDGE_ERROR_COMPILING_VERILOG_FILE
	VJUDGE_ERROR_ASSERTIONS_FILE_WRONG_FORMAT
	VJUDGE_ERROR_ASSERTIONS_FILE_NOT_EXISTS
	VJUDGE_ERROR_OPENING_TEST_DIRECTORY
	VJUDGE_ERROR_OPENING_SRC_DIRECTORY
	VJUDGE_ERROR_HANDLING_TEMP_DIRECTORY
)

func getStatus(code int) ErrorStatus {
	const internal_error_message = "An internal error occurred while processing your commit. Please contact the teaching assistants."
	switch code {
	case C.VJUDGE_NO_ERROR:
		return ErrorStatus{VJUDGE_NO_ERROR, "Your code was judged successfully!"}
	case C.VJUDGE_ERROR_COMPILING_VERILOG_FILE:
		return ErrorStatus{VJUDGE_ERROR_COMPILING_VERILOG_FILE, "An error occurred while compiling your Verilog files. Please make sure your code works correctly."}
	case C.VJUDGE_ERROR_OPENING_SRC_DIRECTORY:
		return ErrorStatus{VJUDGE_ERROR_OPENING_SRC_DIRECTORY, internal_error_message + " But first, make sure your code is inside the `src/` directory."}
	case C.VJUDGE_ERROR_OPENING_VCD_FILE:
		return ErrorStatus{VJUDGE_ERROR_OPENING_VCD_FILE, internal_error_message}
	case C.VJUDGE_ERROR_OPENING_TEST_DIRECTORY:
		return ErrorStatus{VJUDGE_ERROR_OPENING_TEST_DIRECTORY, internal_error_message}
	case C.VJUDGE_ERROR_HANDLING_TEMP_DIRECTORY:
		return ErrorStatus{VJUDGE_ERROR_HANDLING_TEMP_DIRECTORY, internal_error_message}
	case C.VJUDGE_ERROR_ASSERTIONS_FILE_WRONG_FORMAT:
		return ErrorStatus{VJUDGE_ERROR_ASSERTIONS_FILE_WRONG_FORMAT, internal_error_message}
	case C.VJUDGE_ERROR_ASSERTIONS_FILE_NOT_EXISTS:
		return ErrorStatus{VJUDGE_ERROR_ASSERTIONS_FILE_NOT_EXISTS, internal_error_message}

	default:
		return ErrorStatus{-1, "Unknown status code"}
	}
}

func isPassed(testsCount int, passedTestsCount int) bool {
	if testsCount > 0 && testsCount == passedTestsCount {
		return true
	}
	return false
}

func getScore(testsCount int, passedTestsCount int) float64 {
	if testsCount > 0 {
		return util.RoundFloat(float64(passedTestsCount)/float64(testsCount)*100, 2)
	}
	return 0
}

func extractTestCases(result C.judge_result_t) []TestcaseResult {
	// Extract data from the C struct
	var testResults []TestcaseResult
	for i := 0; i < int(result.tests_count); i++ {
		// Accessing C struct fields using pointer dereferencing
		// cTest := (*C.test_t)(unsafe.Pointer(uintptr(unsafe.Pointer(result.tests)) + uintptr(i)*unsafe.Sizeof(result.tests[0])))
		cTest := &result.tests[i]
		testName := C.GoString(cTest.name)
		passed := bool(cTest.passed)

		// Creating TestcaseResult struct in Go
		testResult := TestcaseResult{
			Name:   testName,
			Passed: passed,
		}

		// Append the result to the testResults slice
		testResults = append(testResults, testResult)
	}
	return testResults
}

// This file can contain any functions needed to implement the judge logic.
func JudgeCode(srcDir string, testDir string) *JudgeResult {
	println(srcDir)
	input := C.judge_input_t{
		test_dir_path: C.CString(testDir),
		src_dir_path:  C.CString(srcDir), // Simplifying for one file
	}
	defer C.free(unsafe.Pointer(input.test_dir_path))
	defer C.free(unsafe.Pointer(input.src_dir_path))

	var result C.judge_result_t

	C.vjudge_run(&input, &result)

	fmt.Println(result.passed_tests_count)
	if result.passed {
		fmt.Println("All tests passed!")
	} else {
		fmt.Printf("%d out of %d tests passed\n", result.passed_tests_count, result.tests_count)
	}
	return &JudgeResult{
		Status:           getStatus(int(result.error)),
		Passed:           isPassed(int(result.tests_count), int(result.passed_tests_count)),
		TestsCount:       int(result.tests_count),
		PassedTestsCount: int(result.passed_tests_count),
		Score:            getScore(int(result.tests_count), int(result.passed_tests_count)),
		Testcases:        extractTestCases(result),
	}
}
