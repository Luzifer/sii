package sii

func init() {
	RegisterBlock(&SaveContainer{})
}

type SaveContainer struct {
	SaveName     string   `sii:"name"`
	Time         int64    `sii:"time"`
	FileTime     int64    `sii:"file_time"`
	Version      int      `sii:"version"`
	Dependencies []string `sii:"dependencies"`

	blockName string
}

func (SaveContainer) Class() string { return "save_container" }

func (s *SaveContainer) Init(class, name string) {
	s.blockName = name
}

func (s SaveContainer) Name() string { return s.blockName }
