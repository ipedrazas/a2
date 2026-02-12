package devops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type TerraformCheckTestSuite struct {
	suite.Suite
	check   *TerraformCheck
	tempDir string
}

func (s *TerraformCheckTestSuite) SetupTest() {
	s.check = &TerraformCheck{}
	tempDir, err := os.MkdirTemp("", "terraform-test-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *TerraformCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *TerraformCheckTestSuite) TestID() {
	s.Equal("devops:terraform", s.check.ID())
}

func (s *TerraformCheckTestSuite) TestName() {
	s.Equal("Terraform Configuration", s.check.Name())
}

func (s *TerraformCheckTestSuite) TestRun_NoTerraformFiles() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Terraform files found")
}

func (s *TerraformCheckTestSuite) TestRun_MainTfFile() {
	mainTf := filepath.Join(s.tempDir, "main.tf")
	err := os.WriteFile(mainTf, []byte(`
resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}
`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	// Without terraform installed, should pass with Info
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_VariablesTfFile() {
	varsTf := filepath.Join(s.tempDir, "variables.tf")
	err := os.WriteFile(varsTf, []byte(`
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}
`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_OutputTfFile() {
	outputTf := filepath.Join(s.tempDir, "outputs.tf")
	err := os.WriteFile(outputTf, []byte(`
output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.example.id
}
`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_TerraformTfvars() {
	tfvars := filepath.Join(s.tempDir, "terraform.tfvars")
	err := os.WriteFile(tfvars, []byte(`
instance_type = "t2.micro"
environment  = "prod"
`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_VersionsTfFile() {
	versionsTf := filepath.Join(s.tempDir, "versions.tf")
	err := os.WriteFile(versionsTf, []byte(`
terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}
`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_MultipleTfFiles() {
	// Create multiple .tf files
	mainTf := filepath.Join(s.tempDir, "main.tf")
	err := os.WriteFile(mainTf, []byte(`resource "aws_s3_bucket" "example" {}`), 0644)
	s.Require().NoError(err)

	varsTf := filepath.Join(s.tempDir, "variables.tf")
	err = os.WriteFile(varsTf, []byte(`variable "bucket_name" {}`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_NestedTfFiles() {
	// Create nested directory structure
	modulesDir := filepath.Join(s.tempDir, "modules", "storage")
	err := os.MkdirAll(modulesDir, 0755)
	s.Require().NoError(err)

	storageTf := filepath.Join(modulesDir, "main.tf")
	err = os.WriteFile(storageTf, []byte(`resource "aws_s3_bucket" "example" {}`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_AutoTfvars() {
	autoTfvars := filepath.Join(s.tempDir, "dev.auto.tfvars")
	err := os.WriteFile(autoTfvars, []byte(`
region = "us-west-2"
`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "terraform")
}

func (s *TerraformCheckTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *TerraformCheckTestSuite) TestRun_IgnoresHiddenDirs() {
	// Create .terraform directory (should be ignored)
	dotTerraform := filepath.Join(s.tempDir, ".terraform")
	err := os.MkdirAll(dotTerraform, 0755)
	s.Require().NoError(err)

	// Create a .tf file in a hidden directory (should be ignored)
	hiddenDir := filepath.Join(s.tempDir, ".hidden")
	err = os.MkdirAll(hiddenDir, 0755)
	s.Require().NoError(err)

	hiddenTf := filepath.Join(hiddenDir, "test.tf")
	err = os.WriteFile(hiddenTf, []byte(`resource "aws_s3_bucket" "test" {}`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Terraform files found")
}

func (s *TerraformCheckTestSuite) TestRun_IgnoresVendorDirs() {
	// Create vendor directory (should be ignored)
	vendorDir := filepath.Join(s.tempDir, "vendor")
	err := os.MkdirAll(vendorDir, 0755)
	s.Require().NoError(err)

	vendorTf := filepath.Join(vendorDir, "test.tf")
	err = os.WriteFile(vendorTf, []byte(`resource "aws_s3_bucket" "test" {}`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Terraform files found")
}

func TestTerraformCheckTestSuite(t *testing.T) {
	suite.Run(t, new(TerraformCheckTestSuite))
}
