package honeycomb

import (
	"github.com/cloudfoundry/custom-cats-reporters/honeycomb/client"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
	"strings"
)

type SpecEvent struct {
	Description string
	State       string
}

type honeyCombReporter struct {
	client client.Client
}

func New(client client.Client) honeyCombReporter {

	return honeyCombReporter{client: client}
}

func (hr honeyCombReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	specEvent := SpecEvent{}

	specEvent.State = getTestState(specSummary.State)
	specEvent.Description = createTestDescription(specSummary.ComponentTexts)

	hr.client.SendEvent(specEvent)
}

func (hr honeyCombReporter) SpecSuiteWillBegin(config config.GinkgoConfigType, summary *types.SuiteSummary) {
}
func (hr honeyCombReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {}
func (hr honeyCombReporter) SpecWillRun(specSummary *types.SpecSummary)         {}
func (hr honeyCombReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary)  {}
func (hr honeyCombReporter) SpecSuiteDidEnd(summary *types.SuiteSummary)        {}

func getTestState(state types.SpecState) string {
	switch state {
	case types.SpecStatePassed:
		return "passed"
	case types.SpecStateFailed:
		return "failed"
	case types.SpecStatePending:
		return "pending"
	case types.SpecStateSkipped:
		return "skipped"
	case types.SpecStatePanicked:
		return "panicked"
	case types.SpecStateTimedOut:
		return "timedOut"
	case types.SpecStateInvalid:
		return "invalid"
	default:
		panic("unkown spec state")
	}
}

func createTestDescription(componentTexts []string) string {
	return strings.Join(componentTexts, " | ")
}
