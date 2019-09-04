package stands

import (
	"log"

	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/ec"
)

// StandDTO stores information about stand's ID, product's ID and
// whether stand is active or not.
type StandDTO struct {
	ID        string `json:"ID"`
	ProductID string `json:"productID"`
	Active    bool   `json:"active"`
}

// Repository holds database connection and product fetching service.
type Repository struct {
	productService productService
	database       *db.Database
}

type productService interface {
	GetProductDetailsByID(ID string) (ec.ProductDTO, error)
}

// NewRepository creates new ShelfRepository.
func NewRepository(data *db.Database) *Repository {
	repo := &Repository{database: data}
	repo.CreateTable()
	return repo
}

// GetAllStands returns all IDs assigned to stands.
func (repository *Repository) GetAllStands() ([]StandDTO, error) {
	var out []StandDTO
	rows, err := repository.database.Query("SELECT id, product_id, active FROM stands")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	activeStand := StandDTO{}

	for rows.Next() {
		err = rows.Scan(&activeStand.ID, &activeStand.ProductID, &activeStand.Active)
		if err != nil {
			return nil, err
		}
		out = append(out, activeStand)
	}

	return out, rows.Err()
}

// GetAllActiveStands returns all IDs assigned to active stands.
func (repository *Repository) GetAllActiveStands() ([]StandDTO, error) {
	var out []StandDTO
	rows, err := repository.database.Query("SELECT id, product_id, active FROM stands WHERE active = TRUE")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	activeStand := StandDTO{}

	for rows.Next() {
		err = rows.Scan(&activeStand.ID, &activeStand.ProductID, &activeStand.Active)
		if err != nil {
			return nil, err
		}
		out = append(out, activeStand)
	}

	return out, rows.Err()
}

// GetAllActiveStandsMap returns all IDs assigned to active stands in form of map,
// assigning stand ID to product ID.
func (repository *Repository) GetAllActiveStandsMap() (map[int]string, error) {
	out := make(map[int]string)
	rows, err := repository.database.Query("SELECT id, product_id FROM stands WHERE active = TRUE")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var stand int
	var product string
	for rows.Next() {
		err = rows.Scan(&stand, &product)
		if err != nil {
			return nil, err
		}
		out[stand] = product
	}

	return out, rows.Err()
}

// GetAllInactiveStandsMap returns all IDs assigned to inactive stands in form map,
// assigning stand ID to product ID.
func (repository *Repository) GetAllInactiveStandsMap() (map[int]string, error) {
	out := make(map[int]string)
	rows, err := repository.database.Query("SELECT id, product_id FROM stands WHERE active = FALSE")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var stand int
	var product string
	for rows.Next() {
		err = rows.Scan(&stand, &product)
		if err != nil {
			return nil, err
		}
		out[stand] = product
	}

	return out, rows.Err()
}

// AddStand assigns stand to product in database.
func (repository *Repository) AddStand(standID string, productID string, active bool) (err error) {
	_, err = repository.database.Exec("INSERT INTO stands(id, product_id, active) VALUES ($1, $2, $3)", standID, productID, active)
	return
}

// DropTable drops tables used by stands repository.
func (repository *Repository) DropTable() (err error) {
	_, err = repository.database.Exec("DROP TABLE stands")
	return
}

// CreateTable creates tables used by stands repository.
func (repository *Repository) CreateTable() (err error) {
	_, err = repository.database.Exec("CREATE TABLE IF NOT EXISTS stands (id INTEGER PRIMARY KEY, product_id TEXT, active BOOLEAN)")
	log.Println("Creating stands ", err)
	return
}
