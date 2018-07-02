package processmemory

import (
	"fmt"
	"testing"

	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetEventFromInstance(t *testing.T) {
	Convey("should get event from instance", t, func() {
		Convey("should return event from process memory", func() {
			mock := http.ReponseMock{
				Method: "GET",
				URL:    "*",
				ReponseBody: `
				[
					{
						"name": "consolida.saldo.request",
						"instance_id": null,
						"reference_date": null,
						"tag": "e550b199-7be5-11e8-870d-0242ac12000e",
						"scope": "execution",
						"branch": "master",
						"commands": [],
						"reproduction": {},
						"reprocessing": {},
						"payload": {
							"personId": "30696c2d-2ffc-4a2e-97d7-d5140534d3ec"
						}
					}
				]`,
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				event, err := GetEventByInstance("<process_intances>")
				So(err, ShouldBeNil)
				So(event, ShouldNotBeNil)
				So(event.Name, ShouldEqual, "consolida.saldo.request")
			})
		})

		Convey("should not return event from process memory", func() {
			mock := http.ReponseMock{
				Method:      "GET",
				URL:         "*",
				ReponseBody: `[]`,
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				event, err := GetEventByInstance("<process_intances>")
				So(err.Error(), ShouldEqual, "event not found for instance <process_intances>")
				So(event, ShouldBeNil)
			})
		})

		Convey("should return error from process memory", func() {
			mock := http.ReponseMock{
				Method:        "GET",
				URL:           "*",
				ResponseError: fmt.Errorf("error"),
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				_, err := GetEventByInstance("<process_intances>")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("should return error from process memory when not returning json", func() {
			mock := http.ReponseMock{
				Method:      "GET",
				URL:         "*",
				ReponseBody: `hello`,
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				_, err := GetEventByInstance("<process_intances>")
				So(err, ShouldNotBeNil)
			})
		})
	})

}
