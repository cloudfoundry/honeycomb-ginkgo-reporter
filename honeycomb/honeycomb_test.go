package honeycomb_test

import (
	"github.com/cloudfoundry/custom-cats-reporters/honeycomb"
	"github.com/cloudfoundry/custom-cats-reporters/honeycomb/client/clientfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/ginkgo/types"
	. "github.com/onsi/gomega"
)

var _ = Describe("Honeycomb Reporter", func() {
	var honeycombClient *clientfakes.FakeClient

	BeforeEach(func() {
		honeycombClient = new(clientfakes.FakeClient)
	})

	DescribeTable("SpecDidComplete",
		func(givenSpecState types.SpecState, expectedSpecState string) {
			honeycombReporter := honeycomb.New(honeycombClient)

			specSummary := types.SpecSummary{
				State:          givenSpecState,
				ComponentTexts: []string{"some-it-description", "some-context-description", "some-describe-description"},
			}
			honeycombReporter.SpecDidComplete(&specSummary)

			Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
			Expect(honeycombClient.SendEventArgsForCall(0)).To(Equal(honeycomb.SpecEvent{Description: "some-it-description | some-context-description | some-describe-description", State: expectedSpecState}))

		},
		Entry("with a successful state", types.SpecStatePassed, "passed"),
		Entry("with a failed state", types.SpecStateFailed, "failed"),
		Entry("with a pending state", types.SpecStatePending, "pending"),
		Entry("with a skipped state", types.SpecStateSkipped, "skipped"),
		Entry("with a panicked state", types.SpecStatePanicked, "panicked"),
		Entry("with a timed out state", types.SpecStateTimedOut, "timedOut"),
		Entry("with an invalid state", types.SpecStateInvalid, "invalid"),
	)
})
