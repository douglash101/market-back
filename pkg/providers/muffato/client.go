package muffato

type MuffatoProvider struct {
}

func NewMuffatoProvider() *MuffatoProvider {
	return &MuffatoProvider{}
}

func (p *MuffatoProvider) FetchProducts() ([]string, error) {
	// Implementation to fetch products from Muffato
	return []string{"Product1", "Product2"}, nil
}
