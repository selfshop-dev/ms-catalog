package migrations

import _ "embed"

// CurrentSchemaSQL holds the complete DDL schema for the database.
// Used in integration tests to initialize the test container.
//
//go:embed current_schema.sql
var CurrentSchemaSQL []byte
