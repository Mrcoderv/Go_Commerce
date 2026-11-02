package incoming

type Product interface{}

type Product_Manipulation interface {
	New(product Product) error
	GetByID(id int64) (Product, error)
	GetAll([]Product, error)
	Update(product Product) error
	Delete(id int64) error
}
