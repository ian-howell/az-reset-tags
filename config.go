package main

// Config represents the configuration structure for nc-devbox.
type Config struct {
	TimeToLive        string      `yaml:"time_to_live"`
	Azure             AzureConfig `yaml:"azure"`
	VM                VMConfig    `yaml:"vm"`
	SSH               SSHConfig   `yaml:"ssh"`
	Sync              []string    `yaml:"sync"`
	PostCreateScripts []string    `yaml:"post_create_scripts"`
}

// AzureConfig represents the Azure-specific configuration.
type AzureConfig struct {
	ResourceGroup string `yaml:"resource_group"`
	Location      string `yaml:"location"`
}

// VMConfig represents the virtual machine configuration.
type VMConfig struct {
	Name               string `yaml:"name"`
	Size               string `yaml:"size"`
	Image              string `yaml:"image"`
	DiskSize           int    `yaml:"disk_size"`
	AutoShutdownTime   string `yaml:"auto_shutdown_time"`
	CloudInitPath      string `yaml:"cloud_init_path"`
	EnableSecurityTags bool   `yaml:"enable_security_tags"`
}

// SSHConfig represents the SSH configuration.
type SSHConfig struct {
	PrivateKeyPath string `yaml:"private_key_path"`
}
