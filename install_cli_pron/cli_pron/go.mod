module github.com/colinarticulate/cli_pron

go 1.18

replace github.com/colinarticulate/pron => ../pron

replace github.com/colinarticulate/scanScheduler => ../scanScheduler

require (
	github.com/colinarticulate/pron v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
)

require (
	github.com/colinarticulate/dictionary v0.0.0-20210623083452-03c3a6c454d4 // indirect
	github.com/colinarticulate/scanScheduler v0.0.0-00010101000000-000000000000 // indirect
	github.com/cryptix/wav v0.0.0-20180415113528-8bdace674401 // indirect
	github.com/davidbarbera/articulate-pocketsphinx-go/xyz_plus v0.0.0-20220426123918-5dbb324fc755 // indirect
	github.com/maxhawkins/go-webrtcvad v0.0.0-20210121163624-be60036f3083 // indirect
)
