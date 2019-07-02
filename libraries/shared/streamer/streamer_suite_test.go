package streamer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStreamer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Streamer Suite")
}
