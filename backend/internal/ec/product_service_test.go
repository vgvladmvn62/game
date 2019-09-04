package ec

import (
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
)

func prepareConfig(URL string) *Config {
	var config = Config{}
	config.Host.API = URL
	config.Products.Site = "electronics"
	config.HTTP.Client.Timeout.Seconds = 10
	return &config
}

func prepareProductService(URL string) *ProductService {
	config := prepareConfig(URL)

	return &ProductService{
		Cfg:  config,
		Doer: NewHTTPClient(config.HTTP.Client.Timeout.Seconds),
	}
}

func TestGetExistingProduct(t *testing.T) {
	//given

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		file, err := ioutil.ReadFile("../test_data/occ_correct_get_response.json")
		assert.NoError(t, err)
		w.Write(file)
	}))
	defer ts.Close()

	p := prepareProductService(ts.URL)

	//when

	got, err := p.GetProductDetailsByID("1382080")

	//then

	assert.NoError(t, err)

	wanted := ProductDTO{
		"1382080",
		"The EOS 450D blends uncompromising performance with ease of use in a lightweight, ergonomic body. Capture your world with a 12.2 Megapixel CMOS sensor and the very latest Canon technologies.<br/>Features<br/><br/>- 12.2 MP CMOS sensor<br/>- 3.5fps<br/>- 9-point wide-area AF<br/>- EOS Integrated Cleaning System<br/>- 3.0” LCD with Live View mode<br/>- DIGIC III processor<br/>- Large, bright viewfinder<br/>- Total image control<br/>- Compact and lightweight<br/>- Compatible with EF/EF-S lenses and EX Speedlites",
		Price{
			"USD",
			57488,
			"$574.88",
		},
		"EOS450D + 18-55 IS Kit",
		p.Cfg.Host.Static + "/rest/v2/medias/?context=bWFzdGVyfGltYWdlc3wyNDM2MXxpbWFnZS9qcGVnfGltYWdlcy9oMDcvaGVlLzg3OTY4MjMzMjI2NTQuanBnfDY5NGU2NjU1MDZmMTY5ZGZhZmUxZDhlN2IzYmU2N2NjYjMzNjdmNmNjYjZhZWJhMzkwNjU4YzE0MWYxYWNmMmY",
	}

	assert.Equal(t, wanted, got)
}

func TestGetProductThatDoesNotExistCorrectEndpoint(t *testing.T) {
	//given

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":[{"message":"Product with code '9999999999999' not found!","type":"UnknownIdentifierError"}]}`))
	}))
	defer ts.Close()

	p := prepareProductService(ts.URL)

	//when

	got, err := p.GetProductDetailsByID(fmt.Sprintf("%d", math.MaxInt64))

	//then

	assert.Error(t, err)

	wanted := ProductDTO{}

	assert.Equal(t, wanted, got)
}

func TestGetProductWrongEndpoint(t *testing.T) {
	//given

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"There is no resource for path /rest/v2/electronics/productz/123","type":"UnknownResourceError"}]}`))
	}))
	defer ts.Close()

	p := prepareProductService(ts.URL)

	//when

	got, err := p.GetProductDetailsByID("123")

	//then

	assert.NotNil(t, err)

	wanted := ProductDTO{}

	assert.Equal(t, wanted, got)
}

func TestGetProductEmptyID(t *testing.T) {
	//given

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"There is no resource for path /rest/v2/electronics/products/","type":"UnknownResourceError"}]}`))
	}))
	defer ts.Close()

	p := prepareProductService(ts.URL)

	//when

	got, err := p.GetProductDetailsByID("0")

	//then

	assert.NotNil(t, err)

	wanted := ProductDTO{}

	assert.Equal(t, wanted, got)
}
