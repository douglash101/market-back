package muffato

type MuffatoProductItemsSellers struct {
	SellerDefault   bool `json:"sellerDefault"`
	CommertialOffer struct {
		Price       float64 `json:"Price"`
		LastPrice   float64 `json:"ListPrice"`
		IsAvailable bool    `json:"IsAvailable"`
	} `json:"commertialOffer"`
}

type MuffatoProductItemsImages struct {
	ImageURL string `json:"imageUrl"`
}

type MuffatoProductItems struct {
	Sellers []MuffatoProductItemsSellers `json:"sellers"`
	Images  []MuffatoProductItemsImages  `json:"images"`
}

func (s *MuffatoProductItems) GetDefaultImageURL() string {
	if len(s.Images) > 0 {
		return s.Images[0].ImageURL
	}
	return ""
}

// "input -> /Carnes, Aves e Peixes/Frango/",
// "output -> Frango"
func (p *MuffatoProduct) GetLastCategory() string {
	if len(p.Categories) > 0 {
		lastCategory := p.Categories[len(p.Categories)-1]
		categoryParts := []rune(lastCategory)
		startIndex := -1
		endIndex := -1
		for i, char := range categoryParts {
			if char == '/' {
				if startIndex == -1 {
					startIndex = i
				} else {
					endIndex = i
				}
			}
		}
		if endIndex != -1 {
			return string(categoryParts[startIndex+1 : endIndex])
		}
	}
	return ""
}

type MuffatoProduct struct {
	ProductID   string                `json:"productId"`
	ProductName string                `json:"productName"`
	Categories  []string              `json:"categories"`
	Items       []MuffatoProductItems `json:"items"`
}
