package phone

import (
	"fmt"
	"testing"
)

func TestNormalize(t *testing.T) {
	testCases := []struct {
		nbr       string
		wantedNbr string
	}{
		{"123-123-123", "123123123"},
		{"al;dj24;lkj2409al;sdj//", "242409"},
	}
	fmt.Println("running tests...")

	for _, tc := range testCases {
		newNbr := NormalizeNumber(tc.nbr)
		if tc.wantedNbr != newNbr {
			t.Errorf("%s is not equal %s\n", tc.wantedNbr, newNbr)
		}
	}
}
