package dbtest

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	ctxval "github.com/selfshop-dev/lib-ctxval"
)

// Querier is a constructor function that creates a query executor from a transaction.
// Pass gen.New from your sqlc-generated package.
type Querier[Q any] func(pgx.Tx) Q

// Suite is a generic testify-suite for adapter integration tests.
//
// Each test case gets its own transaction in SetupTest.
// TearDownTest rolls it back — no data leaks between tests, no TRUNCATE needed.
//
// The adapter receives Q through ctxval.Or(ctx, a.g), so the transaction
// is picked up transparently without any changes to production code.
type Suite[Q any] struct {
	suite.Suite

	tx   pgx.Tx
	ctx  context.Context
	Pool *pgxpool.Pool
	q    Querier[Q]

	schema []byte
}

// NewSuite creates a Suite with the given querier constructor.
//
//	type productAdapterSuite struct {
//	    dbtest.Suite[*gen.Queries]
//	}
//
//	func (s *productAdapterSuite) SetupSuite() {
//	    s.Suite = dbtest.NewSuite(migrations.CurrentSchemaSQL, gen.New)
//	    s.Suite.SetupSuite()
//	}
func NewSuite[Q any](schema []byte, q Querier[Q]) Suite[Q] {
	return Suite[Q]{schema: schema, q: q}
}

func (s *Suite[Q]) Init(schema []byte, q Querier[Q]) {
	s.schema = schema
	s.q = q
}

func (s *Suite[Q]) SetupSuite() { s.Pool = MustGetPool(s.T(), s.schema) }

func (s *Suite[Q]) SetupTest() {
	tx, err := s.Pool.Begin(context.Background())
	s.Require().NoError(err)
	s.tx = tx
	s.ctx = ctxval.With(context.Background(), s.q(tx))
}

// Ctx returns the context with the current test transaction.
func (s *Suite[Q]) Ctx() context.Context { return s.ctx }

func (s *Suite[Q]) TearDownTest() { _ = s.tx.Rollback(context.Background()) } //nolint:errcheck // rollback in teardown: error is irrelevant, test already finished
