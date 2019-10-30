package plister

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	filePath = filepath.Join(os.TempDir(), "Info.plist")
	dict     = map[string]interface{}{
		"CFBundlePackageType":     "APPL",
		"CFBundleIconFile":        "icon.icns",
		"CFBundleDisplayName":     "Best App display name",
		"CFBundleExecutable":      "app_binary",
		"CFBundleName":            "BestApp",
		"CFBundleIdentifier":      "com.company.BestApp",
		"LSUIElement":             "NO",
		"LSMinimumSystemVersion":  "10.11",
		"NSHighResolutionCapable": true,
		"NSAppTransportSecurity": map[string]interface{}{
			"NSAllowsArbitraryLoads": true,
		},
		"CFBundleURLTypes": []map[string]interface{}{{
			"CFBundleTypeRole":   "Viewer",
			"CFBundleURLName":    "com.developer.testapp",
			"CFBundleURLSchemes": []interface{}{"testappscheme"},
		}, {
			"CFBundleTypeRole":   "Reader",
			"CFBundleURLName":    "com.developer.testapp",
			"CFBundleURLSchemes": []interface{}{"testappscheme-read"},
		}},
		"SliceDict": []map[string]interface{}{{
			"SliceDictKey1": "SliceDictVal1",
		}, {
			"SliceDictKey2": "SliceDictVal2",
		}, {
			"SliceDictKey3": []interface{}{"val1", "val2"},
		}},
		"Slice": []interface{}{"val1", "val2"},
		"":      "",
	}
)

func TestInfoPlist_Get(t *testing.T) {
	infoPlist := MapToInfoPlist(dict)
	expect := "Best App display name"
	got := infoPlist.Get("CFBundleDisplayName")
	if got != expect {
		t.Errorf("expect %s, got %s", expect, got)
	}

	if got = infoPlist.Get("bad key"); got != nil {
		t.Errorf("expect empty string, got '%s'", got)
	}
}

func TestInfoPlist_Set(t *testing.T) {
	ttData := []struct {
		key   string
		value interface{}
	}{
		{"LSUIElement", "YES"},
		{"testKey", "testValue"},
		{"NSHighResolutionCapable", true},
	}
	infoPlist := MapToInfoPlist(dict)
	for _, tt := range ttData {
		infoPlist.Set(tt.key, tt.value)
		got := infoPlist.Get(tt.key)
		if got != tt.value {
			t.Errorf("expect %s, got %s", tt.value, got)
		}
	}
}

func TestGenerateInfoPlist(t *testing.T) {
	if err := Fprint(ioutil.Discard, MapToInfoPlist(dict)); err != nil {
		t.Error(err)
	}
}

func TestGenerate(t *testing.T) {
	dict := MapToInfoPlist(dict)
	if err := Generate("", dict); err == nil {
		t.Errorf("error not received for empty file path")
	}

	if err := Generate(filePath, dict); err != nil {
		t.Error(err)
	}
	fp, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		fp.Close()
		os.Remove(filePath)
	}()
	fi, err := fp.Stat()
	if err != nil {
		t.Error(err)
	}
	if fi.Size() == 0 {
		t.Errorf("result file size equal 0 bytes")
	}
}

func TestGenerateFromMap(t *testing.T) {
	if err := GenerateFromMap(filePath, dict); err != nil {
		t.Error(err)
	}
	defer os.Remove(filePath)
}

func TestParse(t *testing.T) {
	_, err := Parse("testdata/test.plist")
	if err != nil {
		t.Error(err)
	}
}
