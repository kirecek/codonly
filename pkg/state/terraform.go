package state

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type TerraformState struct {
	Values TerraformValues `json:"values"`
}

type TerraformValues struct {
	Module TerraformModule `json:"root_module"`
}

type TerraformModule struct {
	Resources []TerraformResource `json:"resources"`
}

type TerraformResource struct {
	Type         string `json:"type"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Values       map[string]interface{}
}

type TerraformProvider struct {
	Data *TerraformState
}

func NewTerraformProvider() (*TerraformProvider, error) {
	showCmd := exec.Command("terraform", "show", "-json")

	rawState, err := showCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("'terraform show -json' command failed: %s", err)
	}

	var tfState TerraformState
	if err := json.Unmarshal(rawState, &tfState); err != nil {
		return nil, err
	}

	return &TerraformProvider{
		Data: &tfState,
	}, nil
}

func NewTerraformProviderFromStateOutput(path string) (*TerraformProvider, error) {
	rawState, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tfState TerraformState
	if err := json.Unmarshal(rawState, &tfState); err != nil {
		return nil, err
	}

	return &TerraformProvider{
		Data: &tfState,
	}, nil
}
