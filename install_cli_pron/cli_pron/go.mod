module github.com/colinarticulate/cli_pron

go 1.18

replace github.com/colinarticulate/pron => ../pron

replace github.com/colinarticulate/scanScheduler => ../scanScheduler

require (
	github.com/colinarticulate/pron v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
)
