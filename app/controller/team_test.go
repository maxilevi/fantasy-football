package controller

import (
	"../models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
	"math/rand"
	"net/http/httptest"
	"testing"
)


func TestValidateTeamOwner(t *testing.T) {
	c := Controller{}
	user1 := models.User{
		Model: gorm.Model{
			ID: 1,
		},
	}
	user2 := models.User{
		Model: gorm.Model{
			ID: 2,
		},
	}
	team, _ := models.RandomTeam()
	team.OwnerID = user1.ID
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	tests.AssertEqual(t, c.validateTeamOwner(ctx, user1, team), true)
	tests.AssertEqual(t, c.validateTeamOwner(ctx, user2, team), false)
}

func TestGetTeamPayload(t *testing.T) {
	c := Controller{}
	p, ps := models.RandomTeam()
	p.ID = uint(rand.Int())
	marketValue := 0
	for _, x := range ps {
		marketValue += int(x.MarketValue)
	}

	show := c.getTeamPayload(p)
	tests.AssertEqual(t, show.Country, p.Country)
	tests.AssertEqual(t, show.ID, p.ID)
	tests.AssertEqual(t, show.Budget, p.Budget)
	tests.AssertEqual(t, show.Name, p.Name)
	tests.AssertEqual(t, len(show.Players), len(ps))
	tests.AssertEqual(t, show.MarketValue, marketValue)
	tests.AssertEqual(t, show.ID, p.ID)
}


func TestFillDefaultTeamPayload(t *testing.T) {
	c := Controller{}
	p, _ := models.RandomTeam()
	p.OwnerID = 100
	update := c.fillDefaultTeamPayload(p)
	tests.AssertEqual(t, update.Country, p.Country)
	tests.AssertEqual(t, update.Budget, p.Budget)
	tests.AssertEqual(t, update.Name, p.Name)
	tests.AssertEqual(t, update.Owner, p.OwnerID)
}
