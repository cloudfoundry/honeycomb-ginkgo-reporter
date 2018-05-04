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
			specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
			Expect(specEventArgs).To(Equal(honeycomb.SpecEvent{
				Description: "some-it-description | some-context-description | some-describe-description",
				State: expectedSpecState,
			}))

		},
		Entry("with a successful state", types.SpecStatePassed, "passed"),
		Entry("with a pending state", types.SpecStatePending, "pending"),
		Entry("with a skipped state", types.SpecStateSkipped, "skipped"),
		Entry("with a panicked state", types.SpecStatePanicked, "panicked"),
		Entry("with a timed out state", types.SpecStateTimedOut, "timedOut"),
		Entry("with an invalid state", types.SpecStateInvalid, "invalid"),
	)

	Describe("SpedDidComplete", func(){
		Context("when a spec fails", func(){
			It("tells us the line of code and message of the failure", func(){
				honeycombReporter := honeycomb.New(honeycombClient)
				specSummary := types.SpecSummary{
					State:          types.SpecStateFailed,
					ComponentTexts: []string{"some-it-description", "some-context-description", "some-describe-description"},
					Failure: types.SpecFailure{
						Message: "some-failure-message",
						ComponentCodeLocation: types.CodeLocation{
							FileName: "some-file-name",
							LineNumber: 2,
						},
					},
				}
				honeycombReporter.SpecDidComplete(&specSummary)

				Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
				specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
				Expect(specEventArgs).To(Equal(honeycomb.SpecEvent{
					Description: "some-it-description | some-context-description | some-describe-description",
					State: "failed",
					FailureMessage: "some-failure-message",
					FailureLocation: "some-file-name:2",
				}))
			})
		})
	})

	Describe("SpecSuiteDidEnd", func() {
		It("sends the correct number of flakes to honeycomb", func() {
			honeycombReporter := honeycomb.New(honeycombClient)

			suiteSummary := types.SuiteSummary{
				NumberOfFlakedSpecs: 1,
			}
			honeycombReporter.SpecSuiteDidEnd(&suiteSummary)

			Expect(honeycombClient.SendEventCallCount()).To(Equal(1))

			specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
			Expect(specEventArgs).To(Equal(types.SuiteSummary{
				SuiteDescription: "",
				SuiteSucceeded:   false,
				SuiteID:          "",
				NumberOfSpecsBeforeParallelization: 0,
				NumberOfTotalSpecs:                 0,
				NumberOfSpecsThatWillBeRun:         0,
				NumberOfPendingSpecs:               0,
				NumberOfSkippedSpecs:               0,
				NumberOfPassedSpecs:                0,
				NumberOfFailedSpecs:                0,
				NumberOfFlakedSpecs:                1,
				RunTime:                            0,
			}))
		})

	})

	Describe("with global tags", func() {
		It("should set global tags for each event", func() {
			honeycombReporter := honeycomb.New(honeycombClient)
			globalTags := map[string]interface{}{"some-tag": "some-tag-value"}
			honeycombReporter.SetGlobalTags(globalTags)

			specSummary := types.SpecSummary{
				State:          types.SpecStatePassed,
				ComponentTexts: []string{"some-it-description", "some-context-description", "some-describe-description"},
			}
			honeycombReporter.SpecDidComplete(&specSummary)
			Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
			_, globalTagsArgs, _ := honeycombClient.SendEventArgsForCall(0)
			Expect(globalTagsArgs).To(Equal(globalTags))
		})
	})
	Describe("with custom tags", func() {
		It("should set custom tags for each event", func() {
			honeycombReporter := honeycomb.New(honeycombClient)
			customTags := map[string]interface{}{"some-tag": "some-tag-value"}
			honeycombReporter.SetCustomTags(customTags)

			specSummary := types.SpecSummary{
				State:          types.SpecStatePassed,
				ComponentTexts: []string{"some-it-description", "some-context-description", "some-describe-description"},
			}
			honeycombReporter.SpecDidComplete(&specSummary)
			Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
			_, _, customTagsArgs := honeycombClient.SendEventArgsForCall(0)
			Expect(customTagsArgs).To(Equal(customTags))
		})
	})
})
