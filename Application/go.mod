module github.com/ticotrem/finance

go 1.23

replace github.com/ticotrem/finance/shared => ../Shared

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1
	github.com/ticotrem/finance/shared v0.0.0-unpublished
)
