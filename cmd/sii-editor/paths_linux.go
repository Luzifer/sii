package main

var userConfigPath = "~/.config/sii-editor/config.yml"

var profilePaths = map[string]string{
	// Linux default for non-Steam-profiles
	pathATS:  "~/.local/share/American Truck Simulator/profiles",
	pathETS2: "~/.local/share/Euro Truck Simulator 2/profiles",
}

var gamePaths = map[string]string{
	// Linux default installation path for Steam games
	pathATS:  "~/.local/share/Steam/steamapps/common/American Truck Simulator",
	pathETS2: "~/.local/share/Steam/steamapps/common/Euro Truck Simulator 2",
}
