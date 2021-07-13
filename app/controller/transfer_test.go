package controller

import (
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils/tests"
	"math"
	"net/http/httptest"
	"testing"
)

func TestControllerParseTransferFilters(t *testing.T) {
	c := Controller{}
	params := map[string]interface{}{
		"country": "argentina",
		"team_name": "la seleccion",
		"player_name": "messi",
		"max_age": 38,
		"max_value": 1000000,
		"min_age": 19,
		"min_value": 1000,
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
	tests.AssertEqual(t, filters.MaxAgeFilter, params["max_age"])
	tests.AssertEqual(t, filters.MaxValueFilter, params["max_value"])
	tests.AssertEqual(t, filters.MinAgeFilter, params["min_age"])
	tests.AssertEqual(t, filters.MinValueFilter, params["min_value"])
}

func TestControllerParseEmptyTransferFilters(t *testing.T) {
	c := Controller{}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	filters := c.parseTransferFilters(ctx)

	tests.AssertEqual(t, filters.Country, "")
	tests.AssertEqual(t, filters.TeamName, "")
	tests.AssertEqual(t, filters.PlayerName, "")
	tests.AssertEqual(t, filters.MinAgeFilter, -1)
	tests.AssertEqual(t, filters.MinValueFilter, -1)
	tests.AssertEqual(t, filters.MaxAgeFilter, math.MaxInt32)
	tests.AssertEqual(t, filters.MaxValueFilter, math.MaxInt32)
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
				Age: 26,
			},
		},
		{
			Player:   models.Player{
				FirstName: "timirton",
				Age: 24,
			},
		},
		{
			Player:   models.Player{
				FirstName: "tim",
				Age: 32,
			},
			Ask: 500,
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
		{
			Player:   models.Player{
				FirstName: "tim",
				Age: 55,
			},
		},
		{
			Player:   models.Player{
				FirstName: "tim",
				Age: 30,
			},
			Ask: 5000,
		},
	}
	filter := transferFilters{
		Country:     "",
		TeamName:    "",
		PlayerName:  "ti",
		MinAgeFilter:   24,
		MinValueFilter: -1,
		MaxAgeFilter:   40,
		MaxValueFilter: 1000,
	}

	for _, transfer := range shouldMatch {
		tests.AssertEqual(t, filter.Matches(transfer), true)
	}
	for _, transfer := range shouldNotMatch {
		tests.AssertEqual(t, filter.Matches(transfer), false)
	}
}