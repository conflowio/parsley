package combinator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCombinator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Combinator Suite")
}
