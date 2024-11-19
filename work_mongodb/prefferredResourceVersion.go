package work

import (
	"fmt"

	"k8s.io/client-go/discovery"
	kmapi "kmodules.xyz/client-go/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func GetPreferredResourceVersion(ref kmapi.TypedObjectReference) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	c, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return "", err
	}

	grouplist, err := c.ServerGroups()
	if err != nil {
		return "", err
	}

	var preferredVersion string
	for _, group := range grouplist.Groups {
		if group.Name != ref.APIGroup {
			continue
		}

		var supportedVersions []string
		for _, ver := range group.Versions {
			apiResources, err := c.ServerResourcesForGroupVersion(ver.GroupVersion)
			if err != nil {
				continue
			}

			for _, resource := range apiResources.APIResources {
				if resource.Kind == ref.Kind {
					supportedVersions = append(supportedVersions, ver.Version)
					break
				}
			}
		}

		if len(supportedVersions) == 0 {
			continue
		}

		preferredVersion = supportedVersions[0]
		for _, ver := range supportedVersions {
			if ver == group.PreferredVersion.Version {
				preferredVersion = ver
				break
			}
		}
		if preferredVersion != "" {
			return preferredVersion, nil
		}
	}

	return "", fmt.Errorf("unable to find resource version for %s/%s", ref.APIGroup, ref.Kind)
}
