package data

import (
	"encoding/json"
	"fmt"
	validator "github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"time"
)

// swagger:model
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" `
	Price       float32 `json:"price" validate:"gt=0" `
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}
type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}
func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if (len(matches)) != 1 {
		return false
	}
	//if fl.Field().String() == "invalid" {
	//	return false
	//}
	return true
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}
func GetProducts() Products {
	return productList
}
func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProductById(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p
	return nil
}
func DeleteProduct(id int) error {
	err := removeProductById(id, &productList)
	if err != nil {
		return err
	}
	return nil
}

var ErrorProductNotFound = fmt.Errorf("Product not found")

func findProductById(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrorProductNotFound
}

func removeProductById(id int, products *[]*Product) error {
	_, pos, err := findProductById(id)
	if err != nil {
		return err
	}
	*products = append((*products)[:pos], (*products)[pos+1:]...)
	return nil
}
func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Tasty Latte coffee",
		Price:       45.44,
		SKU:         "asd123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Americano",
		Description: "Tasty Americano coffee",
		Price:       12.44,
		SKU:         "098asd",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   time.Now().UTC().String(),
	},
}
