package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestParsley(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parsley Suite")
}
