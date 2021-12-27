package ini

import (
	"testing"
)

// //file:/C:/CIGAM/CIGAM-11/Magic.ini

// Test LoadIni function
func TestLoadIni(t *testing.T) {
	var config = LoadIni("MagicTest.ini")
	if config.Ini == nil {
		t.Error("LoadMagicIni() failed")
	}
}

func TestTranslate(t *testing.T) {
	var config = LoadIni("MagicTest.ini")
	var result = config.Translate("%CIGAM_INSTAL%")
	if result == "" {
		t.Error("Translate() failed")
	}
}

func TestCigamSql(t *testing.T) {
	var config = LoadIni("MagicTest.ini")
	var CIGAM_SQL = config.Get("MAGIC_DATABASES", "CIGAM_SQL")
	if CIGAM_SQL == "" {
		t.Error("CIGAM_SQL not found")
	}
}
