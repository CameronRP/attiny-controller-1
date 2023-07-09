module github.com/TheCacophonyProject/attiny-controller

go 1.20

require (
	github.com/TheCacophonyProject/event-reporter/v3 v3.4.0
	github.com/TheCacophonyProject/go-api v1.0.2
	github.com/TheCacophonyProject/go-config v1.8.3
	github.com/TheCacophonyProject/modemd v1.5.1
	github.com/TheCacophonyProject/window v0.0.0-20200312071457-7fc8799fdce7
	github.com/alexflint/go-arg v1.4.3
	github.com/c9s/goprocinfo v0.0.0-20190309065803-0b2ad9ac246b
	github.com/godbus/dbus v4.1.0+incompatible
	github.com/stretchr/testify v1.8.1
	golang.org/x/sys v0.6.0
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/host/v3 v3.8.2
)

require (
	github.com/alexflint/go-scalar v1.1.0 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/nathan-osman/go-sunrise v0.0.0-20171121204956-7c449e7c690b // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.9.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.64.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/TheCacophonyProject/event-reporter/v3 => ../event-reporter

replace github.com/TheCacophonyProject/go-config => ../go-config

replace github.com/TheCacophonyProject/go-api => ../go-api

replace github.com/TheCacophonyProject/modemd => ../modemd
