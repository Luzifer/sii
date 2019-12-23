package main

var userConfigPath = "~/.config/sii-editor/config.yml"
var userCachePath = "~/.cache/sii-editor"

var profilePaths = []string{
	// Linux default for non-Steam-profiles
	"~/.local/share/Euro Truck Simulator 2/profiles",
}

var gamePaths = []string{
	// Linux default installation path for Steam games
	"~/.local/share/Steam/steamapps/common/Euro Truck Simulator 2",
}
