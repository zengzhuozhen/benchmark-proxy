package resources

import (
	_ "embed"
)

//go:embed ca.crt
var CaCrt []byte

//go:embed ca.key
var CaKey []byte

var AuthKey string
