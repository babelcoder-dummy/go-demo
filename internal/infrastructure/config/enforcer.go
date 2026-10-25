package config

import (
	_ "embed"
	"fmt"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	stringadapter "github.com/casbin/casbin/v3/persist/string-adapter"
)

//go:embed casbin_model.conf
var modelText string

//go:embed casbin_policy.csv
var policyText string

func NewCasbinEnforcer() (*casbin.Enforcer, error) {
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, fmt.Errorf("failed to load casbin model: %w", err)
	}

	adapter := stringadapter.NewAdapter(policyText)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	return enforcer, nil
}
