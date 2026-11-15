package main

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewPodExcluder(t *testing.T) {
	testCases := []struct {
		pattern string
		ok      bool
	}{
		{pattern: "", ok: true},
		{pattern: "^awesome-app-web-", ok: true},
		{pattern: "(", ok: false},
	}

	for _, tc := range testCases {
		excludePods, err := NewPodExcluder(tc.pattern)

		if tc.ok {
			if err != nil {
				t.Errorf("NewPodExcluder(%q) returned an unexpected error: %s", tc.pattern, err)
				continue
			}

			if tc.pattern == "" && excludePods != nil {
				t.Errorf("NewPodExcluder(%q) should return nil pattern", tc.pattern)
			}

			continue
		}

		if err == nil {
			t.Errorf("NewPodExcluder(%q) did not return an error", tc.pattern)
		}
	}
}

func TestIsPodExcluded(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  string
		pod      string
		excluded bool
	}{
		{name: "no pattern", pattern: "", pod: "awesome-app-web-1", excluded: false},
		{name: "matched", pattern: "^awesome-app-web-", pod: "awesome-app-web-1", excluded: true},
		{name: "not matched", pattern: "^awesome-app-web-", pod: "awesome-app-worker-1", excluded: false},
		{name: "matched in the middle", pattern: "worker", pod: "awesome-app-worker-1", excluded: true},
	}

	for _, tc := range testCases {
		excludePods, err := NewPodExcluder(tc.pattern)
		if err != nil {
			t.Errorf("%s: NewPodExcluder(%q) returned an unexpected error: %s", tc.name, tc.pattern, err)
			continue
		}

		pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: tc.pod}}

		if got := isPodExcluded(pod, excludePods); got != tc.excluded {
			t.Errorf("%s: isPodExcluded(%q, %q) = %t, want %t", tc.name, tc.pattern, tc.pod, got, tc.excluded)
		}
	}
}
