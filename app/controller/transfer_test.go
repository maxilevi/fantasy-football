package controller

import (
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils/tests"
	"net/http/httptest"
	"testing"
)

func TestControllerParseTransferFilters(t *testing.T) {
	c := Controller{}
	params := map[string]interface{}{
		"country": "argentina",
		"team_name": "la seleccion",
		"player_name": "messi",
		"age": 38,
		"value": 1000000,
	}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	var p []gin.Param
	for k, v := range params {
		p = append(p, gin.Param{
			Key: k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	ctx.Params = p
	filters := c.parseTransferFilters(ctx)

	tests.AssertEqual(t, filters.Country, params["country"])
	tests.AssertEqual(t, filters.TeamName, params["team_name"])
	tests.AssertEqual(t, filters.PlayerName, params["player_name"])
	tests.AssertEqual(t, filters.AgeFilter, params["age"])
	tests.AssertEqual(t, filters.ValueFilter, params["value"])
}

func TestControllerParseEmptyTransferFilters(t *testing.T) {
	c := Controller{}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	filters := c.parseTransferFilters(ctx)

	tests.AssertEqual(t, filters.Country, "")
	tests.AssertEqual(t, filters.TeamName, "")
	tests.AssertEqual(t, filters.PlayerName, "")
	tests.AssertEqual(t, filters.AgeFilter, -1)
	tests.AssertEqual(t, filters.ValueFilter, -1)
}

func TestEmptyFiltersMatchEverything(t *testing.T) {
	shouldMatch := []models.Transfer{
		{
			Player: models.Player{
				FirstName: "tito",
			},
		},
		{
			Player:   models.Player{
				LastName: "messi",
			},
		},
		{
			Player:   models.Player{
				LastName: "messi",
				Age: 0,
			},
		},
		{
			Player:   models.Player{
				Age: 123,
			},
		},
		{
			Player:   models.Player{
				Country: "bolivia",
			},
		},
	}
	c := Controller{}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	filters := c.parseTransferFilters(ctx)

	for _, transfer := range shouldMatch {
		tests.AssertEqual(t, filters.Matches(transfer), true)
	}
}

func TestControllerTransferFilterMatches(t *testing.T) {
	shouldMatch := []models.Transfer{
		{
			Player:   models.Player{
				FirstName: "tito",
			},
		},
		{
			Player:   models.Player{
				FirstName: "timirton",
				Age: 24,
			},
		},
	}
	shouldNotMatch := []models.Transfer{
		{
			Player:   models.Player{
				FirstName: "tobias",
			},
		},
		{
			Player:   models.Player{
				FirstName: "tim",
				Age: 18,
			},
		},
	}
	filter := transferFilters{
		Country:     "",
		TeamName:    "",
		PlayerName:  "ti",
		AgeFilter:   24,
		ValueFilter: -1,
	}

	for _, transfer := range shouldMatch {
		tests.AssertEqual(t, filter.Matches(transfer), true)
	}
	for _, transfer := range shouldNotMatch {
		tests.AssertEqual(t, filter.Matches(transfer), false)
	}
}