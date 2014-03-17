package util

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSlice(t *testing.T) {
	Convey("Slice utils should", t, func() {
		input := []string{"this", "is", "a", "Test", "SLICE"}
		Convey("check to see if it contains an item", func() {
			So(Contains("a", input...), ShouldBeTrue)
			So(Contains("b", input...), ShouldBeFalse)
			So(Contains("slice", input...), ShouldBeFalse)
			So(Contains("SLICE", input...), ShouldBeTrue)
		})
	})
}
