package pkg

import (
	"fmt"
	"testing"
)

func TestParseDependencies(t *testing.T) {
	deps, err := parseDependencies(map[string]string{
		"ChrisMcKenzie/test":     "0.0.1",
		"ChrisMcKenzie/testa":    "1.x",
		"ChrisMcKenzie/testb":    "~1.2",
		"github.com/spf13/cobra": "git+https://github.com/spf13/cobra#master",
	})

	if err != nil {
		t.Error(err)
	}

	if len(deps) != 4 {
		t.Fail()
	}

	fmt.Println(deps)
}
