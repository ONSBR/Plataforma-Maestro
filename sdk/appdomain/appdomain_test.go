package appdomain

import (
	"fmt"
	"testing"

	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetEntitiesOnDomain(t *testing.T) {
	Convey("should get entities from domain", t, func() {
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
		Convey("should return entities from domain", func() {
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
				]
				`,
			}

			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&mock2)
				entities, err := GetEntitiesByProcessInstance("<system_id>", "<process_instance>")
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
				URL:    "http://localhost:9110/instance/<process_instance>/entities",
				ReponseBody: `
				[]
				`,
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&mock2)
				entities, err := GetEntitiesByProcessInstance("<system_id>", "<process_instance>")
				So(err, ShouldBeNil)
				So(len(entities), ShouldEqual, 0)
			})
		})

		Convey("should return error from domain", func() {
			mock := http.ReponseMock{
				Method:        "GET",
				URL:           "http://localhost:9110/instance/<process_instance>/entities",
				ResponseError: fmt.Errorf("erro"),
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&mock2)
				entities, err := GetEntitiesByProcessInstance("<system_id>", "<process_instance>")
				So(err, ShouldNotBeNil)
				So(entities, ShouldBeNil)
			})
		})

		Convey("should return error from apicore", func() {
			mock := http.ReponseMock{
				Method:        "GET",
				URL:           "http://localhost:9110/instance/<process_instance>/entities",
				ResponseError: fmt.Errorf("erro"),
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&http.ReponseMock{
					Method:        mock2.Method,
					URL:           mock2.URL,
					ResponseError: fmt.Errorf("error"),
				})
				_, err := GetEntitiesByProcessInstance("<system_id>", "<process_instance>")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("should return error when app not found on apicore", func() {
			mock := http.ReponseMock{
				Method:        "GET",
				URL:           "http://localhost:9110/instance/<process_instance>/entities",
				ResponseError: fmt.Errorf("erro"),
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&mock)
				ctx.RegisterMock(&http.ReponseMock{
					Method:      mock2.Method,
					URL:         mock2.URL,
					ReponseBody: `[]`,
				})
				_, err := GetEntitiesByProcessInstance("<system_id>", "<process_instance>")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("should return error when message format is not json", func() {
			mock := http.ReponseMock{
				Method:        "GET",
				URL:           "http://localhost:9110/instance/<process_instance>/entities",
				ResponseError: fmt.Errorf("erro"),
			}
			http.With(t, func(ctx *http.MockContext) {
				ctx.RegisterMock(&http.ReponseMock{
					Method:      mock.Method,
					URL:         mock.URL,
					ReponseBody: "asd",
				})
				ctx.RegisterMock(&mock2)
				_, err := GetEntitiesByProcessInstance("<system_id>", "<process_instance>")
				So(err, ShouldNotBeNil)
			})
		})
	})
}
