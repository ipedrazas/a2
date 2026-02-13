package devops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type AnsibleCheckTestSuite struct {
	suite.Suite
	check   *AnsibleCheck
	tempDir string
}

func (s *AnsibleCheckTestSuite) SetupTest() {
	s.check = &AnsibleCheck{}
	tempDir, err := os.MkdirTemp("", "ansible-test-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *AnsibleCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *AnsibleCheckTestSuite) TestID() {
	s.Equal("devops:ansible", s.check.ID())
}

func (s *AnsibleCheckTestSuite) TestName() {
	s.Equal("Ansible Configuration", s.check.Name())
}

func (s *AnsibleCheckTestSuite) TestRun_NoAnsibleFiles() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "No Ansible files found")
}

func (s *AnsibleCheckTestSuite) TestRun_AnsibleCfg() {
	ansibleCfg := filepath.Join(s.tempDir, "ansible.cfg")
	err := os.WriteFile(ansibleCfg, []byte(`[defaults]
inventory = ./inventory
remote_user = ansible`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_PlaybookYml() {
	playbookYml := filepath.Join(s.tempDir, "playbook.yml")
	err := os.WriteFile(playbookYml, []byte(`---
- name: Configure servers
  hosts: webservers
  tasks:
    - name: Install nginx
      ansible.builtin.apt:
        name: nginx
        state: present`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_PlaybookYaml() {
	playbookYaml := filepath.Join(s.tempDir, "playbook.yaml")
	err := os.WriteFile(playbookYaml, []byte(`---
- name: Deploy application
  hosts: all
  roles:
    - myrole`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_SiteYml() {
	siteYml := filepath.Join(s.tempDir, "site.yml")
	err := os.WriteFile(siteYml, []byte(`---
- name: Site setup
  hosts: localhost
  tasks:
    - name: Debug
      ansible.builtin.debug:
        msg: "Hello"`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_SiteYaml() {
	siteYaml := filepath.Join(s.tempDir, "site.yaml")
	err := os.WriteFile(siteYaml, []byte(`---
- hosts: all
  tasks:
    - name: Ping
      ansible.builtin.ping:`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_MainYml() {
	mainYml := filepath.Join(s.tempDir, "main.yml")
	err := os.WriteFile(mainYml, []byte(`---
- name: Main playbook
  hosts: all
  become: yes
  gather_facts: yes`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_RolesDirectory() {
	rolesDir := filepath.Join(s.tempDir, "roles")
	err := os.MkdirAll(rolesDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_AnsibleDirectory() {
	ansibleDir := filepath.Join(s.tempDir, "ansible")
	err := os.MkdirAll(ansibleDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_DotAnsibleDirectory() {
	dotAnsibleDir := filepath.Join(s.tempDir, ".ansible")
	err := os.MkdirAll(dotAnsibleDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_PlaybookInAnsibleDir() {
	ansibleDir := filepath.Join(s.tempDir, "ansible")
	err := os.MkdirAll(ansibleDir, 0755)
	s.Require().NoError(err)

	playbookYaml := filepath.Join(ansibleDir, "playbook.yaml")
	err = os.WriteFile(playbookYaml, []byte(`---
- name: Deploy
  hosts: all
  tasks:
    - name: Task
      ansible.builtin.debug:
        msg: "Deploying"`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_PlaybookInPlaybooksDir() {
	playbooksDir := filepath.Join(s.tempDir, "playbooks")
	err := os.MkdirAll(playbooksDir, 0755)
	s.Require().NoError(err)

	playbookYaml := filepath.Join(playbooksDir, "deploy.yml")
	err = os.WriteFile(playbookYaml, []byte(`---
- name: Deploy app
  hosts: webservers
  roles:
    - deploy`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_RoleInRolesDir() {
	roleDir := filepath.Join(s.tempDir, "roles", "myrole", "tasks")
	err := os.MkdirAll(roleDir, 0755)
	s.Require().NoError(err)

	mainYaml := filepath.Join(roleDir, "main.yml")
	err = os.WriteFile(mainYaml, []byte(`---
- name: Install package
  ansible.builtin.apt:
    name: mypackage
    state: present`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_DoesNotConfuseWithNonAnsibleYaml() {
	// Create a YAML file that's not an Ansible playbook
	configYaml := filepath.Join(s.tempDir, "config.yaml")
	err := os.WriteFile(configYaml, []byte(`database:
  host: localhost
  port: 5432`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	// Should not find ansible files since config.yaml doesn't have ansible markers
	s.NotContains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *AnsibleCheckTestSuite) TestRun_IgnoresHiddenDirs() {
	// Create .hidden directory with playbook (should be ignored)
	hiddenDir := filepath.Join(s.tempDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0755)
	s.Require().NoError(err)

	playbookYaml := filepath.Join(hiddenDir, "playbook.yaml")
	err = os.WriteFile(playbookYaml, []byte(`---
- name: Test
  hosts: all
  tasks: []`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.NotContains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_IgnoresVendorDirs() {
	// Create vendor directory with playbook (should be ignored)
	vendorDir := filepath.Join(s.tempDir, "vendor")
	err := os.MkdirAll(vendorDir, 0755)
	s.Require().NoError(err)

	playbookYaml := filepath.Join(vendorDir, "playbook.yaml")
	err = os.WriteFile(playbookYaml, []byte(`---
- name: Test
  hosts: all
  tasks: []`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.NotContains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_SkipsRequirementsAndGalaxyFiles() {
	// Create requirements and galaxy files which should be skipped
	requirementsYml := filepath.Join(s.tempDir, "requirements.yml")
	err := os.WriteFile(requirementsYml, []byte(`roles:
  - name: geerlingguy.nginx`), 0644)
	s.Require().NoError(err)

	galaxyYml := filepath.Join(s.tempDir, "galaxy.yml")
	err = os.WriteFile(galaxyYml, []byte(`roles:
  - name: common`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.NotContains(result.Reason, "ansible-lint")
}

func (s *AnsibleCheckTestSuite) TestRun_MultipleAnsibleIndicators() {
	// Create multiple indicators (should still work)
	ansibleCfg := filepath.Join(s.tempDir, "ansible.cfg")
	err := os.WriteFile(ansibleCfg, []byte(`[defaults]`), 0644)
	s.Require().NoError(err)

	playbookYaml := filepath.Join(s.tempDir, "site.yml")
	err = os.WriteFile(playbookYaml, []byte(`---
- hosts: all
  tasks: []`), 0644)
	s.Require().NoError(err)

	rolesDir := filepath.Join(s.tempDir, "roles")
	err = os.MkdirAll(rolesDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "ansible-lint")
}

func TestAnsibleCheckTestSuite(t *testing.T) {
	suite.Run(t, new(AnsibleCheckTestSuite))
}
