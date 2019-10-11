package ec

import "encoding/json"

// ProductDTO stores details about the product.
type ProductDTO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Price       Price  `json:"price"`
	Name        string `json:"name"`
	Image       string `json:"image"`
}

type productDTO struct {
	ID          string `json:"code"`
	Description string `json:"description"`
	Price       Price  `json:"price"`
	Name        string `json:"name"`
	Image       Image  `json:"images"`
}

// UnmarshalJSON provides custom unmarshal for ProductDTO.
func (p *ProductDTO) UnmarshalJSON(data []byte) error {
	var productDTO productDTO

	if err := json.Unmarshal(data, &productDTO); err != nil {
		return err
	}

	p.ID = productDTO.ID
	p.Description = productDTO.Description
	p.Price = productDTO.Price
	p.Name = productDTO.Name
	p.Image = productDTO.Image.Image

	return nil
}

// String returns ProductDTO string representation.
func (p *ProductDTO) String() string {
	return "ID: " + p.ID + "\n" +
		"Description: " + p.Description + "\n" +
		"Price: " + "\n" + p.Price.String() + "\n" +
		"Name: " + p.Name + "\n" +
		"Image: " + p.Image
}

// Equals returns true if two products are the same or
// false otherwise.
func (p *ProductDTO) Equals(product ProductDTO) bool {
	if p.Description == product.Description &&
		p.ID == product.ID &&
		p.Name == product.Name &&
		p.Image == product.Image &&
		p.Price == product.Price {
		return false
	}

	return true
}

// IsEmpty checks if product has any description, ID, name and image.
func (p *ProductDTO) IsEmpty() bool {
	if p.Description == "" ||
		p.ID == "" ||
		p.Name == "" ||
		p.Image == "" {
		return true
	}

	return false
}
