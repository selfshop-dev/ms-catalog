package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/testhelpers/e2etest"
)

func TestProduct_CreateAndGet(t *testing.T) {
	e := e2etest.NewExpect(t)

	id := e.POST("/api/v1/products").
		WithJSON(map[string]any{
			"name":        "Widget",
			"slug":        "widget",
			"price_cents": 1000,
			"currency":    "USD",
			"status":      "draft",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		Value("data").Object().
		Value("id").String().Raw()

	e.GET("/api/v1/products/"+id).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		HasValue("id", id).
		HasValue("name", "Widget").
		HasValue("slug", "widget")
}

func TestProduct_DuplicateSlug_Conflict(t *testing.T) {
	e := e2etest.NewExpect(t)

	body := map[string]any{
		"name":        "Widget",
		"slug":        fmt.Sprintf("widget-conflict-%s", uuid.New().String()[:8]),
		"price_cents": 0,
		"currency":    "USD",
		"status":      "draft",
	}

	e.POST("/api/v1/products").
		WithJSON(body).
		Expect().
		Status(http.StatusCreated)

	e.POST("/api/v1/products").
		WithJSON(body).
		Expect().
		Status(http.StatusConflict).
		HasContentType("application/problem+json")
}

func TestProduct_FullLifecycle(t *testing.T) {
	e := e2etest.NewExpect(t)

	slug := fmt.Sprintf("lifecycle-%s", uuid.New().String()[:8])

	id := e.POST("/api/v1/products").
		WithJSON(map[string]any{
			"name":        "Lifecycle Product",
			"slug":        slug,
			"price_cents": 500,
			"currency":    "USD",
			"status":      "draft",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		Value("data").Object().
		Value("id").String().Raw()

	// activate
	e.PUT("/api/v1/products/" + id + "/status").
		WithJSON(map[string]any{"status": "active"}).
		Expect().
		Status(http.StatusOK)

	// appears in active list
	e.GET("/api/v1/products/active").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		Value("items").Array().
		NotEmpty()

	// update
	updatedSlug := fmt.Sprintf("lifecycle-updated-%s", uuid.New().String()[:8])
	e.PUT("/api/v1/products/"+id).
		WithJSON(map[string]any{
			"name":        "Lifecycle Product Updated",
			"slug":        updatedSlug,
			"price_cents": 999,
			"currency":    "USD",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		HasValue("id", id)

	// get by updated slug
	e.GET("/api/v1/products/"+updatedSlug+"/by-slug").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		HasValue("slug", updatedSlug)

	// delete
	e.DELETE("/api/v1/products/" + id).
		Expect().
		Status(http.StatusNoContent)

	// gone
	e.GET("/api/v1/products/" + id).
		Expect().
		Status(http.StatusNotFound).
		HasContentType("application/problem+json")
}

func TestProduct_NotFound(t *testing.T) {
	e := e2etest.NewExpect(t)

	e.GET("/api/v1/products/" + uuid.New().String()).
		Expect().
		Status(http.StatusNotFound).
		HasContentType("application/problem+json")
}

func TestProduct_InvalidUUID(t *testing.T) {
	e := e2etest.NewExpect(t)

	e.GET("/api/v1/products/not-a-uuid").
		Expect().
		Status(http.StatusBadRequest).
		HasContentType("application/problem+json")
}

func TestProduct_ValidationError(t *testing.T) {
	e := e2etest.NewExpect(t)

	e.POST("/api/v1/products").
		WithJSON(map[string]any{
			"price_cents": 1000,
		}).
		Expect().
		Status(http.StatusUnprocessableEntity).
		HasContentType("application/problem+json").
		Body().Contains(`"fields"`)
}

func TestProduct_GetActiveProducts_Pagination(t *testing.T) {
	e := e2etest.NewExpect(t)

	// create 5 active products with unique slugs
	for i := range 5 {
		id := e.POST("/api/v1/products").
			WithJSON(map[string]any{
				"name":        fmt.Sprintf("Pagination Product %d", i),
				"slug":        fmt.Sprintf("pagination-%s-%d", uuid.New().String()[:8], i),
				"price_cents": 0,
				"currency":    "USD",
				"status":      "draft",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().
			Value("data").Object().
			Value("id").String().Raw()

		e.PUT("/api/v1/products/" + id + "/status").
			WithJSON(map[string]any{"status": "active"}).
			Expect().
			Status(http.StatusOK)
	}

	e.GET("/api/v1/products/active").
		WithQuery("limit", 3).
		WithQuery("offset", 0).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		Value("items").Array().
		Length().IsEqual(3)
}
