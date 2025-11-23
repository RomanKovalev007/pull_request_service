package integration

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardown()
    os.Exit(code)
}


