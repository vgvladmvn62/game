package server_test

import (
	"testing"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/matching"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/server"
	"github.com/stretchr/testify/assert"
)

type hardwarerMock struct {
	TurnedOnStands []byte
}

func (hw *hardwarerMock) TurnOffLights() error {
	return nil
}

func (hw *hardwarerMock) TurnOnLight(platformID byte) error {
	hw.TurnedOnStands = append(hw.TurnedOnStands, platformID)
	return nil
}

func TestIfMatchedProductsAreSetToBeHighlighted(t *testing.T) {
	//given

	matches := []matching.MatchedProductDTO{{Score: 100, StandID: "1"}, {Score: 100, StandID: "2"}, {Score: 50, StandID: "3"}}
	hw := &hardwarerMock{}

	//when

	err := server.HighlightMatchedProducts(hw, matches)

	//then

	assert.Equal(t, hw.TurnedOnStands, []byte{1, 2})
	assert.Equal(t, nil, err)
}
