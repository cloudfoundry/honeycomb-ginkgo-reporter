package honeycomb_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHoneycomb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Honeycomb Suite")
}
