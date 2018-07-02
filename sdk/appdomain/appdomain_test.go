package appdomain

import (
	"fmt"
	"testing"

	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetEntitiesOnDomain(t *testing.T) {
	Convey("should get entities from domain", t, func() {

		Convey("should return entities from domain", func() {
			mock := http.ReponseMock{
				Method: "GET",
				URL:    "*",
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
				]
				`,
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				entities, err := GetEntitiesByProcessInstance("<system_id>", "<process_intances>")
				So(err, ShouldBeNil)
				So(len(entities), ShouldEqual, 1)
				So(entities[0]["id"], ShouldEqual, "30696c2d-2ffc-4a2e-97d7-d5140534d3ec")
				So(entities[0]["saldo"], ShouldEqual, 60)
				So(entities[0]["_metadata"].(map[string]interface{})["type"], ShouldEqual, "conta")
			})
		})

		Convey("should not return entities", func() {
			mock := http.ReponseMock{
				Method: "GET",
				URL:    "*",
				ReponseBody: `
				[]
				`,
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				entities, err := GetEntitiesByProcessInstance("<system_id>", "<process_intances>")
				So(err, ShouldBeNil)
				So(len(entities), ShouldEqual, 0)
			})
		})

		Convey("should return error from domain", func() {
			mock := http.ReponseMock{
				Method:        "GET",
				URL:           "*",
				ResponseError: fmt.Errorf("erro"),
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				entities, err := GetEntitiesByProcessInstance("<system_id>", "<process_intances>")
				So(err, ShouldNotBeNil)
				So(entities, ShouldBeNil)
			})
		})
	})
}
