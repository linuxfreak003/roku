package roku_test

import (
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
	"jonwillia.ms/roku"
)

const ROKU_ADDR = "10.70.145.26:8060"

func TestRoku(t *testing.T) {
	G := NewGomegaWithT(t)

	r, err := roku.NewRemote(ROKU_ADDR)
	G.Expect(err).To(BeNil())
	err = r.VolumeUp()
	G.Expect(err).To(BeNil())
}

func TestLaunchWithValues(t *testing.T) {
	r, err := roku.NewRemote(ROKU_ADDR)
	if err != nil {
		panic(err)
	}
	err = r.LaunchWithValues(&roku.App{Id: "63218", Name: "Roku Stream Tester"},
		url.Values{
			"live":          {"true"},
			"autoCookie":    {"true"},
			"debugVideoHud": {"false"},
			"url":           {"https://tv.nknews.org/tvhls/stream.m3u8"},
			"fmt":           {"HLS"},
			"drmParams":     {"{}"},
			"headers":       {`{"Referer":"https://kcnawatch.org/korea-central-tv-livestream/"}`}, // TODO
			"metadata":      {`{"isFullHD":false}`},
			"cookies":       {"[]"},
		},
	)

	if err != nil {
		panic(err)
	}
	r.Home()
}

func TestScan(t *testing.T) {
	devs, err := roku.FindRokuDevices()
	if err != nil {
		panic(err)
	}
	if len(devs) == 0 {
		panic("no devs")
	}
}
