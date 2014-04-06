package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNoye(t *testing.T) {
	conf := &Config{}
	conf.init()

	Convey("Given a configuration", t, func() {
		Convey("it should convert to a k/v map", func() {
			m := conf.toMap()
			res := map[string]string{
				"Auth":     `["museun"]`,
				"Channels": `["#noye"]`,
				"Nick":     `"noye"`,
				"User":     `"museun"`,
				"Server":   `"localhost:6667"`,
			}
			So(m, ShouldResemble, res)
		})
	})
}
