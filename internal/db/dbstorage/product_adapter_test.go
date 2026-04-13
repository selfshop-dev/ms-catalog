package dbstorage_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"

	apperr "github.com/selfshop-dev/lib-apperr"

	"github.com/selfshop-dev/ms-catalog/internal/db/dbstorage"
	"github.com/selfshop-dev/ms-catalog/internal/db/gen"
	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/migrations"
	"github.com/selfshop-dev/ms-catalog/pkg/db/dbtest"
)

type productAdapterSuite struct {
	dbtest.Suite[*gen.Queries]
	r domain.ProductRepository
}

func (s *productAdapterSuite) SetupSuite() {
	s.Init(
		migrations.CurrentSchemaSQL,
		func(tx pgx.Tx) *gen.Queries { return gen.New(tx) },
	)
	s.Suite.SetupSuite()
	s.r = dbstorage.NewProductAdapter(gen.New(s.Pool))
}

func (s *productAdapterSuite) newProduct(name, slug string, status domain.ProductStatus) *domain.Product {
	p, err := domain.NewProduct(
		name, slug,
		new("Description for "+name), nil, nil,
		0, "USD",
		new(status),
	)
	s.Require().NoError(err)
	return p
}

func (s *productAdapterSuite) createProduct(name, slug string, status domain.ProductStatus) *domain.Product {
	created, err := s.r.Create(s.Ctx(), s.newProduct(name, slug, status))
	s.Require().NoError(err)
	return created
}

func TestProductAdapter(t *testing.T) {
	suite.Run(t, new(productAdapterSuite))
}

func (s *productAdapterSuite) TestCreate_ok() {
	p := s.newProduct("Widget", "widget", domain.ProductStatusDraft)

	got, err := s.r.Create(s.Ctx(), p)

	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, got.ID())
	s.Require().Equal("Widget", got.Name())
	s.Require().Equal("widget", got.Slug())
	s.Require().Equal(domain.ProductStatusDraft, got.Status())
	s.Require().False(got.CreatedAt().IsZero())
	s.Require().False(got.UpdatedAt().IsZero())
	s.Require().Nil(got.DeletedAt())
}

func (s *productAdapterSuite) TestCreate_withOptionalFields() {
	p, err := domain.NewProduct(
		"Full Product",
		"full-product",
		new("Full description"),
		new("Short desc"),
		new("https://cdn.example.com/image.jpg"),
		1999,
		"EUR",
		new(domain.ProductStatusDraft),
	)
	s.Require().NoError(err)

	got, err := s.r.Create(s.Ctx(), p)

	s.Require().NoError(err)
	s.Require().Equal("Short desc", *got.ShortDescription())
	s.Require().NotNil(got.DisplayImageURL())
	s.Require().Equal("https://cdn.example.com/image.jpg", *got.DisplayImageURL())
	s.Require().Equal(int64(1999), got.PriceCents())
	s.Require().Equal("EUR", got.Currency())
}

func (s *productAdapterSuite) TestCreate_duplicateSlug_conflict() {
	s.createProduct("First", "dup-slug", domain.ProductStatusDraft)

	_, err := s.r.Create(s.Ctx(), s.newProduct("Second", "dup-slug", domain.ProductStatusDraft))

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindConflict))
}

func (s *productAdapterSuite) TestGetByID_ok() {
	created := s.createProduct("Gadget", "gadget", domain.ProductStatusDraft)

	got, err := s.r.GetByID(s.Ctx(), created.ID())

	s.Require().NoError(err)
	s.Require().Equal(created.ID(), got.ID())
	s.Require().Equal("Gadget", got.Name())
	s.Require().Equal("gadget", got.Slug())
}

func (s *productAdapterSuite) TestGetByID_notFound() {
	_, err := s.r.GetByID(s.Ctx(), uuid.New())

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindNotFound))
}

func (s *productAdapterSuite) TestGetBySlug_ok() {
	s.createProduct("Thingamajig", "thingamajig", domain.ProductStatusDraft)

	got, err := s.r.GetBySlug(s.Ctx(), "thingamajig")

	s.Require().NoError(err)
	s.Require().Equal("thingamajig", got.Slug())
}

func (s *productAdapterSuite) TestGetBySlug_notFound() {
	_, err := s.r.GetBySlug(s.Ctx(), "no-such-slug")

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindNotFound))
}

