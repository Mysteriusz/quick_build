package ar

/*	

FIELDS:
	'Mode' -> ar-compatbile mode ex: rcs
	'OutputExt' -> ar-compatbile extension of the output archive
	'OutputName' -> ar-compatbile name of the output archive

*/
type PolicyConfig struct{
	Mode		string 	`toml:"mode"`
	OutputExt	string 	`toml:"output_ext"`
	OutputName	string 	`toml:"output_name"`
}

