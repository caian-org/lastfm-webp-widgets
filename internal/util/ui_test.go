package util

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBoolToYesOrNo(t *testing.T) {
	tests := []struct {
		name string
		in   bool
		want string
	}{
		{name: "true returns yes", in: true, want: "yes"},
		{name: "false returns no", in: false, want: "no"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, BoolToYesOrNo(tc.in))
		})
	}
}

func TestUnixToHumanReadable(t *testing.T) {
	// The function uses time.Unix(...,0).Format("02 Jan 2006") which is
	// timezone-sensitive. Avoid asserting fixed strings tied to the local TZ
	// and instead assert the output matches the documented "DD Mon YYYY"
	// shape so the test stays portable across CI and dev machines.
	pattern := regexp.MustCompile(`^\d{2} [A-Z][a-z]{2} \d{4}$`)

	t.Run("epoch zero matches expected format", func(t *testing.T) {
		got := UnixToHumanReadable(0)
		assert.Regexp(t, pattern, got)
	})

	t.Run("arbitrary timestamp matches expected format", func(t *testing.T) {
		got := UnixToHumanReadable(1700000000)
		assert.Regexp(t, pattern, got)
	})
}

func TestDurationToHumanReadable(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		durationType time.Duration
		word         string
		want         string
	}{
		{
			name:         "singular hour",
			duration:     1 * time.Hour,
			durationType: time.Hour,
			word:         "hour",
			want:         "1 hour",
		},
		{
			name:         "plural hours",
			duration:     2 * time.Hour,
			durationType: time.Hour,
			word:         "hour",
			want:         "2 hours",
		},
		{
			name:         "math.Ceil rounds 90 minutes up to 2 hours",
			duration:     90 * time.Minute,
			durationType: time.Hour,
			word:         "hour",
			want:         "2 hours",
		},
		{
			name:         "singular minute",
			duration:     1 * time.Minute,
			durationType: time.Minute,
			word:         "minute",
			want:         "1 minute",
		},
		{
			name:         "plural minutes",
			duration:     5 * time.Minute,
			durationType: time.Minute,
			word:         "minute",
			want:         "5 minutes",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := durationToHumanReadable(tc.duration, tc.durationType, tc.word)
			assert.Equal(t, tc.want, got)
		})
	}
}
