package ini

import "testing"

// Test LoadIni function
func TestLoadIni(t *testing.T) {
	var config = LoadIni("Magic.ini")
	if config.Ini == nil {
		t.Error("LoadMagicIni() failed")
	}
}
