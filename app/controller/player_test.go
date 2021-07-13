package controller

import (
	"../models"
	"gorm.io/gorm/utils/tests"
	"testing"
)

func TestFillDefaultPlayerPayload(t *testing.T) {
	c := Controller{}
	p := models.RandomPlayer(2)
	payload := c.fillDefaultPlayerPayload(p)
	tests.AssertEqual(t, payload.Country, p.Country)
	tests.AssertEqual(t, payload.FirstName, p.FirstName)
	tests.AssertEqual(t, payload.Age, p.Age)
	tests.AssertEqual(t, payload.LastName, p.LastName)
	tests.AssertEqual(t, payload.Position, p.Position)
	tests.AssertEqual(t, payload.MarketValue, p.MarketValue)
}

func TestGetPlayerPayload(t *testing.T) {
	c := Controller{}
	p := models.RandomPlayer(2)
	show := c.getPlayerPayload(p)
	tests.AssertEqual(t, show.Country, p.Country)
	tests.AssertEqual(t, show.FirstName, p.FirstName)
	tests.AssertEqual(t, show.Age, p.Age)
	tests.AssertEqual(t, show.ID, p.ID)
	tests.AssertEqual(t, show.LastName, p.LastName)
	tests.AssertEqual(t, show.Position, p.Position)
	tests.AssertEqual(t, show.MarketValue, p.MarketValue)
}