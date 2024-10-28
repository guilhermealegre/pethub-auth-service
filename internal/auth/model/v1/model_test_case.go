package v1

import (
	"github.com/guilhermealegre/go-clean-arch-core-lib/test"
	_ "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

type TestCase struct {
	test.BaseTestCase
	Streaming  test.MapCall
	Repository test.MapCall
}

func testCaseScannerLoginWithSuccess() *TestCase {

	return &TestCase{}
}

func testCaseScannerLoginWithError() *TestCase {

	return &TestCase{}
}
