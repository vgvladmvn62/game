package matching

import (
	"sort"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/stand"

	"strconv"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/attributes"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
)

type standService interface {
	GetActiveProductsWithData() ([]stand.WithProductDetailsDTO, error)
}

type attributesRepository interface {
	GetAttributes(ID string) ([]attributes.Attribute, error)
}

// Service performs products matching.
type Service struct {
	standService         standService
	attributesRepository attributesRepository
}

// NewService returns new Service.
func NewService(standService standService, attributesRepository attributesRepository) *Service {
	return &Service{standService: standService, attributesRepository: attributesRepository}
}

// MatchedAttributeDTO stores information about selected attributes from the specific product.
type MatchedAttributeDTO struct {
	Attribute attributes.Attribute `json:"attribute"`
	Found     bool                 `json:"found"`
}

// MatchedProductDTO stores information about product.
type MatchedProductDTO struct {
	Score      int                   `json:"score"`
	Product    ec.ProductDTO         `json:"product"`
	Attributes []MatchedAttributeDTO `json:"attributes"`
	StandID    string
}

// MatchProducts select products depending on attributes. SelectedTags must be unique.
func (s *Service) MatchProducts(selectedTags []attributes.Attribute) ([]MatchedProductDTO, error) {
	var matchedProducts []MatchedProductDTO

	stands, err := s.standService.GetActiveProductsWithData()

	if err != nil {
		return nil, err
	}

	for i := range stands {
		matchedProduct, err := s.createMatchedProductFromStand(selectedTags, stands[i])

		if err != nil {
			return nil, err
		}

		matchedProducts = append(matchedProducts, matchedProduct)
	}

	sortMatchedProductsDesc(matchedProducts)

	return matchedProducts, nil
}

func calculateScorePercentage(foundTagCount float64, tagCount float64) int {
	score := int(foundTagCount / tagCount * 100.0)
	return score
}

func calculateTagCounts(selectedTags []attributes.Attribute, tagCount float64, productAttributes []attributes.Attribute, foundTagCount float64) (float64, float64) {
	for _, selectedTag := range selectedTags {
		tagCount++
		foundTagCount = countTagOccurrences(productAttributes, selectedTag, foundTagCount)
	}

	return tagCount, foundTagCount
}

func countTagOccurrences(productAttributes []attributes.Attribute, selectedTag attributes.Attribute, foundTagCount float64) float64 {
	for _, productTag := range productAttributes {
		if selectedTag.Eq(productTag) {
			foundTagCount++
			break
		}
	}

	return foundTagCount
}

func sortMatchedProductsDesc(matchedProducts []MatchedProductDTO) {
	sort.Slice(matchedProducts, func(i, j int) bool {
		return matchedProducts[i].Score > matchedProducts[j].Score
	})
}

func (s *Service) createMatchedProductFromStand(selectedTags []attributes.Attribute, stand stand.WithProductDetailsDTO) (MatchedProductDTO, error) {
	var matchedAttributes []MatchedAttributeDTO
	productAttributes, err := s.attributesRepository.GetAttributes(stand.Product.ID)

	if err != nil {
		return MatchedProductDTO{}, err
	}

	tagCount := 0.0
	foundTagCount := 0.0

	tagCount, foundTagCount = calculateTagCounts(selectedTags, tagCount, productAttributes, foundTagCount)

	matchedAttributes = s.generateMatchedAttributes(productAttributes, selectedTags, matchedAttributes)

	return MatchedProductDTO{Score: calculateScorePercentage(foundTagCount, tagCount), Product: stand.Product, Attributes: matchedAttributes, StandID: strconv.Itoa(stand.ID)}, nil
}

func (s *Service) generateMatchedAttributes(productAttributes []attributes.Attribute, selectedTags []attributes.Attribute, matchedAttributes []MatchedAttributeDTO) []MatchedAttributeDTO {
	for _, productTag := range productAttributes {
		found := false

		for _, selectedTag := range selectedTags {
			if productTag.Eq(selectedTag) {
				found = true
				break
			}
		}

		matchedAttributes = append(matchedAttributes, MatchedAttributeDTO{Attribute: productTag, Found: found})
	}

	return matchedAttributes
}
