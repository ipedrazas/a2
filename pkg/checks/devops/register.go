package devops

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all DevOps check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	return []checker.CheckRegistration{
		{
			Checker: &TerraformCheck{},
			Meta: checker.CheckMeta{
				ID:          "devops:terraform",
				Name:        "Terraform Configuration",
				Description: "Validates Terraform configurations for syntax and best practices.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       950,
				Suggestion:  "Fix terraform validation errors",
			},
		},
		{
			Checker: &AnsibleCheck{},
			Meta: checker.CheckMeta{
				ID:          "devops:ansible",
				Name:        "Ansible Configuration",
				Description: "Validates Ansible playbooks and roles using ansible-lint.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       960,
				Suggestion:  "Fix ansible-lint issues",
			},
		},
		{
			Checker: &PulumiCheck{},
			Meta: checker.CheckMeta{
				ID:          "devops:pulumi",
				Name:        "Pulumi Configuration",
				Description: "Validates Pulumi infrastructure as code configurations.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       970,
				Suggestion:  "Fix pulumi validation errors",
			},
		},
		{
			Checker: &HelmCheck{},
			Meta: checker.CheckMeta{
				ID:          "devops:helm",
				Name:        "Helm Charts",
				Description: "Validates Helm charts using helm lint.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       980,
				Suggestion:  "Fix helm lint issues",
			},
		},
		{
			Checker: &K8sCheck{},
			Meta: checker.CheckMeta{
				ID:          "devops:k8s",
				Name:        "Kubernetes Ready",
				Description: "Checks for Kubernetes manifests or Helm charts for deployment.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1030,
				Suggestion:  "Add Kubernetes manifests for deployment",
			},
		},
	}
}
