module pvault

go 1.24

require github.com/binary-soup/go-commando v1.2.2

replace github.com/binary-soup/go-commando => ../go-command

require (
	github.com/atotto/clipboard v0.1.4
	golang.org/x/crypto v0.40.0
)

require (
	github.com/google/uuid v1.6.0
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/term v0.33.0
)
