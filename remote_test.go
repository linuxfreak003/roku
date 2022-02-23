package roku_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"jonwillia.ms/roku"
)

func TestRoku(t *testing.T) {
	G := NewGomegaWithT(t)

	r, err := roku.NewRemote("192.168.1.51")
	G.Expect(err).To(BeNil())
	err = r.VolumeUp()
	G.Expect(err).To(BeNil())
}
