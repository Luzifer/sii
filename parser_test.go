package sii

import (
	"strings"
	"testing"
)

var testSii = `SiiNunit
{
/*
 * Multi-line comment
 */
some_unit : .my_mod.unit
{
		// Single line comment
		# Single line comment
		/*
		 * In-Block multi-line string
		 */
    attribute_number: 40
    attribute_string: "TEST STRING"
    attribute_unit: test.unit
    attribute_float3: (1.0, 1.0, 1.0)
    attribute_float_number_ieee754: &40490f5a
}

save_container : _nameless.1c57,b4b0 {
 name: ""
 time: 96931
 file_time: 1572907597
 version: 42
 dependencies: 14
 dependencies[0]: "mod|promods-assets-v242|ProMods Assets Package"
 dependencies[1]: "mod|promods-model1-v242|ProMods Models Package 1"
 dependencies[2]: "mod|promods-model2-v242|ProMods Models Package 2"
 dependencies[3]: "mod|promods-model3-v242|ProMods Models Package 3"
 dependencies[4]: "mod|promods-media-v242|ProMods Media Package"
 dependencies[5]: "mod|promods-map-v242|ProMods Map Package"
 dependencies[6]: "mod|promods-def-v242|ProMods Definition Package"
 dependencies[7]: "dlc|eut2_balt|DLC - Beyond the Baltic Sea"
 dependencies[8]: "dlc|eut2_east|DLC - Going East!"
 dependencies[9]: "dlc|eut2_fr|DLC - Vive la France !"
 dependencies[10]: "dlc|eut2_it|DLC - Italia"
 dependencies[11]: "rdlc|eut2_metallics|DLC - Metallic Paint Jobs"
 dependencies[12]: "dlc|eut2_north|DLC - Scandinavia"
 dependencies[13]: "rdlc|eut2_rocket_league|DLC - Rocket League"
}
}
`

func TestParseUnit(t *testing.T) {
	unit, err := parseSIIPlainFile(strings.NewReader(testSii))
	if err != nil {
		t.Fatalf("parseSIIPlainFile caused an error: %s", err)
	}

	if len(unit.Entries) != 2 {
		t.Errorf("Expected 1 block, got %d", len(unit.Entries))
	}

	t.Logf("%#v", unit)
	t.Logf("%#v", unit.Entries[0])
	t.Logf("%#v", unit.Entries[1])
}
