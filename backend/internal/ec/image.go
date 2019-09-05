package ec

import (
	"encoding/json"
)

// Image stores URL of the picture
type Image struct {
	Image string `json:"image"`
}

type imagesSlice struct {
	Images []imageJSON `json:"images"`
}

type imageJSON struct {
	URL       string `json:"url"`
	ImageType string `json:"imageType"`
}

// UnmarshalJSON provides custom unmarshal for Image
func (i *Image) UnmarshalJSON(data []byte) error {
	var s imagesSlice

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	i.Image = s.getFirstPrimaryImageURL()

	return nil
}

// UnmarshalJSON provides custom unmarshal for imageSlice
func (i *imagesSlice) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &i.Images); err != nil {
		return err
	}

	return nil
}

func (i *imagesSlice) getFirstPrimaryImageURL() string {
	if len(i.Images) > 0 {
		for idx := range i.Images {
			if i.Images[idx].ImageType == "PRIMARY" || i.Images[idx].ImageType == "GALLERY" {
				return i.Images[idx].URL
			}
		}
	}

	return ""
}
