package tremble_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTremble(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tremble Suite")
}
