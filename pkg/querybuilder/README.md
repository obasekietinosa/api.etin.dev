# `pkg/querybuilder`

This package provides composable helpers for building SQL statements. The core types expose fluent builders for the standard CRUD operations:

- `select.go` constructs `SELECT` queries with filtering, ordering, and pagination helpers.
- `insert.go` builds `INSERT` statements with support for returning clauses.
- `update.go` assembles `UPDATE` queries with conditional sets and filters.
- `delete.go` creates `DELETE` statements.

See the accompanying tests for usage examples that cover the supported query patterns.
