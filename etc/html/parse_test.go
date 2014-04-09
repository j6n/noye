package html

import (
	"github.com/robertkrimen/otto"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestParser(t *testing.T) {
	ctx := otto.New()

	Convey("Given a parser", t, func() {
		parser, err := NewParser("http://music.naver.com/promotion/specialContent.nhn?articleId=4569", ctx)
		So(err, ShouldBeNil)
		So(parser, ShouldNotBeNil)

		parser.doc.Find("li[id^='videoContent']").
			Each(parser.get("id"))

		for k, v := range parser.results {
			t.Logf("> %s: %v\n", k, v)
		}
	})
}
