package state

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// terraformState is JSON representation of terraform show command.
// https://developer.hashicorp.com/terraform/internals/json-format
type terraformState struct {
	Values terraformValues `json:"values"`
}

type terraformValues struct {
	Module terraformModule `json:"root_module"`
}

type terraformModule struct {
	Resources []terraformResource `json:"resources"`
}

type terraformResource struct {
	Type         string `json:"type"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Values       map[string]interface{}
}

type TerraformProvider struct {
	data *terraformState
}

func NewTerraformProvider() (*TerraformProvider, error) {
	showCmd := exec.Command("terraform", "show", "-json")

	rawState, err := showCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("'terraform show -json' command failed: %s", err)
	}

	var tfState terraformState
	if err := json.Unmarshal(rawState, &tfState); err != nil {
		return nil, err
	}

	return &TerraformProvider{
		data: &tfState,
	}, nil
}

func NewTerraformProviderFromStateOutput(path string) (*TerraformProvider, error) {
	rawState, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tfState terraformState
	if err := json.Unmarshal(rawState, &tfState); err != nil {
		return nil, err
	}

	return &TerraformProvider{
		data: &tfState,
	}, nil
}

func (tf *TerraformProvider) Contains(r *Resource) bool {
	for _, resource := range tf.data.Values.Module.Resources {
		if resource.Type == r.Type && resource.Values[r.IDKey].(string) == r.IDValue {
			return true
		}
	}

	return false
}
