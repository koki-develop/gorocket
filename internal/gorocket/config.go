package gorocket

import (
	_ "embed"
)

//go:embed .gorocket.yml
var DefaultConfig []byte
