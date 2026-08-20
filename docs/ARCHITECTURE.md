# Application Architecture

## Status and source

This project adopts [HUMQ](https://github.com/kodaimura/humq) v1.1.0 and maps
its responsibilities to Go, Gin, and GORM. It was reviewed against upstream
commit [`d3c9150`](https://github.com/kodaimura/humq/commit/d3c9150a2b824e6197fbc87230a1dc6940631313).
This document is the local contract; a later HUMQ version does not change the
scaffold automatically.

## Responsibility mapping

| HUMQ responsibility | This repository |
| --- | --- |
| Handler | `internal/handler/<domain>/` |
| Usecase | Public actions in `internal/usecase/<domain>/` |
| Module | Table-oriented interfaces in `internal/module/` |
| Query | Read-specific implementations in `internal/query/` |
| Policy / Operation | Narrow unexported functions in the owning usecase package |
| External client | Interfaces and adapters with protocol-specific implementations |

A Go `Usecase` interface may group related public actions for dependency
injection. Each public method still represents one independently explainable
Usecase action and must keep its flow visible.

## Dependency direction

```text
Gin route
  -> Handler
       -> Usecase action
            -> Module -> GORM model / PostgreSQL
            -> Query  -> read model / PostgreSQL
            -> Policy
            -> Operation -> Module / Query / Policy
            -> External client
```

- Handler calls Usecase, not Module, Query, GORM, or external clients.
- Usecase may call Module, Query, Policy, Operation, and external clients.
- Module owns persistence for one table by default and does not call Usecase.
- Query is read-only and named after the view or observation it returns.
- `internal/core` contains infrastructure primitives, not business rules.
- Dependencies are constructed in `internal/app`; packages must not use hidden
  mutable globals for business behavior.

## Handler

Handlers own HTTP concerns: Gin parameters and context, request binding, DTO
conversion, cookies, status codes, and response serialization. Middleware or a
guard may authenticate the request and attach the caller identity.

Handlers must not query GORM, mutate models, decide business authorization,
open transactions, or coordinate multiple business steps. Pass the caller
identity and target identity to the Usecase when authorization is required.

## Usecase

Each exported Usecase method owns one business action. It controls validation
that depends on business state, authorization, ordering, external-I/O policy,
and transaction success or rollback.

Usecase may hold interfaces for its dependencies. Keep DTOs and result types
close to the action or domain package. Avoid generic service, manager, or helper
packages that hide the sequence. An exported Usecase action must not call
another exported action as a reuse mechanism. Reusable pure decisions belong in
an unexported policy function.

Keep database-dependent behavior in each action by default. Use a narrow
unexported operation only when multiple actions must preserve the same
invariant and divergent validation, errors, locks, or update order would cause
a concrete inconsistency. Similar code or a long action is not sufficient.

Usecase must not issue GORM queries directly. The exception is opening
`db.Transaction` to own a transaction boundary; all reads and writes inside it
must use Modules or Queries bound to the transaction.

## Module and Query

A Module owns basic reads and every write for one table. It may expose
`WithTx(*gorm.DB)` so the Usecase can bind all participants to the same
transaction. Module owns row locks, conditional updates, and GORM details, but
never commits or starts a transaction.

Use Query for joins, aggregation, reports, complex searches, or specialized
read models. Query never changes state. Basic lookup and list operations remain
in the table's Module.

GORM models under `internal/model` describe persistence. Do not expose secret
fields through HTTP DTOs.

## Transactions and concurrency

The Usecase starts `db.Transaction` when multiple changes must succeed or fail
together. Construct transaction-bound Modules with `WithTx` inside the
callback. Let returned errors roll back the callback.

HUMQ does not structurally guarantee multi-table consistency. A Usecase can
omit a required Module call while following every dependency rule. Protect
representable invariants with database constraints and test business branches
and rollback behavior. Use a structurally protective design for a domain that
cannot accept implementation-level enforcement.

Row locking and conditional writes belong to the owning Module. Read a value
with a lock inside the same transaction that validates and updates it. Database
constraints remain the final guard for representable invariants.

External I/O is not made atomic by a database transaction. Send best-effort
notifications after commit, or use an outbox when delivery must not be lost.

## Authentication and authorization

Authentication verifies credentials and identifies the caller. Authorization
is a separate Usecase decision. Define ownership, role, and administrative
rules for the product and deny access when no rule allows an operation.
Collection routes and routes accepting a target account ID require explicit
authorization before production use.

Middleware may verify token format and signature, but account-state validation
such as existence, disabled state, and token version goes through an
authentication Usecase rather than calling a Module directly.

## Testing and evolution

- Test pure policies and branch-heavy Usecase actions with Go tests.
- Test Module SQL, transactions, locks, and constraints against PostgreSQL when
  those details matter.
- Test request validation, authentication, response contracts, and complete
  flows through API E2E tests.
- Run `make check`; run `make test_e2e` for API, database, transaction, or
  migration changes.

A deliberate exception must be narrow, documented here with its reason and
covered at the level where its risk appears. Treat increasing Operations or
shared invariants that dominate a domain as a signal to consider an
aggregate-centered or other domain-specific architecture for that domain.
