package test

import "testing"

func True(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Errorf("Value exptected to be True")
	}
}

func False(t *testing.T, b bool) {
	t.Helper()
	if b {
		t.Errorf("Value exptected to be False")
	}
}

func NilErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Exptected err to be nil, got: %v", err)
	}
}
