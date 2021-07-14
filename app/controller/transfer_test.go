package controller

import (
	"../models"
	"fmt"
	"gorm.io/gorm/utils/tests"
	"math"
	"net/url"
	"testing"
)

func TestControllerParseTransferFilters(t *testing.T) {
	c := Controller{}
	params := map[string]interface{}{
		"country":     "argentina",
		"team_name":   "la seleccion",
		"player_name": "messi",
		"max_age":     38,
		"max_value":   1000000,
		"min_age":     19,
		"min_value":   1000,
		"value_type":  "market",
	}
	p := url.Values{}
	for k, v := range params {
		p.Add(k, fmt.Sprintf("%v", v))
	}
	filters := c.parseTransferFilters(p)

	tests.AssertEqual(t, filters.Country, params["country"])
	tests.AssertEqual(t, filters.TeamName, params["team_name"])
	tests.AssertEqual(t, filters.PlayerName, params["player_name"])
	tests.AssertEqual(t, filters.MaxAgeFilter, params["max_age"])
	tests.AssertEqual(t, filters.MaxValueFilter, params["max_value"])
	tests.AssertEqual(t, filters.MinAgeFilter, params["min_age"])
	tests.AssertEqual(t, filters.MinValueFilter, params["min_value"])
	tests.AssertEqual(t, filters.ValueType, params["value_type"])
}

func TestControllerParseEmptyTransferFilters(t *testing.T) {
	c := Controller{}
	filters := c.parseTransferFilters(url.Values{})

	tests.AssertEqual(t, filters.Country, "")
	tests.AssertEqual(t, filters.TeamName, "")
	tests.AssertEqual(t, filters.PlayerName, "")
	tests.AssertEqual(t, filters.MinAgeFilter, -1)
	tests.AssertEqual(t, filters.MinValueFilter, -1)
	tests.AssertEqual(t, filters.MaxAgeFilter, math.MaxInt32)
	tests.AssertEqual(t, filters.MaxValueFilter, math.MaxInt32)
	tests.AssertEqual(t, filters.ValueType, "")
}

func TestEmptyFiltersMatchEverything(t *testing.T) {
	shouldMatch := []models.Transfer{
		{
			Player: models.Player{
				FirstName: "tito",
			},
		},
		{
			Player: models.Player{
				LastName: "messi",
			},
		},
		{
			Player: models.Player{
				LastName: "messi",
				Age:      0,
			},
		},
		{
			Player: models.Player{
				Age: 123,
			},
		},
		{
			Player: models.Player{
				Country: "bolivia",
			},
		},
	}
	c := Controller{}
	filters := c.parseTransferFilters(url.Values{})

	for _, transfer := range shouldMatch {
		tests.AssertEqual(t, filters.Matches(transfer), true)
	}
}

func TestControllerTransferFilterMatches(t *testing.T) {
	shouldMatch := []models.Transfer{
		{
			Player: models.Player{
				FirstName: "tito",
				Age:       26,
			},
		},
		{
			Player: models.Player{
				FirstName: "timirton",
				Age:       24,
			},
		},
		{
			Player: models.Player{
				FirstName: "tim",
				Age:       32,
			},
			Ask: 500,
		},
	}
	shouldNotMatch := []models.Transfer{
		{
			Player: models.Player{
				FirstName: "tobias",
			},
		},
		{
			Player: models.Player{
				FirstName: "tim",
				Age:       18,
			},
		},
		{
			Player: models.Player{
				FirstName: "tim",
				Age:       55,
			},
		},
		{
			Player: models.Player{
				FirstName: "tim",
				Age:       30,
			},
			Ask: 5000,
		},
	}
	filter := transferFilters{
		Country:        "",
		TeamName:       "",
		PlayerName:     "ti",
		MinAgeFilter:   24,
		MinValueFilter: -1,
		MaxAgeFilter:   40,
		MaxValueFilter: 1000,
		ValueType:      "ask",
	}

	for _, transfer := range shouldMatch {
		tests.AssertEqual(t, filter.Matches(transfer), true)
	}
	for _, transfer := range shouldNotMatch {
		tests.AssertEqual(t, filter.Matches(transfer), false)
	}
}
