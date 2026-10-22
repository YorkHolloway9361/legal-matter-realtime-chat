package main

import (
	"testing"
	"time"
)

func TestFollowUpEvent(t *testing.T) {
	now := time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		signed bool
		offset time.Duration
		want   bool
	}{
		{"signed and near", true, 24 * time.Hour, true},
		{"unsigned", false, 24 * time.Hour, false},
		{"too far", true, 72 * time.Hour, false},
		{"past deadline", true, -time.Hour, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, got := FollowUpEvent(FollowUp{MatterID: "m-1", Deadline: now.Add(tc.offset), Signed: tc.signed}, now)
			if got != tc.want {
				t.Fatalf("want %t, got %t", tc.want, got)
			}
		})
	}
}
