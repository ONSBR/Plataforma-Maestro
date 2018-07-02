package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetReprocessingInstances(t *testing.T) {
	mock := http.ReponseMock{
		Method: "GET",
		URL:    "http://localhost:9110/instance/<process_instance>/entities",
		ReponseBody: `
		[
			{
				"_metadata": {
					"type": "conta"
				},
				"branch": "master",
				"id": "30696c2d-2ffc-4a2e-97d7-d5140534d3ec",
				"meta_instance_id": "fe292699-1905-4807-9e4b-e0a4c4aa6cbf",
				"saldo": 60
			}
		]`,
	}
	mock2 := http.ReponseMock{
		Method: "GET",
		URL:    "http://localhost:9110/core/installedApp?filter=bySystemIdAndType&systemId=<system_id>&type=domain",
		ReponseBody: `
		[
			{
				"_metadata": {
					"type": "installedApp"
				},
				"host": "localhost",
				"port": 9110
			}
		]
		`,
	}
	Convey("should get reprocessing instances", t, func() {

		Convey("should return error when event has no instance id", func() {
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&mock2)
				evt := new(domain.Event)
				evt.Name = "bla"
				_, err := GetReprocessingInstances(evt)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("should return error when event has no system id", func() {
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&mock2)
				evt := new(domain.Event)
				evt.Name = "bla"
				evt.InstanceID = "<process_instance>"
				_, err := GetReprocessingInstances(evt)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("should return list of process instances", func() {

			a := 1
			a++
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&mock2)
				evt := new(domain.Event)
				evt.Name = "bla"
				evt.InstanceID = "<process_instance>"
				evt.SystemID = "<system_id>"
				list, err := GetReprocessingInstances(evt)
				So(err, ShouldBeNil)
				So(len(list), ShouldBeGreaterThan, 0)
			})
		})

	})
}
