package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/kirecek/codonly/internal/pkg/state"
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

type tfStateWrapper struct {
	data *terraformState
}

// NewTerraformProvider uses terraform binary to read terraform state into typed structure.
func New() (state.State, error) {
	showCmd := exec.Command("terraform", "show", "-json")

	rawState, err := showCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("'terraform show -json' command failed: %s", err)
	}

	var tfState terraformState
	if err := json.Unmarshal(rawState, &tfState); err != nil {
		return nil, err
	}

	return &tfStateWrapper{
		data: &tfState,
	}, nil
}

// NewTerraformProviderFromStateOutput reads terraform json state from a given file path.
func NewFromStateOutput(path string) (state.State, error) {
	rawState, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tfState terraformState
	if err := json.Unmarshal(rawState, &tfState); err != nil {
		return nil, err
	}

	return &tfStateWrapper{
		data: &tfState,
	}, nil
}

// Contains checks if a given Resource exists in terraform state.
func (tf *tfStateWrapper) Contains(r *state.Resource) bool {
	for _, resource := range tf.data.Values.Module.Resources {
		key := "id"
		if r.IDKey != nil {
			key = *r.IDKey
		}
		if resource.Type == r.Type && resource.Values[key].(string) == r.IDValue {
			return true
		}
	}

	return false
}
