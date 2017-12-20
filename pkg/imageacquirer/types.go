package imageacquirer

import (
	docker "github.com/fsouza/go-dockerclient"
	iiapi "github.com/openshift/image-inspector/pkg/api"
	iicmd "github.com/openshift/image-inspector/pkg/cmd"
)

// ImageAcquirer abstract getting an image and extracting it in a given directory
type ImageAcquirer interface {
	// Acquire gets the image from `source` and extract it in `dest` which is the first output
	Acquire(source string) (string, docker.Image, iiapi.ScanResult, iiapi.FilesFilter, error)
}

// AuthsOptions is a struct to hold options to retrieve authentication credentials to registries and services.
type AuthsOptions struct {
	// DockerCfg is the location of the docker config file.
	DockerCfg iicmd.MultiStringVar
	// Username is the username for authenticating to the docker registry.
	Username string
	// PasswordFile is the location of the file containing the password for authentication to the
	// docker registry.
	PasswordFile string
}

func NewDockerImageAcquirer(dockerSocket string,
	preferedDestination string,
	pullPolicy string,
	auths AuthsOptions) ImageAcquirer {
	return &dockerImageAcquirer{dockerSocket, preferedDestination, pullPolicy, auths}
}

func NewContainerLibImageAcquirer(dstPath string, registryCertPath string, auths AuthsOptions) ImageAcquirer {
	return &containerLibImageAcquirer{dstPath, registryCertPath, auths}
}

func NewDockerContainerImageAcquirer(dockerSocket string, scanContainerChanges bool) ImageAcquirer {
	return &dockerContainerImageAcquirer{dockerSocket, scanContainerChanges}
}
