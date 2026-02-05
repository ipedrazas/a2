package godogs

import (
	"sync"
)

// State holds the test state across scenarios
type State struct {
	mu              sync.RWMutex
	tempDir         string // temporary directory for scenario (created in Before, removed in After)
	originalDir     string // working directory before scenario (restored in After)
	projectDir      string
	a2Installed     bool
	configFile      string
	lastCommand     string
	lastOutput      string
	lastExitCode    int
	currentLanguage string
	maturityScore   int
	checkResults    map[string]string
	issuesDetected  []string
	beforeMaturity  int // for comparing scores
}

var (
	sharedState *State
	once        sync.Once
)

// GetState returns the shared state instance
func GetState() *State {
	once.Do(func() {
		sharedState = &State{
			checkResults:   make(map[string]string),
			issuesDetected: make([]string, 0),
		}
	})
	return sharedState
}

// ResetState resets the state for a new scenario
func ResetState() {
	s := GetState()
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tempDir = ""
	s.originalDir = ""
	s.projectDir = ""
	s.a2Installed = false
	s.configFile = ""
	s.lastCommand = ""
	s.lastOutput = ""
	s.lastExitCode = 0
	s.currentLanguage = ""
	s.maturityScore = 0
	s.beforeMaturity = 0
	s.checkResults = make(map[string]string)
	s.issuesDetected = make([]string, 0)
}

// Setters
func (s *State) SetTempDir(dir string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tempDir = dir
}

func (s *State) SetOriginalDir(dir string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.originalDir = dir
}

func (s *State) SetProjectDir(dir string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projectDir = dir
}

func (s *State) SetA2Installed(installed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.a2Installed = installed
}

func (s *State) SetConfigFile(file string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.configFile = file
}

func (s *State) SetLastCommand(cmd string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastCommand = cmd
}

func (s *State) SetLastOutput(output string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastOutput = output
}

func (s *State) SetLastExitCode(code int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastExitCode = code
}

func (s *State) SetCurrentLanguage(lang string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentLanguage = lang
}

func (s *State) SetMaturityScore(score int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maturityScore = score
}

func (s *State) SetBeforeMaturity(score int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.beforeMaturity = score
}

func (s *State) AddCheckResult(check, result string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkResults[check] = result
}

func (s *State) AddIssue(issue string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.issuesDetected = append(s.issuesDetected, issue)
}

// Getters
func (s *State) GetTempDir() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tempDir
}

func (s *State) GetOriginalDir() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.originalDir
}

func (s *State) GetProjectDir() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projectDir
}

func (s *State) GetA2Installed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.a2Installed
}

func (s *State) GetConfigFile() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.configFile
}

func (s *State) GetLastCommand() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastCommand
}

func (s *State) GetLastOutput() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastOutput
}

func (s *State) GetLastExitCode() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastExitCode
}

func (s *State) GetCurrentLanguage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentLanguage
}

func (s *State) GetMaturityScore() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maturityScore
}

func (s *State) GetBeforeMaturity() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.beforeMaturity
}

func (s *State) GetCheckResult(check string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.checkResults[check]
	return result, ok
}

func (s *State) GetIssuesDetected() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.issuesDetected
}
