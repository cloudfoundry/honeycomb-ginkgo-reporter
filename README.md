# Honeycomb Ginkgo Reporter

## Getting Started
You should be able to get set up in less than 30 minutes.
If it takes you longer, 
please open an issue,
and I'll help you out
and then update this README.

To use this reporter, 
you'll need a Honeycomb account.
You should be able to [make an account on Honeycomb.io](https://ui.honeycomb.io/signup)
and get a dataset and a write key following their documentation.

Then, in your suite test file,
import this reporter and the honeycomb go library

```
"github.com/cloudfoundry/honeycomb-ginkgo-reporter/honeycomb"
"github.com/cloudfoundry/honeycomb-ginkgo-reporter/honeycomb/client"
"github.com/honeycombio/libhoney-go"
```

create a Honeycomb client
with the write key and the data set name,
and pass use that client to set up.

```
honeyCombClient := client.New(libhoney.Config{
	WriteKey: "write-key",
	Dataset:  "my-dataset",
})
honeyCombReporter := honeycomb.New(honeyCombClient)
```
Then run your Ginkgo specs with the custom reporter.

```
rs := []Reporter{}
rs = append(rs, honeyCombReporter)

RunSpecsWithDefaultAndCustomReporters(t, "suite-name", rs)
```
(This will replace the default `RunSpecs` function.)

## Usage Example
We use this reporter to track flakes in the Cloud Foundry Acceptance Tests (CATs).

Here's [the suite test code from that project.](https://github.com/cloudfoundry/cf-acceptance-tests/blob/master/cats_suite_test.go#L145)

The [cf-deployment Concourse task for running CATs](https://github.com/cloudfoundry/cf-deployment-concourse-tasks/blob/master/run-cats/task#L18)
contains an example of passing in configuration,
and adds some additional fields at runtime.

## Running Queries
To get started, trying breaking down by `Description`
Calculate `COUNT` per group
add a Filter for `state = failed`
and a Filter for `Description != ''`
and order by `COUNT desc`.

This will generate a list of tests
with the most failingest at the top.

Then, try filtering by that tests's description,
and break down by `ComponentCodeDescription`
to understand exactly how that test fails.

## Adding Moar Data

You'll probably want to add some tags to your events,
so you can break down your test failures by independent variables,
and figure out what's causing your test failures.

You may also want to filter out tests from specific runs,
from specific environments,
or tests that were run under specific conditions.

You can add tags in two sets: Global and Custom.
```
globalTags := map[string]interface{}{
	"run_id":  os.Getenv("RUN_ID"),
	"env_api": Config.GetApiEndpoint(),
}

honeyCombReporter := honeycomb.New(honeyCombClient)
honeyCombReporter.SetGlobalTags(globalTags)
honeyCombReporter.SetCustomTags(reporterConfig.CustomTags)
```
 
These are two somewhat arbitrary buckets.
In CATs, "global" tags are tags that the test suite expects to be present.
"Custom" tags are passed in ad-hoc.

I recommend at least setting unique `run_id`
one or more tags describing the environment the test was run in,
and as much information about the version of the code under test as you have.
