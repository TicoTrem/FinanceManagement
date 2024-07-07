module Service

go 1.22.3

require github.com/go-sql-driver/mysql v1.8.1
replace github.com/ticotrem/shared => ../shared

require (
    filippo.io/edwards25519 v1.1.0 // indirect
	github.com/ticotrem/shared v0.0.0-unpublished
)