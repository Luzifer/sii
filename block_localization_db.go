package sii

import "strings"

func init() {
	RegisterBlock(&LocalizationDB{})
}

type LocalizationDB struct {
	Keys []string `sii:"key"`
	Vals []string `sii:"val"`

	blockName string
}

func (LocalizationDB) Class() string { return "localization_db" }

func (l *LocalizationDB) Init(class, name string) {
	l.blockName = name
}

func (l LocalizationDB) Name() string { return l.blockName }

func (l LocalizationDB) GetTranslation(key string) string {
	key = strings.Trim(key, "@")

	for i, k := range l.Keys {
		if k == key {
			val := l.Vals[i]

			if strings.HasPrefix(val, "@@") {
				// Some translations are translation keys themselves
				return l.GetTranslation(val)
			}

			return val
		}
	}

	return ""
}
