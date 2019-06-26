package config

type configFile struct {
	label          string
	sourcefilename string
	targetfilename string
	required       bool
	transformation bool
}

type configFiles []*configFile
