module github.com/colinarticulate/cli_pron

go 1.13

replace github.com/colinarticulate/pron => ../pron

replace github.com/colinarticulate/scanScheduler => ../scanScheduler

require (
	github.com/colinarticulate/pron v1.1.0
	github.com/google/uuid v1.3.0
)
