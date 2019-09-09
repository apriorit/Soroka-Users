package endpoints

import (
	"testing"
)

func TestCheckPriveledge_ExpectedTrueForAll(t *testing.T) {
	var tests = []struct {
		role   int64
		method string
	}{
		{2048, "changeRole"},
		{16384, "userList"},
		{8192, "userProfile"},
		{4096, "disableUser"},
		{4096, "enableUser"},
	}

	for _, test := range tests {
		if res := CheckPriveledge(test.role, test.method); !res {
			t.Errorf("Mask %d has no priveledge for %s", test.role, test.method)
		}
	}
}

func TestCheckPriveledge_EmptyMask(t *testing.T) {
	var tests = []struct {
		role   int64
		method string
	}{
		{2048, ""},
		{8, ""},
		{64, ""},
		{4096, ""},
		{512, ""},
	}

	for _, test := range tests {
		if res := CheckPriveledge(test.role, test.method); res {
			t.Errorf("Expected false, find %v", res)
		}
	}
}

func TestCheckPriveledge_HasNoPriveledge(t *testing.T) {
	var tests = []struct {
		role   int64
		method string
	}{
		{8, "userList"},
		{64, "userProfile"},
		{512, "enableUser"},
	}

	for _, test := range tests {
		if res := CheckPriveledge(test.role, test.method); res {
			t.Errorf("Expected false, find %v", res)
		}
	}
}
