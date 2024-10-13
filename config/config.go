package config

var Config FullConfig

type FullConfig struct {
	General GeneralConfig
}

type GeneralConfig struct {
	NoteLocation string
}
