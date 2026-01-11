package board

import (
	"chess/magic"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Initialize magic bitboards before running tests
	if err := magic.Prepare(); err != nil {
		panic("Failed to load magic data: " + err.Error())
	}
	os.Exit(m.Run())
}
