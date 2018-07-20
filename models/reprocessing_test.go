package models

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldOrderEventsByTimestamp(t *testing.T) {

	Convey("should order events by timestamp", t, func() {
		list := make([]*domain.Event, 2)
		evt := new(domain.Event)
		evt.Timestamp = "2018-07-20T15:08:41.884268"
		evt.Name = "b"
		evt1 := new(domain.Event)
		evt1.Timestamp = "2018-07-20T15:07:41.884268"
		evt1.Name = "a"
		list[0] = evt
		list[1] = evt1

		Convey("should order events list", func() {
			rep := NewReprocessing(evt)
			rep.Sort(&list)

			So(list[0].Name, ShouldEqual, "a")
			So(list[1].Name, ShouldEqual, "b")
		})

		Convey("should parse event string date", func() {
			time, err := evt.GetTimestamp()
			So(err, ShouldBeNil)
			So(time.Year(), ShouldEqual, 2018)
			So(time.Hour(), ShouldEqual, 15)
		})
	})
}
