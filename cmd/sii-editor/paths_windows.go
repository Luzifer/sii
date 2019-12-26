package main

var userConfigPath = `~\documents\sii-editor\config.yml`

var profilePaths = map[string]string{
	// Windows default for non-Steam-profiles
	pathATS:  `~\documents\American Truck Simulator\profiles`,
	pathETS2: `~\documents\Euro Truck Simulator 2\profiles`,
}

var gamePaths = map[string]string{
	// Windows default installation path for Steam games
	pathATS:  `C:\Program Files (x86)\Steam\steamapps\common\American Truck Simulator`,
	pathETS2: `C:\Program Files (x86)\Steam\steamapps\common\Euro Truck Simulator 2`,
}
