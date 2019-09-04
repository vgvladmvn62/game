package attributes

import (
	"database/sql"
	"strings"

	database "github.wdf.sap.corp/Magikarpie/bullseye/internal/db"
)

// Repository holds an instance of a database.
type Repository struct {
	db *database.Database
}

// NewRepository creates new Repository.
func NewRepository(db *database.Database) *Repository {
	return &Repository{db}
}

// Attribute describes product feature.
type Attribute string

// Eq compares two attributes returning true if they are equal
func (attr Attribute) Eq(b Attribute) bool {
	return strings.ToUpper(string(attr)) == strings.ToUpper(string(b))
}

func (repo *Repository) attributeName(id int) (name Attribute, err error) {
	err = repo.db.QueryRow("SELECT attribute FROM attributes WHERE id = $1", id).Scan(&name)
	return
}

// CreateTable creates tables used by this repository
func (repo *Repository) CreateTable() error {
	_, err1 := repo.db.Exec("CREATE TABLE attributes (id SERIAL NOT NULL PRIMARY KEY, attribute TEXT NOT NULL UNIQUE)")
	_, err2 := repo.db.Exec("CREATE TABLE product_attributes (product_id TEXT NOT NULL, attribute_id INTEGER NOT NULL)")

	if err1 == nil {
		return err2
	}
	return err1
}

// DropTable drops tables used by this repository
func (repo *Repository) DropTable() error {
	_, err1 := repo.db.Exec("DROP TABLE attributes")
	_, err2 := repo.db.Exec("DROP TABLE product_attributes")

	if err1 == nil {
		return err2
	}
	return err1
}

// AddAttributes assigns attributes to the product.
func (repo *Repository) AddAttributes(id string, attrs []Attribute) error {
	stmt, err := repo.db.Prepare("INSERT INTO product_attributes(product_id, attribute_id) VALUES ($1, $2);")
	if err != nil {
		return err
	}

	for i := range attrs {
		attrID, err := repo.attrIDOrInsert(attrs[i])
		if err != nil {
			return err
		}

		_, err = stmt.Exec(id, attrID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repository) attrIDOrInsert(name Attribute) (id int, err error) {
	err = repo.db.QueryRow("SELECT id FROM attributes WHERE attribute = $1", name).Scan(&id)
	if err == sql.ErrNoRows {
		err = repo.db.QueryRow("INSERT INTO attributes(id, attribute) VALUES (DEFAULT, $1) RETURNING id", name).Scan(&id)
	}
	return
}

// GetAttributes fetches attributes about the product.
func (repo *Repository) GetAttributes(id string) (attributes []Attribute, err error) {
	rows, err := repo.db.Query("SELECT attribute_id FROM product_attributes WHERE product_id LIKE $1", id)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()
	var attributeIds []int
	var attr int
	for rows.Next() {
		err = rows.Scan(&attr)
		if err != nil {
			return nil, err
		}
		attributeIds = append(attributeIds, attr)
	}

	attributes = make([]Attribute, len(attributeIds))
	for i := range attributeIds {
		attributes[i], err = repo.attributeName(attributeIds[i])
		if err != nil {
			return nil, err
		}
	}

	return attributes, rows.Err()
}
