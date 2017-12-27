package imageacquirer

import (
	iiapi "github.com/openshift/image-inspector/pkg/api"
	iicmd "github.com/openshift/image-inspector/pkg/cmd"
)

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
	auths AuthsOptions) iiapi.ImageAcquirer {
	return &dockerImageAcquirer{dockerSocket, preferedDestination, pullPolicy, auths}
}

func NewDockerContainerImageAcquirer(dockerSocket string, scanContainerChanges bool) iiapi.ImageAcquirer {
	return &dockerContainerImageAcquirer{dockerSocket, scanContainerChanges}
}
