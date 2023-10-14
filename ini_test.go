package ini

import (
	"os"
	"testing"
)

func TestNewMagicIni(t *testing.T) {
	ini := NewMagicIni()
	if ini.TranslatableSection != DEFAULT_TRANSLATABLE_SECTION {
		t.Fatalf("Expected default translatable section to be %s, got %s", DEFAULT_TRANSLATABLE_SECTION, ini.TranslatableSection)
	}
}

func TestLoadIni(t *testing.T) {
	content := "[section1]\nkey1=value1\n"
	tmpfile, err := os.CreateTemp("", "example.*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString(content)

	ini := NewMagicIni()
	if err := ini.LoadIni(tmpfile.Name()); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val := ini.Get("section1", "key1"); val != "value1" {
		t.Fatalf("Expected value1, got %s", val)
	}
}

func TestParseSection(t *testing.T) {
	ini := NewMagicIni()
	ini.Parse("[section2]")
	if _, exists := ini.Ini["section2"]; !exists {
		t.Fatal("section2 should exist")
	}
}

func TestParseSectionInlineKeyValue(t *testing.T) {
	ini := NewMagicIni()
	ini.Parse("/[section3] key=value")
	if val := ini.Get("section3", "key"); val != "value" {
		t.Fatalf("Expected value, got %s", val)
	}
}

func TestParseKeyValue(t *testing.T) {
	ini := NewMagicIni()
	ini.currentParsedSection = "section4"
	ini.Parse("key2=value2")
	if val := ini.Get("section4", "key2"); val != "value2" {
		t.Fatalf("Expected value2, got %s", val)
	}
}

func TestTranslate(t *testing.T) {
	ini := NewMagicIni()
	ini.Ini[DEFAULT_TRANSLATABLE_SECTION] = map[string]string{
		"test": "translated",
	}
	result := ini.Translate("%test%")
	if result != "translated" {
		t.Fatalf("Expected translated, got %s", result)
	}
}

func TestSetAndGet(t *testing.T) {
	ini := NewMagicIni()
	ini.Set("section5", "key3", "value3")
	if val := ini.Get("section5", "key3"); val != "value3" {
		t.Fatalf("Expected value3, got %s", val)
	}
}

func TestSave(t *testing.T) {
	ini := NewMagicIni()
	ini.Set("section6", "key4", "value4")
	tmpfile, err := os.CreateTemp("", "save_example.*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if err := ini.Save(tmpfile.Name()); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
