package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := NewConfigWithFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
}
