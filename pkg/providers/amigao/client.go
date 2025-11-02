package amigao

type AmigaoProvider struct {
}

func NewAmigaoProvider() *AmigaoProvider {
	return &AmigaoProvider{}
}

func (p *AmigaoProvider) FetchProducts() ([]string, error) {
	// Implementation to fetch products from Amigao
	return []string{"Product1", "Product2"}, nil
}
