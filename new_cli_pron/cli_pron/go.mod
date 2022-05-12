module github.com/colinarticulate/cli_pron

go 1.13

replace github.com/colinarticulate/pron => ../pron

replace github.com/colinarticulate/scanScheduler => ../scanScheduler

require (
	github.com/colinarticulate/pron v0.0.0-00010101000000-000000000000
	github.com/davidbarbera/xyz_plus v1.1.2 // indirect
	github.com/google/uuid v1.3.0
)
