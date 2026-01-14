# Plan

This document tracks proposed new changes for a2.

### Priority Levels
- **P0 (High)**: Critical for a2 adoption
- **P1 (Medium)**: Important for a2 adoption
- **P2 (Low)**: Nice-to-have for a2 adoption


## Languages to support in a2

1. Java [P0] ✅ 
2. Rust [P1] ✅
3. Typescript [P0] ✅
4. Ruby [P2]
5. Swift [P2] ✅ 


## Tasks

- New command `run` to run a particular check. For example, when we do `a2 check` we might get this result: 
```
! WARN Go Race Detection (2.4s)
    Race condition detected
```
    Good because it tells us what's wrong, but not enough details.
- New command `explain` to provide details on what a particular check does. something like`a2 explain CHECK`. We could also add a flag to `a2 list checks --explain` to list all the checks with a more detailed explanation of what it does. 
- Validation command: a2 profiles validate / a2 targets validate to check user definitions
- Add DevOps checks: ansible, helm, terraform, pulumi...
- Does it make sense to categorise the checks? for example, pre-commit hooks, editorconfig is local dev, k8s, retry, telemtery, signals is prod-ready... 
- Add the check name next to ` PASS Go Build (1.4s)` as ` PASS Go Build (1.4s) - go:build` so we know the name of the check.
- Add a `--verbose` flag to add for info to the output

