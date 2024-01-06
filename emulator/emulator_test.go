package emulator

import (
	"reflect"
	"testing"
)

func TestClearScreen(t *testing.T) {
	screen := Screen{0, 1, 0, 1, 1, 1}
	t.Run("screen is cleared when 00e0 instruction runs", func(t *testing.T) {
		got := ClearScreen(&screen)
		want := &Screen{0, 0, 0, 0, 0, 0}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
