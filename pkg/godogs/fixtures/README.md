# Godogs Test Fixtures

Fixture data for BDD scenarios.

## Structure

- **simple-go-project/** – Minimal Go project (go.mod, main.go, README.md) for quick-start and “clean” runs
- **with-issues/** – Go project with a failing test (main_test.go) so A2 reports FAIL; used by “Fix and re-run checks”
- **passing-go-project/** – Go project with .a2.yaml so a2 check passes (100%); used by "Quick pre-push validation"
- **multi-language-project/** – Placeholder for multi-language test projects
- **config-examples/** – Sample .a2.yaml and configs (api-profile, etc.)
- **expected-outputs/** – Expected test outputs for comparison (JSON, TOON)

Used by BeforeScenario temp dir and step helpers; scenarios copy these into the scenario temp dir as needed.
