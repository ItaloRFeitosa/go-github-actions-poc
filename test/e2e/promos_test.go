package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-resty/resty/v2"
	"github.com/italorfeitosa/go-github-actions-poc/internal/app"
	"github.com/italorfeitosa/go-github-actions-poc/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromosEndpoints(t *testing.T) {
	url := os.Getenv("TEST_SERVER_BASE_URL")
	if url == "" {
		url = "http://localhost:8080"
	}
	r := resty.New().SetBaseURL(url)
	t.Run("POST /promos", func(t *testing.T) {
		t.Run("given valid promo params when persisted sucessfully should return promo data with status code 200", func(t *testing.T) {
			createPromoReq := app.PromoModel{
				ProductName: gofakeit.BookTitle(),
				Link:        gofakeit.URL(),
			}
			promoDataResponse, err := CreatePromo(r, createPromoReq)
			require.NoError(t, err)
			assert.Equal(t, createPromoReq.ProductName, promoDataResponse.Data.ProductName)
			assert.Equal(t, createPromoReq.Link, promoDataResponse.Data.Link)
			assert.False(t, promoDataResponse.Data.CreatedAt.IsZero())
			assert.Greater(t, promoDataResponse.Data.ID, uint(0))
		})
	})

	t.Run("GET /promos", func(t *testing.T) {
		t.Run("after create a promo should return new promo in the top of list", func(t *testing.T) {
			createPromoReq := app.PromoModel{
				ProductName: gofakeit.BookTitle(),
				Link:        gofakeit.URL(),
			}
			promoDataResponse, err := CreatePromo(r, createPromoReq)
			require.NoError(t, err)

			promosDataResponse, err := GetPromos(r)
			require.NoError(t, err)

			assert.Equal(t, promoDataResponse.Data.ID, promosDataResponse.Data[0].ID)
			assert.Equal(t, promoDataResponse.Data.Link, promosDataResponse.Data[0].Link)
			assert.Equal(t, promoDataResponse.Data.ProductName, promosDataResponse.Data[0].ProductName)
			assert.Equal(t, promoDataResponse.Data.CreatedAt.Unix(), promosDataResponse.Data[0].CreatedAt.Unix())
		})
	})

	t.Run("GET /promos/{id}", func(t *testing.T) {
		t.Run("after create a promo should return get promo with id", func(t *testing.T) {
			createPromoReq := app.PromoModel{
				ProductName: gofakeit.BookTitle(),
				Link:        gofakeit.URL(),
			}
			expectedPromoDataResponse, err := CreatePromo(r, createPromoReq)
			require.NoError(t, err)

			promoDataResponse, err := GetPromo(r, expectedPromoDataResponse.Data.ID)
			require.NoError(t, err)

			assert.Equal(t, expectedPromoDataResponse.Data.ID, promoDataResponse.Data.ID)
			assert.Equal(t, expectedPromoDataResponse.Data.Link, promoDataResponse.Data.Link)
			assert.Equal(t, expectedPromoDataResponse.Data.ProductName, promoDataResponse.Data.ProductName)
			assert.Equal(t, expectedPromoDataResponse.Data.CreatedAt.Unix(), promoDataResponse.Data.CreatedAt.Unix())
		})
	})
}

func CreatePromo(c *resty.Client, req app.PromoModel) (PromoDataResponse, error) {
	return httpclient.Resty[PromoDataResponse](c).
		Post("/promos").
		Body(req).
		StatusCode(http.StatusCreated).
		Exec()
}

func GetPromo(c *resty.Client, id uint) (PromoDataResponse, error) {
	return httpclient.Resty[PromoDataResponse](c).
		Get("/promos/{id}").
		Param("id", id).
		StatusCode(http.StatusOK).
		Exec()
}

type PromoDataResponse struct {
	Data app.PromoModel `json:"data"`
}

func GetPromos(c *resty.Client) (PromosDataResponse, error) {
	return httpclient.Resty[PromosDataResponse](c).
		Get("/promos").
		StatusCode(http.StatusOK).
		Exec()
}

type PromosDataResponse struct {
	Data []app.PromoModel `json:"data"`
}
