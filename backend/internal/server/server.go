package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/stands"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/stand"

	"github.com/gorilla/mux"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"

	"github.com/justinas/alice"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/attributes"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/matching"
)

// Server struct holds connectors and settings
type Server struct {
	Router               mux.Router
	middleware           alice.Chain
	cfg                  *Config
	productCacheService  ProductCacheService
	standsRepository     standsRepository
	standService         standService
	matchingService      matchingService
	attributesRepository attributesRepository
	questionsService     questionsRepository
	hardwareService      hardwareService
}

//go:generate mockery -name=ProductCacheService -output=mocks -outpkg=mocks -case=underscore

type standsRepository interface {
	DropTable() error
	CreateTable() error
	AddStand(string, string, bool) error
	GetAllStands() ([]stands.StandDTO, error)
}

type attributesRepository interface {
	DropTable() error
	CreateTable() error
	GetAttributes(id string) ([]attributes.Attribute, error)
	AddAttributes(id string, attrs []attributes.Attribute) error
}

type questionsRepository interface {
	DropTable() error
	CreateTable() error
	AddQuestion(text string, answers []string) error
	GetAllQuestions() ([]QuestionDTO, error)
}

// ProductCacheService fetches products' data from repository.
// If no data is available then appropriate services are called.
// Introduced in order to reduce amount of calls to external services.
type ProductCacheService interface {
	GetProductDetailsByID(ID string) (ec.ProductDTO, error)
	UpdateProducts() error
	ForceUpdateProducts() error
	ForceUpdateProductByID(ID string) error
}

type standService interface {
	GetAllProductsWithData() ([]stand.WithProductDetailsDTO, error)
	GetActiveProductsWithData() ([]stand.WithProductDetailsDTO, error)
}

type matchingService interface {
	MatchProducts(selectedTags []attributes.Attribute) ([]matching.MatchedProductDTO, error)
}

type hardwareService interface {
	TurnOffLights() error
	TurnOnGreenLight(byte) error
}

// NewServer initializes server with needed connectors and default settings.
func NewServer(config *Config, productCacheService ProductCacheService, standsRepository standsRepository, standService standService, matchingService matchingService, attributesRepository attributesRepository, hardwareService hardwareService, questionsRepository questionsRepository) *Server {
	router := *mux.NewRouter()
	middleware := alice.New(
		Header("Access-Control-Allow-Origin", "*"),
		Header("Content-Type", "application/json"),
		Log(config.Logger.Type),
	)

	srv := Server{
		Router:               router,
		middleware:           middleware,
		cfg:                  config,
		productCacheService:  productCacheService,
		standsRepository:     standsRepository,
		attributesRepository: attributesRepository,
		questionsService:     questionsRepository,
		standService:         standService,
		matchingService:      matchingService,
		hardwareService:      hardwareService,
	}

	srv.Router.HandleFunc("/product/{id:[0-9]+}", srv.productGETHandler).Methods("GET")

	srv.Router.HandleFunc("/update", srv.updateGETHandler).Methods("GET")

	srv.Router.HandleFunc("/update/{id:[0-9]+}", srv.updateByIDGETHandler).Methods("GET")

	srv.Router.HandleFunc("/turnofflights", srv.turnOffLightsHandler).Methods("GET")

	srv.Router.HandleFunc("/results", srv.resultsPOSTHandler).Methods("POST", "OPTIONS")

	srv.Router.HandleFunc("/questions", srv.questionsGETHandler).Methods("GET")

	srv.Router.HandleFunc("/questions", srv.questionsPOSTHandler).Methods("POST")

	srv.Router.HandleFunc("/stands", srv.standsPOSTHandler).Methods("POST")

	srv.Router.HandleFunc("/stands", srv.standsGETHandler).Methods("GET")

	srv.Router.HandleFunc("/attributes", srv.attributesPOSTHandler).Methods("POST")

	srv.Router.HandleFunc("/attributes", srv.attributesGETHandler).Methods("GET")

	return &srv
}

// Start runs server.
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%v", s.cfg.IP, s.cfg.Port)
	err := s.productCacheService.UpdateProducts()
	if err != nil {
		log.Println("Products cache update error: ", err)
		return err
	}

	return http.ListenAndServe(address, s.middleware.Then(&s.Router))
}
