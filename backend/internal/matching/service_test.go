package matching

import (
	"errors"
	"testing"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/stand"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/attributes"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
	"github.com/stretchr/testify/assert"
)

type shelfFetcherMock struct {
}

func (s *shelfFetcherMock) GetActiveProductsWithData() ([]stand.WithProductDetailsDTO, error) {
	return []stand.WithProductDetailsDTO{{0, ec.ProductDTO{ID: "0"}}, {1, ec.ProductDTO{ID: "1"}},
		{2, ec.ProductDTO{ID: "2"}}, {3, ec.ProductDTO{ID: "3"}}}, nil
}

type shelfFetcherMock2 struct {
}

func (s *shelfFetcherMock2) GetActiveProductsWithData() ([]stand.WithProductDetailsDTO, error) {
	return []stand.WithProductDetailsDTO{{0, ec.ProductDTO{ID: "3"}}, {1, ec.ProductDTO{ID: "0"}},
		{2, ec.ProductDTO{ID: "2"}}}, nil
}

type attributeFetcherMock struct {
}

func (a attributeFetcherMock) GetAttributes(id string) ([]attributes.Attribute, error) {
	switch id {
	case "0":
		return []attributes.Attribute{"one", "two", "something"}, nil
	case "1":
		return []attributes.Attribute{"one", "two"}, nil
	case "2":
		return []attributes.Attribute{"something", "else"}, nil
	case "3":
		return []attributes.Attribute{"completely", "different"}, nil
	default:
		return nil, errors.New("GetAttributes error")
	}
}

func TestIfMatchScoreIsCorrect(t *testing.T) {
	//given

	s := NewService(&shelfFetcherMock{}, attributeFetcherMock{})
	attributes := []attributes.Attribute{"one", "two", "something"}

	//when

	matches, err := s.MatchProducts(attributes)

	//then

	assert.NoError(t, err)

	wanted := MatchedProductDTO{Score: 100, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[0].Score)

	wanted = MatchedProductDTO{Score: 66, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[1].Score)

	wanted = MatchedProductDTO{Score: 33, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[2].Score)

	wanted = MatchedProductDTO{Score: 0, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[3].Score)
}

func TestCaseSensitiveness(t *testing.T) {
	//given

	s := NewService(&shelfFetcherMock{}, attributeFetcherMock{})
	attributes := []attributes.Attribute{"ONE", "TWO", "SOMETHING"}

	//when

	matches, err := s.MatchProducts(attributes)

	//then

	assert.NoError(t, err)

	wanted := MatchedProductDTO{Score: 100, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[0].Score)

	wanted = MatchedProductDTO{Score: 66, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[1].Score)

	wanted = MatchedProductDTO{Score: 33, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[2].Score)

	wanted = MatchedProductDTO{Score: 0, Product: ec.ProductDTO{}}
	assert.Equal(t, wanted.Score, matches[3].Score)
}

func TestIfMatchedProductsAreSortedCorrectly(t *testing.T) {
	//given

	s := NewService(&shelfFetcherMock2{}, attributeFetcherMock{})
	attributes := []attributes.Attribute{"one", "two", "something"}

	//when

	matches, err := s.MatchProducts(attributes)

	//then

	assert.NoError(t, err)

	assert.Equal(t, matches[0].Score, 100)

	assert.Equal(t, matches[1].Score, 33)

	assert.Equal(t, matches[2].Score, 0)
}
