# pg-go

**pg-go** compares three Go libraries for Postgres access: [goqu](https://github.com/doug-martin/goqu), [squirrel](https://github.com/Masterminds/squirrel) and [gorm](https://github.com/go-gorm/gorm).

## How the comparison works

- [`repository`](repository) defines one contract: domain types, a `Store` interface with user, post and comment repositories, pagination parameters, and sentinel errors.
- [`goqurepo`](goqurepo), [`squirrelrepo`](squirrelrepo) and [`gormrepo`](gormrepo) implement that contract once each. Every implementation owns its row types and struct tags and maps them to the shared domain types.
- [`integration`](integration) runs the same test suite against all three:
  - Postgres runs in a container started with testcontainers. It is migrated with goose and seeded once, then restored from a snapshot before every test.
  - The expected results of each scenario are stored once in [`integration/testdata/results`](integration/testdata/results), and all three implementations must return exactly that data.
  - A pgx query tracer records the SQL that each library sends to Postgres. It is stored per scenario and implementation in [`integration/testdata/sql`](integration/testdata/sql).
- Every library runs with its default options. Only the wiring needed to talk to Postgres is configured: the goqu Postgres dialect and squirrel's `$n` placeholders.
- All three libraries use the same driver: pgx through `database/sql`.

The seed in [`seed/seed.sql`](seed/seed.sql) is plain SQL. It uses `generate_series` with modular arithmetic, so the data is identical on every run, and the IDs are readable fixed UUIDs. It generates about 4,100 rows without difficulty.

## What is taken into consideration?

- **CRUD operations:** create, read, update and delete, with `RETURNING`.
- **Bulk insert:** 100 rows in one statement.
- **Runtime query building:** a REST-style list with optional filters, sorting by two fields, limit, offset and a total count.
- **Joins and nested mapping:** one-to-one into a nested struct, and one-to-many two levels deep.
- **Aggregates:** `COUNT(DISTINCT)`, `MAX`, `GROUP BY`, `HAVING`, and window functions.
- **Postgres features:** `ON CONFLICT` upsert, `WITH`, `WITH RECURSIVE`, and `FOR UPDATE SKIP LOCKED`.
- **Escape hatches:** how much raw SQL each scenario needs.
- **Errors:** not found and unique violation.
- **Transactions:** commit, rollback on error, savepoints, and one transaction shared by several repositories.

## Results

### How much raw SQL each scenario needs

- `native`: the library's API expresses the construct.
- `fragment`: the library's API plus raw SQL pieces for what it can't express, such as `goqu.L`, `squirrel.Expr`/`Prefix`/`Suffix`, or a gorm string condition with a Postgres function.
- `raw`: the construct under test is handwritten SQL run through the library.

Plain column conditions such as gorm's `Where("id = ?", id)` or squirrel's `OrderBy("id ASC")` count as `native`, because strings are the normal API of those libraries.

Each scenario name links to the SQL that the three libraries generate.

| Scenario | Criterion | goqu | squirrel | gorm |
| --- | --- | --- | --- | --- |
| [`users_create`](integration/testdata/sql/users_create) | CRUD, `RETURNING` | [native](goqurepo/users.go) | [fragment¹](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_create_many`](integration/testdata/sql/users_create_many) | Bulk insert | [native](goqurepo/users.go) | [fragment¹](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_get_by_id`](integration/testdata/sql/users_get_by_id) | CRUD | [native](goqurepo/users.go) | [native](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_update`](integration/testdata/sql/users_update) | CRUD, partial update | [native](goqurepo/users.go) | [fragment¹](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_delete`](integration/testdata/sql/users_delete) | CRUD, `RETURNING` | [native](goqurepo/users.go) | [fragment¹](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_upsert_update`](integration/testdata/sql/users_upsert_update) | `ON CONFLICT` | [native](goqurepo/users.go) | [fragment](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_list_filtered`](integration/testdata/sql/users_list_filtered) | Runtime query, pagination | [native](goqurepo/users.go) | [native](squirrelrepo/users.go) | [native](gormrepo/users.go) |
| [`users_list_eligible_adults`](integration/testdata/sql/users_list_eligible_adults) | Nested `AND`/`OR` conditions | [fragment²](goqurepo/users.go) | [fragment²](squirrelrepo/users.go) | [fragment²](gormrepo/users.go) |
| [`users_unverify_minors`](integration/testdata/sql/users_unverify_minors) | Update many, `RETURNING` | [fragment²](goqurepo/users.go) | [fragment¹²](squirrelrepo/users.go) | [fragment²](gormrepo/users.go) |
| [`users_top_age_groups`](integration/testdata/sql/users_top_age_groups) | `WITH`, window function | [fragment²](goqurepo/users.go) | [fragment](squirrelrepo/users.go) | [raw³](gormrepo/users.go) |
| [`users_post_stats`](integration/testdata/sql/users_post_stats) | Join, `GROUP BY`, `HAVING` | [native](goqurepo/users.go) | [fragment](squirrelrepo/users.go) | [fragment](gormrepo/users.go) |
| [`posts_get_with_author`](integration/testdata/sql/posts_get_with_author) | Join, one-to-one mapping | [native](goqurepo/posts.go) | [fragment⁴](squirrelrepo/posts.go) | [native](gormrepo/posts.go) |
| [`posts_list_recent_with_comments`](integration/testdata/sql/posts_list_recent_with_comments) | One-to-many mapping, two levels | [native⁵](goqurepo/posts.go) | [fragment⁴⁵](squirrelrepo/posts.go) | [native](gormrepo/posts.go) |
| [`comments_thread`](integration/testdata/sql/comments_thread) | `WITH RECURSIVE` | [fragment](goqurepo/comments.go) | [raw](squirrelrepo/comments.go) | [raw](gormrepo/comments.go) |
| [`tx_commit`](integration/testdata/sql/tx_commit) | Transaction across repositories | [native](goqurepo/store.go) | [native](squirrelrepo/store.go) | [native](gormrepo/store.go) |
| [`tx_savepoint`](integration/testdata/sql/tx_savepoint) | Savepoints | [raw⁶](goqurepo/store.go) | [raw⁶](squirrelrepo/store.go) | [native](gormrepo/store.go) |
| [`tx_skip_locked`](integration/testdata/sql/tx_skip_locked) | `FOR UPDATE SKIP LOCKED` | [native](goqurepo/users.go) | [fragment](squirrelrepo/users.go) | [native](gormrepo/users.go) |

1. squirrel has no `RETURNING` builder, so it is added with `Suffix("RETURNING *")`.
2. `age()` and interval literals have no builder API in any of the three libraries.
3. The `WITH` statement is written with `Raw`. The inner query is built with gorm and passed in as a subquery argument.
4. sqlx scans nested structs, but the aliased column list (`post.id AS "post.id"`) is built by hand.
5. The joined rows are grouped into nested structs in Go.
6. The library has no savepoint API, so nested `InTx` runs `SAVEPOINT`, `ROLLBACK TO SAVEPOINT` and `RELEASE SAVEPOINT` by hand.

All three libraries return identical results for every scenario. The full list of scenarios, including the not-found and invalid-input cases, is in [`integration`](integration).

### Round trips

| Scenario | goqu | squirrel | gorm |
| --- | --- | --- | --- |
| `users_list_filtered` (count + page) | 2 | 2 | 2 |
| `posts_get_with_author` | 1 | 1 | 1 |
| `posts_list_recent_with_comments` | 1 | 1 | 3 |
| `users_create` | 1 | 1 | 3 (`BEGIN`, `INSERT`, `COMMIT`) |

### Errors

| Need | goqu | squirrel | gorm |
| --- | --- | --- | --- |
| Single row not found | `ScanStruct` returns `found == false` | `sqlx.GetContext` returns `sql.ErrNoRows` | `First` returns `gorm.ErrRecordNotFound`; `RowsAffected == 0` for updates and deletes |
| Unique violation | `errors.As` into `*pgconn.PgError` with code `23505` | Same | Same, because `TranslateError` is off by default |

### Size

| Package | Non-blank lines |
| --- | --- |
| `goqurepo` | 627 |
| `squirrelrepo` | 615 |
| `gormrepo` | 521 |

## Notes per library

### goqu

Version v9.19.0, the latest release, published 2023-11-21.

- **Easy:**
  - Column lists come from struct tags, including nested structs scanned through dotted aliases (`"author"."id" AS "author.id"`).
  - `WITH`, `WITH RECURSIVE`, `UNION ALL`, `ON CONFLICT`, `RETURNING` and `FOR UPDATE SKIP LOCKED` all have builder methods.
- **Workarounds:**
  - There is no savepoint API.
  - `WithTx` has no context parameter, so transactions use `BeginTx`.
  - Postgres functions, intervals and arithmetic such as `"t"."depth" + 1` need `goqu.L`.
- **Defaults that shape the SQL:**
  - Values are written directly into the SQL text, because prepared mode is off.
  - `Eq(true)` renders as `IS TRUE`.
  - `ScanStruct` and `Count` add `LIMIT 1`.

### squirrel

Version v1.5.4, the latest release, published 2023-03-17. Its README declares the project "complete": bug fixes are still merged, slowly.

- **Easy:**
  - It has the smallest API, with conditions as maps (`Eq`, `Like`, `Gt`).
  - Builders are immutable, so the list filters are built once and reused for both the count and the page.
- **Workarounds:**
  - squirrel doesn't scan rows; sqlx does.
  - There are no builders for `RETURNING`, `ON CONFLICT`, `FOR UPDATE`, `WITH` or `UNION`; they all go through `Prefix`/`Suffix` strings.
  - A builder nested inside `Expr` is rendered with its own placeholder format. CTE subqueries are therefore built with `squirrel.Select` (`?` placeholders), and only the outer statement converts them to `$n`.
- **Defaults that shape the SQL:**
  - Placeholders default to `?`, which Postgres rejects, so `PlaceholderFormat(squirrel.Dollar)` is required.

### gorm

Version v1.31.2, the latest release, published 2026-06-22.

- **Easy:**
  - `Create`, `Updates` and `Delete` combine with `clause.Returning`, `clause.OnConflict` and `clause.Locking`.
  - `Joins("Author")` loads a one-to-one relation in the same query.
  - `Preload("Comments.Reactions")` loads nested one-to-many relations without manual grouping.
  - A nested `Transaction` uses a savepoint automatically.
- **Workarounds:**
  - Row types need `TableName`, `default:uuid_generate_v4()` on IDs, and `autoUpdateTime:false` on `UpdatedAt`. Without that tag, gorm sets `updated_at` itself on insert.
  - A chained query must become a `Session` before it can be reused for both the count and the page.
  - There is no CTE support, so `WITH` statements use `Raw`.
- **Defaults that shape the SQL:**
  - Every single write is wrapped in `BEGIN`/`COMMIT`.
  - `Create` returns only `"id"`.
  - `First` adds `ORDER BY` primary key and `LIMIT 1`.
  - `Preload` runs one query per level.
  - Nested transactions use random savepoint names. They are stored as `sp<random>` in the SQL files.
  - The default logger prints errors to stdout, including "record not found".

## Running it

Prerequisites: a Docker-compatible runtime and [mise](https://mise.jdx.dev), which installs the Go toolchain and the goose CLI.

```sh
mise install
go test ./...
```

Rewrite the stored SQL for every implementation, and the expected results from goqu:

```sh
go test ./integration -update
```

Run one scenario for one library:

```sh
go test ./integration -run 'TestUsers/users_list_filtered/gorm'
```

Create a new migration:

```sh
goose --dir migrations create <name> sql
```

## Project layout

```
migrations/      goose migrations, embedded for the tests
seed/            plain SQL seed and fixture IDs
repository/      shared contract
goqurepo/        goqu implementation
squirrelrepo/    squirrel + sqlx implementation
gormrepo/        gorm implementation
integration/     shared test suite, expected results and recorded SQL
```
