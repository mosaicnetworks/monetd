package config

type configFile struct {
	label          string
	sourcefilename string
	targetfilename string
	subfolder      string
	required       bool
	transformation bool
}
