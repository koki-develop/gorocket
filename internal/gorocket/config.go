package gorocket

import (
	_ "embed"
)

//go:embed config_default.yaml
var DefaultConfig []byte