func (s *productAdapterSuite) TestUpdate_ok() {
	created := s.createProduct("Old Name", "old-name", domain.ProductStatusDraft)

	updated := s.newProduct("New Name", "new-name", domain.ProductStatusDraft)
	got, err := s.r.Update(s.Ctx(), created.ID(), updated)

	s.Require().NoError(err)
	s.Require().Equal(created.ID(), got.ID())
	s.Require().Equal("New Name", got.Name())
	s.Require().Equal("new-name", got.Slug())
}

func (s *productAdapterSuite) TestUpdate_priceCents() {
	created := s.createProduct("Priced", "priced-product", domain.ProductStatusDraft)

	p, err := domain.NewProduct(
		"Priced", "priced-product",
		new("desc"), nil, nil,
		750, "RUB",
		new(domain.ProductStatusDraft),
	)
	s.Require().NoError(err)

	got, err := s.r.Update(s.Ctx(), created.ID(), p)

	s.Require().NoError(err)
	s.Require().Equal(int64(750), got.PriceCents())
	s.Require().Equal("RUB", got.Currency())
}

func (s *productAdapterSuite) TestUpdate_notFound() {
	_, err := s.r.Update(s.Ctx(), uuid.New(), s.newProduct("X", "x-slug", domain.ProductStatusDraft))

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindNotFound))
}

func (s *productAdapterSuite) TestUpdate_duplicateSlug_conflict() {
	s.createProduct("Alpha", "alpha", domain.ProductStatusDraft)
	beta := s.createProduct("Beta", "beta", domain.ProductStatusDraft)

	_, err := s.r.Update(s.Ctx(), beta.ID(), s.newProduct("Beta Renamed", "alpha", domain.ProductStatusDraft))

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindConflict))
}

func (s *productAdapterSuite) TestUpdateStatus_allTransitions() {
	transitions := []struct {
		name   string
		status domain.ProductStatus
	}{
		{"to active", domain.ProductStatusActive},
		{"to inactive", domain.ProductStatusInactive},
		{"to archived", domain.ProductStatusArchived},
		{"back to draft", domain.ProductStatusDraft},
	}

	created := s.createProduct("Status Product", "status-product", domain.ProductStatusDraft)

	for _, tr := range transitions {
		got, err := s.r.UpdateStatus(s.Ctx(), created.ID(), tr.status)
		s.Require().NoError(err)
		s.Require().Equal(tr.status, got.Status(), tr.name)
	}
}

func (s *productAdapterSuite) TestUpdateStatus_notFound() {
	_, err := s.r.UpdateStatus(s.Ctx(), uuid.New(), domain.ProductStatusActive)

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindNotFound))
}

func (s *productAdapterSuite) TestDelete_ok() {
	created := s.createProduct("To Delete", "to-delete", domain.ProductStatusDraft)

	err := s.r.Delete(s.Ctx(), created.ID())
	s.Require().NoError(err)

	_, err = s.r.GetByID(s.Ctx(), created.ID())
	s.Require().True(apperr.IsKind(err, apperr.KindNotFound))
}

func (s *productAdapterSuite) TestDelete_notFound() {
	err := s.r.Delete(s.Ctx(), uuid.New())

	s.Require().Error(err)
	s.Require().True(apperr.IsKind(err, apperr.KindNotFound))
}

func (s *productAdapterSuite) TestGetListActiveProducts_emptyResult() {
	products, err := s.r.GetListActiveProducts(s.Ctx(), 10, 0)

	s.Require().NoError(err)
	s.Require().Empty(products)
}

func (s *productAdapterSuite) TestGetListActiveProducts_returnsOnlyActive() {
	active := s.createProduct("Active One", "active-one", domain.ProductStatusActive)
	s.createProduct("Draft", "draft-only", domain.ProductStatusDraft)
	s.createProduct("Inactive", "inactive-one", domain.ProductStatusInactive)

	products, err := s.r.GetListActiveProducts(s.Ctx(), 10, 0)

	s.Require().NoError(err)
	s.Require().Len(products, 1)
	s.Require().Equal(active.ID(), products[0].ID())
}

func (s *productAdapterSuite) TestGetListActiveProducts_pagination() {
	for i := range 5 {
		s.createProduct(
			fmt.Sprintf("Product %d", i),
			fmt.Sprintf("product-%d", i), domain.ProductStatusActive,
		)
	}

	page1, err := s.r.GetListActiveProducts(s.Ctx(), 3, 0)
	s.Require().NoError(err)
	s.Require().Len(page1, 3)

	page2, err := s.r.GetListActiveProducts(s.Ctx(), 3, 3)
	s.Require().NoError(err)
	s.Require().Len(page2, 2)
}
