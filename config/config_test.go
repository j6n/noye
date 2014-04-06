package config

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	conf := &Config{}
	conf.init()

	Convey("Given a configuration", t, func() {
		Convey("it should convert to a k/v map", func() {
			m := conf.ToMap()
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
