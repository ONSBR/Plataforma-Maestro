package actions

import (
	"fmt"
	"testing"

	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetEventsFromListOfInstances(t *testing.T) {
	Convey("should get list of events from instances", t, func() {
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
			events, err := GetEventsFromInstances([]string{"<process_intances>"})
			So(err, ShouldBeNil)
			So(len(events), ShouldEqual, 1)
			event := events[0]
			So(event, ShouldNotBeNil)
			So(event.Name, ShouldEqual, "consolida.saldo.request")
		})
	})

	Convey("should return error on get list of events from instances", t, func() {
		mock := http.ReponseMock{
			Method:        "GET",
			URL:           "*",
			ResponseError: fmt.Errorf("error"),
		}
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			_, err := GetEventsFromInstances([]string{"<process_intances>"})
			So(err, ShouldNotBeNil)
		})
	})
}
