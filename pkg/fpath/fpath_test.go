package fpath

import "testing"

func TestSplitPath(t *testing.T) {
	cases := []struct {
		Input        string
		Want1, Want2 string
	}{
		{"/aaa/bbb.txt", "aaa", "bbb.txt"},
	}

	for _, tc := range cases {
		if got1, got2 := SplitPath(tc.Input); got1 != tc.Want1 || got2 != tc.Want2 {
			t.Errorf("SplitPath=%s %s, want=%s %s", got1, got2, tc.Want1, tc.Want2)
		}
	}

}

func TestSplitName(t *testing.T) {
	cases := []struct {
		Input string
		Want  string
	}{
		{"/aaa", "aaa"},
		{"/aaa/", "aaa"},
	}

	for _, tc := range cases {
		if got := SplitName(tc.Input); got != tc.Want {
			t.Errorf("SplitName=%s, want=%s", got, tc.Want)
		}
	}
}

func TestJoinPath(t *testing.T) {
	cases := []struct {
		Input1, Input2 string
		Want           string
	}{
		{"", "/aaa", "/aaa"},
		{"/aaa", "bbb/ccc.txt", "/aaa/bbb/ccc.txt"},
	}

	for _, tc := range cases {
		if got := JoinPath(tc.Input1, tc.Input2); got != tc.Want {
			t.Errorf("JoinPath=%s, want=%s", got, tc.Want)
		}
	}
}
