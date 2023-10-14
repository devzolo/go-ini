// Package ini provides structures and methods for parsing and managing INI files with special enhancements.
package ini

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const DEFAULT_TRANSLATABLE_SECTION = "MAGIC_LOGICAL_NAMES"

// MagicIni represents the structured content of an INI file.
// It provides functionalities for parsing, modifying, and saving INI content.
type MagicIni struct {
	Ini                  map[string]map[string]string // Parsed INI content organized as Section -> Key -> Value.
	currentParsedSection string                       // Tracks the currently parsed section while reading an INI file.
	TranslatableSection  string                       // Specifies the section containing translatable strings.
	SectionsOrder        []string                     // Maintains the order of sections as they appear in the INI file.
	KeysOrder            map[string][]string          // Maintains the order of keys within sections.
}

// NewMagicIni initializes a new MagicIni instance with default values.
// By default, it sets the translatable section as "MAGIC_LOGICAL_NAMES".
func NewMagicIni() *MagicIni {
	return &MagicIni{
		Ini:                 make(map[string]map[string]string),
		TranslatableSection: DEFAULT_TRANSLATABLE_SECTION,
		KeysOrder:           make(map[string][]string),
	}
}

// LoadIni reads and parses the content of the INI file located at the specified path.
func (ini *MagicIni) LoadIni(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}
	defer f.Close()
	source, err := io.ReadAll(f)

	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	ini.Parse(string(source))
	return nil
}

// LoadAdditionalIni merges content from another INI file into the current instance.
// Useful for layering configurations or adding supplemental data.
func (ini *MagicIni) LoadAdditionalIni(path string) error {
	return ini.LoadIni(path)
}

// replaceLineScapes removes line continuation sequences from a given string.
func replaceLineScapes(s string) string {
	s = strings.Replace(s, "+\r\n", "", -1)
	s = strings.Replace(s, "+\n", "", -1)
	return s
}

// breakLines breaks a string into individual lines, considering both Linux and Windows line endings.
func breakLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		lines = append(lines, line)
	}
	return lines
}

// ParseSection handles the parsing of a section declaration line.
func (ini *MagicIni) Parse(source string) {
	buf := replaceLineScapes(source)
	lines := breakLines(buf)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == ';' || line[0] == '#' {
			continue
		}
		if line[0] == '[' && line[len(line)-1] == ']' {
			ini.ParseSection(line)
		} else if (line[0] == '[' || (line[0] == '/' && line[1] == '[')) && strings.Contains(line, "]") {
			ini.ParseSectionInlineKeyValue(line)
		} else {
			ini.ParseKeyValue(line)
		}
	}
}

func trimSectionKey(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "/[")
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	// s = strings.TrimSuffix(s, "\n")
	// s = strings.TrimSuffix(s, "\r")
	s = strings.TrimSpace(s)
	return s
}

func trimKey(s string) string {
	s = strings.TrimSpace(s)
	return s
}

func trimValue(s string) string {
	s = strings.TrimSpace(s)
	return s
}

// ParseSection handles the parsing of a section declaration line.
func (ini *MagicIni) ParseSection(line string) {
	line = trimSectionKey(line)
	if len(line) == 0 {
		return
	}
	if ini.Ini[line] == nil {
		ini.Ini[line] = make(map[string]string)
		ini.SectionsOrder = append(ini.SectionsOrder, line)
	}
	ini.currentParsedSection = line
}

// ParseSectionInlineKeyValue handles special case sections where a key-value pair is declared inline.
func (ini *MagicIni) ParseSectionInlineKeyValue(line string) {
	line = trimSectionKey(line)
	if len(line) == 0 {
		return
	}

	sectionAndKeyValue := strings.SplitN(line, "]", 2)
	if len(sectionAndKeyValue) > 1 {
		ini.currentParsedSection = trimSectionKey(sectionAndKeyValue[0])

		if ini.Ini[ini.currentParsedSection] == nil {
			ini.Ini[ini.currentParsedSection] = make(map[string]string)
			ini.SectionsOrder = append(ini.SectionsOrder, ini.currentParsedSection)
		}

		ini.ParseKeyValue(sectionAndKeyValue[1])
	}
}

func handleAsteriskValue(s string) string {
	if strings.HasPrefix(s, "*") {
		return strings.TrimPrefix(s, "*")
	}
	return s
}

// ParseKeyValue processes a single key-value pair line.
func (ini *MagicIni) ParseKeyValue(line string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}
	keyValue := strings.SplitN(line, "=", 2)
	if len(keyValue) > 1 {
		key := trimKey(keyValue[0])
		value := trimValue(keyValue[1])
		value = handleAsteriskValue(value)

		if ini.Ini[ini.currentParsedSection] == nil {
			ini.Ini[ini.currentParsedSection] = make(map[string]string)
			ini.SectionsOrder = append(ini.SectionsOrder, ini.currentParsedSection)
		}

		// Append key to order if it doesn't exist
		if _, exists := ini.Ini[ini.currentParsedSection][key]; !exists {
			ini.KeysOrder[ini.currentParsedSection] = append(ini.KeysOrder[ini.currentParsedSection], key)
		}

		ini.Ini[ini.currentParsedSection][key] = value
	}
}

// ParseKeyValue processes a single key-value pair line.
func (ini *MagicIni) GetSections() []string {
	var sections []string
	for k := range ini.Ini {
		sections = append(sections, k)
	}
	return sections
}

// GetSectionKeys retrieves all key names from a specified section.
func (ini *MagicIni) GetSectionKeys(section string) []string {
	var keys []string
	for k := range ini.Ini[section] {
		keys = append(keys, k)
	}
	return keys
}

// Get fetches the value associated with a given key within a specified section.
func (ini *MagicIni) Get(section string, key string) string {
	if ini.Ini[section] == nil {
		return ""
	}
	return ini.Ini[section][key]
}

// Translate attempts to replace placeholders in a string with their corresponding translations from the TranslatableSection.
func (ini *MagicIni) Translate(str string) string {
	re := regexp.MustCompile(`%(.*?)%`)
	return re.ReplaceAllStringFunc(str, func(s string) string {
		return ini.Translate(ini.Get(ini.TranslatableSection, strings.Trim(s, "%")))
	})
}

// Set assigns a value to a key within a specific section.
func (ini *MagicIni) Set(section string, key string, value string) {
	if ini.Ini[section] == nil {
		ini.Ini[section] = make(map[string]string)
	}
	ini.Ini[section][key] = value
}

// ParseSection handles the parsing of a section declaration line.
func (ini *MagicIni) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}
	defer f.Close()

	for _, section := range ini.SectionsOrder {
		_, err := f.WriteString("[" + section + "]\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}

		for _, key := range ini.KeysOrder[section] {
			value := ini.Ini[section][key]
			_, err := f.WriteString(key + "=" + value + "\n")
			if err != nil {
				return fmt.Errorf("error writing to file: %v", err)
			}
		}
		_, err = f.WriteString("\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
