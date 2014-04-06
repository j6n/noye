package store

import (
	"encoding/json"
	"os"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDb(t *testing.T) {
	if err := os.Remove("noye.db"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	Convey("Given a database", t, func() {
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

		Convey("it should set a value", func() {
			err := Set("foo", "bar", string(data))
			So(err, ShouldBeNil)
		})

		Convey("it should get a value", func() {
			res, err := Get("foo", "bar")
			So(err, ShouldBeNil)
			So(res, ShouldResemble, string(data))
		})

		Convey("it should update a value", func() {
			f.Date = time.Now()
			f.Enabled = true

			data, err = json.Marshal(&f)
			So(err, ShouldBeNil)

			err := Set("foo", "bar", string(data))
			So(err, ShouldBeNil)

			res, err := Get("foo", "bar")
			So(err, ShouldBeNil)
			So(res, ShouldResemble, string(data))
		})
	})
}
