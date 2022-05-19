package ddbmock

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_testreporter.go -package=mocks . TestReporter

// A TestReporter is something that can be used to report test failures.  It
// is satisfied by the standard library's *testing.T.
type TestReporter interface {
	Fatal(args ...interface{})
}
