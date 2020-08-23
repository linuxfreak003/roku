package roku_test

import (
	"testing"

	"github.com/linuxfreak003/roku"
	. "github.com/onsi/gomega"
)

func TestRoku(t *testing.T) {
	G := NewGomegaWithT(t)

	r, err := roku.NewRemote("192.168.1.51")
	G.Expect(err).To(BeNil())
	err = r.VolumeUp()
	G.Expect(err).To(BeNil())
}
