package ddbmock

import "fmt"

// A TestReporter is something that can be used to report test failures.  It
// is satisfied by the standard library's *testing.T.
type TestReporter interface {
	Fatalf(format string, args ...interface{})
}

// mockTestReporter meets the TestReporter interface and is used in tests.
type mockTestReporter struct {
	Logs []string
}

func (m *mockTestReporter) Fatalf(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	m.Logs = append(m.Logs, log)
}
