package linode

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/linode/linodego"
	v1 "k8s.io/api/core/v1"
)

const providerIDPrefixLinode = "linode://"
const providerIDPrefixIBM = "ibm://"

type invalidProviderIDError struct {
	value string
}

func (e invalidProviderIDError) Error() string {
	return fmt.Sprintf("invalid provider ID %q", e.value)
}

func parseProviderID(providerID string) (int, error) {
	id, err := strconv.Atoi(strings.TrimPrefix(providerID, providerIDPrefixLinode))
	if err != nil {
		return 0, invalidProviderIDError{providerID}
	}
	return id, nil
}

func getLinodeIDforSatellite(node *v1.Node, client Client) (int, error) {
	for _, address := range node.Status.Addresses {
		if address.Type == v1.NodeHostName {
			ipaddress := strings.ReplaceAll(address.Address, "-", ".")
			instances, err := client.ListInstances(context.Background(), &linodego.ListOptions{Filter: "{\"ipv4\": \"" + ipaddress + "\"}"})
			if err != nil {
				return 0, fmt.Errorf("Error getting instances: %s", err.Error())
			}
			if len(instances) != 1 {
				return 0, fmt.Errorf("Error finding a single instance with IP address %s. Found %d instances", ipaddress, len(instances))
			}
			return instances[0].ID, nil
		}
	}
	return 0, fmt.Errorf("Error detecting Linode ID for IBM Cloud Satellite")
}
