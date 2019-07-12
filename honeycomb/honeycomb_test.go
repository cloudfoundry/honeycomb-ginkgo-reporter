package honeycomb_test

import (
	"time"

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
				State:       expectedSpecState,
			}))

		},
		Entry("with a pending state", types.SpecStatePending, "pending"),
		Entry("with a skipped state", types.SpecStateSkipped, "skipped"),
		Entry("with a panicked state", types.SpecStatePanicked, "panicked"),
		Entry("with a timed out state", types.SpecStateTimedOut, "timedOut"),
		Entry("with an invalid state", types.SpecStateInvalid, "invalid"),
	)

	Describe("SpecDidComplete", func() {
		Context("when a spec fails", func() {
			It("tells us about the failure", func() {
				honeycombReporter := honeycomb.New(honeycombClient)
				specSummary := types.SpecSummary{
					State:          types.SpecStateFailed,
					CapturedOutput: "some-failed-test-output",
					ComponentTexts: []string{"some-it-description", "some-context-description", "some-describe-description"},
					Failure: types.SpecFailure{
						Message: "some-failure-message",
						Location: types.CodeLocation{
							FileName:   "failure-location-file-name",
							LineNumber: 77,
						},
						ComponentCodeLocation: types.CodeLocation{
							FileName:   "component-location-file-name",
							LineNumber: 2,
						},
						ComponentType: types.SpecComponentTypeIt,
					},
					RunTime: 1234 * time.Millisecond,
				}
				honeycombReporter.SpecDidComplete(&specSummary)

				Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
				specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
				Expect(specEventArgs).To(Equal(honeycomb.SpecEvent{
					Description:           "some-it-description | some-context-description | some-describe-description",
					State:                 "failed",
					FailureOutput:         "some-failed-test-output",
					FailureMessage:        "some-failure-message",
					FailureLocation:       "failure-location-file-name:77",
					ComponentCodeLocation: "component-location-file-name:2",
					ComponentType:         "it",
					RunTimeInSeconds:      "1.234000",
				}))
			})
		})

		Context("when a spec passes", func() {
			It("reports RunTime", func() {
				honeycombReporter := honeycomb.New(honeycombClient)
				specSummary := types.SpecSummary{
					State:          types.SpecStatePassed,
					CapturedOutput: "some-test-output",
					ComponentTexts: []string{"some-it-description", "some-context-description", "some-describe-description"},
					RunTime:        1234 * time.Millisecond,
				}
				honeycombReporter.SpecDidComplete(&specSummary)

				Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
				specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
				Expect(specEventArgs).To(Equal(honeycomb.SpecEvent{
					Description:      "some-it-description | some-context-description | some-describe-description",
					State:            "passed",
					RunTimeInSeconds: "1.234000",
				}))
			})
		})
	})

	DescribeTable("BeforeSuiteDidRun",
		func(givenComponentType types.SpecComponentType, expectedType string) {
			honeycombReporter := honeycomb.New(honeycombClient)

			setupSummary := types.SetupSummary{
				ComponentType: givenComponentType,
			}
			honeycombReporter.BeforeSuiteDidRun(&setupSummary)

			Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
			specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
			Expect(specEventArgs).To(Equal(honeycomb.SpecEvent{
				State:         "invalid",
				ComponentType: expectedType,
			}))

		},
		Entry("with an invalid type", types.SpecComponentTypeInvalid, "invalid"),
		Entry("with a container", types.SpecComponentTypeContainer, "container"),
		Entry("with a before suite", types.SpecComponentTypeBeforeSuite, "beforeSuite"),
		Entry("with an after suite", types.SpecComponentTypeAfterSuite, "afterSuite"),
		Entry("with a before each", types.SpecComponentTypeBeforeEach, "beforeEach"),
		Entry("with a just before each", types.SpecComponentTypeJustBeforeEach, "justBeforeEach"),
		Entry("with an after each", types.SpecComponentTypeAfterEach, "afterEach"),
		Entry("with an it", types.SpecComponentTypeIt, "it"),
		Entry("with a measure", types.SpecComponentTypeMeasure, "measure"),
	)

	Describe("BeforeSuiteDidRun", func() {
		Context("when a spec fails", func() {
			It("tells us the component type, location (line of code) of the failure and component and message of the failure", func() {
				honeycombReporter := honeycomb.New(honeycombClient)
				setupSummary := types.SetupSummary{
					State:          types.SpecStateFailed,
					CapturedOutput: "some-failed-test-output",
					ComponentType:  types.SpecComponentTypeBeforeSuite,
					Failure: types.SpecFailure{
						Message: "some-failure-message",
						Location: types.CodeLocation{
							FileName:   "failure-location-file-name",
							LineNumber: 77,
						},
						ComponentCodeLocation: types.CodeLocation{
							FileName:   "component-location-file-name",
							LineNumber: 2,
						},
					},
				}
				honeycombReporter.BeforeSuiteDidRun(&setupSummary)

				Expect(honeycombClient.SendEventCallCount()).To(Equal(1))
				specEventArgs, _, _ := honeycombClient.SendEventArgsForCall(0)
				Expect(specEventArgs).To(Equal(honeycomb.SpecEvent{
					State:                 "failed",
					FailureOutput:         "some-failed-test-output",
					ComponentType:         "beforeSuite",
					FailureMessage:        "some-failure-message",
					FailureLocation:       "failure-location-file-name:77",
					ComponentCodeLocation: "component-location-file-name:2",
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
