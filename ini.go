package ini

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func getEOLForCurrentOS() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

var EOL = getEOLForCurrentOS()

type MagicIni struct {
	Ini                  map[string]map[string]string
	currentParsedSection string
}

func LoadIni(path string) *MagicIni {
	ini := new(MagicIni)
	ini.Ini = make(map[string]map[string]string)
	ini.LoadIni(path)
	return ini
}

func (ini *MagicIni) LoadIni(path string) {
	fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("unable to read file: %v", err)
	}
	defer f.Close()
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("unable to read file: %v", err)
		}
		if n == 0 {
			break
		}
		ini.Parse(string(buf))
	}
}

func (ini *MagicIni) Parse(buf string) {
	lines := strings.Split(buf, EOL)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == ';' || line[0] == '#' {
			continue
		}
		if line[0] == '[' {
			ini.ParseSection(line)
		} else {
			ini.ParseKeyValue(line)
		}
	}
}

func trimSectionKey(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
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

func (ini *MagicIni) ParseSection(line string) {
	line = trimSectionKey(line)
	if len(line) == 0 {
		return
	}
	ini.currentParsedSection = line
	ini.Ini[line] = make(map[string]string)
}

func (ini *MagicIni) ParseKeyValue(line string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}
	keyValue := strings.SplitN(line, "=", 2)
	if len(keyValue) > 1 {
		key := trimKey(keyValue[0])
		value := trimValue(keyValue[1])
		ini.Ini[ini.currentParsedSection][key] = value
	}
}

func (ini *MagicIni) GetSections() []string {
	var sections []string
	for k := range ini.Ini {
		sections = append(sections, k)
	}
	return sections
}

func (ini *MagicIni) GetSectionKeys(section string) []string {
	var keys []string
	for k := range ini.Ini[section] {
		keys = append(keys, k)
	}
	return keys
}

func (ini *MagicIni) Get(section string, key string) string {
	if ini.Ini[section] == nil {
		return ""
	}
	return ini.Ini[section][key]
}

func (ini *MagicIni) Translate(str string) string {
	re := regexp.MustCompile(`%(.*?)%`)
	return re.ReplaceAllStringFunc(str, func(s string) string {
		return ini.Translate(ini.Get("MAGIC_LOGICAL_NAMES", strings.Trim(s, "%")))
	})
}
