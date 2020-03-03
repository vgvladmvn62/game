package products

import (
	"database/sql"
	"encoding/json"
	"log"

	//"github.com/lib/pq"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
)

// Repository holds a client of database.
type Repository struct {
	database *db.Database
}

// NewRepository creates new Products Repository.
func NewRepository(db *db.Database) *Repository {
	repo := &Repository{db}

	err := repo.CreateTable()
	if err != nil {
		log.Println(err)
	}

	return repo
}

// CreateTable creates products table in database.
func (repository *Repository) CreateTable() error {
	_, err := repository.database.Exec("CREATE TABLE IF NOT EXISTS products (ID TEXT, data TEXT)")

	return err
}

// DropTable drops products table.
func (repository *Repository) DropTable() error {
	_, err := repository.database.Exec("DROP TABLE products")

	return err
}

// AddProduct adds new product to database.
func (repository *Repository) AddProduct(ID string, data string) error {
	stmt, err := repository.database.Prepare("INSERT INTO products(ID, data) VALUES ($1, $2);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ID, data)
	if err != nil {
		return err
	}

	return nil
}

// GetProductByID returns ProductDTO object with its data from database.
func (repository *Repository) GetProductByID(ID string) (ec.ProductDTO, error) {
	product := ec.ProductDTO{}

	rows, err := repository.database.Query("SELECT ID, data FROM products WHERE ID LIKE $1;", ID)
	if err != nil {
		return ec.ProductDTO{}, err
	}

	defer func() { _ = rows.Close() }()

	var id string
	var data string
	for rows.Next() {
		err = rows.Scan(&id, &data)
		if err != nil {
			return ec.ProductDTO{}, err
		}

		err = json.Unmarshal([]byte(data), &product)
		if err != nil {
			return ec.ProductDTO{}, ec.UnmarshalDataFailedError
		}
	}

	return product, nil
}

// GetAllProducts returns all products from database.
func (repository *Repository) GetAllProducts() ([]ec.ProductDTO, error) {
	var product ec.ProductDTO
	var products []ec.ProductDTO

	rows, err := repository.database.Query("SELECT ID, data FROM products;")
	if err != nil {
		return []ec.ProductDTO{}, err
	}
	defer func() { _ = rows.Close() }()

	var id string
	var data string
	for rows.Next() {
		err = rows.Scan(&id, &data)
		if err != nil {
			return []ec.ProductDTO{}, err
		}

		err = json.Unmarshal([]byte(data), &product)
		if err != nil {
			return []ec.ProductDTO{}, ec.UnmarshalDataFailedError
		}

		products = append(products, product)
	}

	return products, nil
}

// UpdateProductDataByID using information fetched from EC.
func (repository *Repository) UpdateProductDataByID(ID string, newData string) error {
	stmt, err := repository.database.Prepare("UPDATE products SET data = $2 WHERE ID LIKE $1;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ID, newData)
	if err != nil {
		return err
	}

	return nil
}

// RemoveProductByID from products table.
func (repository *Repository) RemoveProductByID(ID string) error {
	stmt, err := repository.database.Prepare("DELETE FROM products WHERE ID LIKE $1;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if table products exists.
func (repository *Repository) Exists() (bool, error) {
	var rows *sql.Rows
	var err error
	rows, err = repository.database.Query("SELECT EXISTS ( SELECT 1 FROM   information_schema.tables WHERE  table_schema = 'public' AND table_name = 'products' );")
	if err != nil {
		return false, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var text string
		err = rows.Scan(&text)
		if err != nil {
			return false, err
		}

		if text == "true" {
			return true, nil
		}
	}

	return false, nil
}
