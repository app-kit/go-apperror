package apperror_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoApperror(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoApperror Suite")
}
