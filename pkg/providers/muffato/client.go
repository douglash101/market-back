package muffato

import (
	"fmt"
	"market/pkg/debug"
	"market/pkg/request"
	"time"

	"github.com/google/uuid"
)

type MuffatoCategoryDump struct {
	From int       `json:"from"`
	To   uuid.UUID `json:"to"`
}

type muffatoProvider struct {
	FetchProductsURL     string
	MuffatoCategoryDumps []MuffatoCategoryDump
}

func NewMuffatoProvider() *muffatoProvider {

	categories := []MuffatoCategoryDump{}
	categories = append(categories,
		MuffatoCategoryDump{From: 8, To: uuid.MustParse("52fb6e4d-d739-49fd-ac00-97215732c79f")},
		MuffatoCategoryDump{From: 53, To: uuid.MustParse("d1ec80e3-95af-41fd-a18e-d4cf056e0a57")},
		MuffatoCategoryDump{From: 72, To: uuid.MustParse("e3ab3538-bec0-4f6a-bcb5-c6d5bfae1440")},
		MuffatoCategoryDump{From: 78, To: uuid.MustParse("636a0d16-a13e-45e4-86cd-3e20a88ca571")},
		MuffatoCategoryDump{From: 99, To: uuid.MustParse("10c4f7aa-b79e-4e81-b96a-eb20b538e37c")},
		MuffatoCategoryDump{From: 134, To: uuid.MustParse("eeb79bb5-881b-49bb-9c5e-f6f9989e9ec9")},
		MuffatoCategoryDump{From: 140, To: uuid.MustParse("af852929-24a9-4a54-9b34-baeef6b3ea75")},
		MuffatoCategoryDump{From: 168, To: uuid.MustParse("9e9be38d-2b0d-482c-ab1a-e3127837b403")},
		MuffatoCategoryDump{From: 181, To: uuid.MustParse("8fe03a58-282b-4261-a447-bfbad503c7c9")},
		MuffatoCategoryDump{From: 186, To: uuid.MustParse("08110ad4-0381-467c-b5fd-7e56e00b5722")},
		MuffatoCategoryDump{From: 551, To: uuid.MustParse("7954b8d2-feaf-43b4-9887-bb011678704b")},
		MuffatoCategoryDump{From: 607, To: uuid.MustParse("960bd9a1-30d9-454b-b0ef-6220dff73085")},
		MuffatoCategoryDump{From: 684, To: uuid.MustParse("37fa473f-f323-4347-9277-6ceba95d5175")},
	)

	return &muffatoProvider{
		FetchProductsURL:     "https://www.supermuffato.com.br/api/catalog_system/pub/products/search/",
		MuffatoCategoryDumps: categories,
	}
}

func (p *muffatoProvider) FetchProducts() error {

	client := request.NewRequest[[]MuffatoProduct](request.RequestParams{
		Name:    "Fetch Muffato Products",
		Method:  request.GET,
		URL:     p.FetchProductsURL,
		Retries: 3,
	})

	allProducts := []MuffatoProduct{}
	for _, category := range p.MuffatoCategoryDumps {
		from := 0
		to := 49
		client.URL = p.FetchProductsURL + "?fq=C:%d&_from=%d&_to=%d"

		for {
			client.URL = fmt.Sprintf(p.FetchProductsURL+"?fq=C:%d&_from=%d&_to=%d", category.From, from, to)
			products, err := client.Execute()

			if err != nil {
				return err
			}

			if len(*products) == 0 {
				break
			}
			allProducts = append(allProducts, *products...)
			from += 50
			to += 50

			time.Sleep(200)
			debug.PrintStructJson(products)
		}
	}

	return nil
}
