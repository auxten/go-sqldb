package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtils_Length(t *testing.T) {
	Convey("Length", t, func() {
		a := [16]byte{'a'}
		So(Length(a[:]), ShouldEqual, 1)
		b := [16]byte{}
		So(Length(b[:]), ShouldEqual, 0)
	})
}
