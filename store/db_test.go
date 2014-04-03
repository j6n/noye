package store

import (
	"encoding/json"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDb(t *testing.T) {
	Convey("Given a database", t, func() {
		Convey("it should set a value", func() {
			type foo struct {
				Name    string
				Enabled bool
				What    []string
				Where   map[string]string
				Date    time.Time
			}

			f := &foo{"something", false, []string{"a", "b", "c"}, map[string]string{
				"asdf": "fdsa",
				"baz":  "quux",
			}, time.Now()}

			data, err := json.Marshal(&f)
			So(err, ShouldBeNil)

			err = Set("foo", string(data))
			So(err, ShouldBeNil)

			res, err := Get("foo")
			So(err, ShouldBeNil)
			So(res, ShouldResemble, string(data))
		})
	})
}
