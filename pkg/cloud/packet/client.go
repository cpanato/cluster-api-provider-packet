package packet

import (
	"fmt"
	"os"
	"strings"

	"github.com/packethost/cluster-api-provider-packet/pkg/cloud/packet/util"
	"github.com/packethost/packngo"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	apiTokenVarName = "PACKET_API_KEY"
)

type PacketClient struct {
	*packngo.Client
}

// NewClient creates a new Client for the given Packet credentials
func NewClient(packetAPIKey string) *PacketClient {
	token := strings.TrimSpace(packetAPIKey)

	if token != "" {
		return &PacketClient{packngo.NewClientWithAuth("gardener", token, nil)}
	}

	return nil
}

func GetClient() (*PacketClient, error) {
	token := os.Getenv(apiTokenVarName)
	if token == "" {
		return nil, fmt.Errorf("env var %s is required", apiTokenVarName)
	}
	return NewClient(token), nil
}
func (p *PacketClient) GetDevice(machine *clusterv1.Machine) (*packngo.Device, error) {
	c, err := util.MachineProviderFromProviderConfig(machine.Spec.ProviderSpec)
	if err != nil {
		return nil, fmt.Errorf("Failed to process config for providerSpec: %v", err)
	}
	tag := util.GenerateMachineTag(string(machine.UID))
	return p.GetDeviceByTags(c.ProjectID, []string{tag})
}
func (p *PacketClient) GetDeviceByTags(project string, tags []string) (*packngo.Device, error) {
	devices, _, err := p.Devices.List(project, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving devices: %v", err)
	}
	// returns the first one that matches all of the tags
	for _, device := range devices {
		if util.ItemsInList(device.Tags, tags) {
			return &device, nil
		}
	}
	return nil, nil
}
