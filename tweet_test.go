package tweethog

import (
	"testing"
)

func TestGetRandomInt(t *testing.T) {
	min := 5
	max := 50

	result := GetRandomInt(min, max)

	if result > max {
		t.Errorf("Random result must not be bigger than %d", max)
	} else if result < min {
		t.Errorf("Random result must not be smaller than %d", min)
	}
}