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


- Add DevOps checks: ansible, helm, terraform, pulumi...
- Does it make sense to categorise the checks? for example, pre-commit hooks, editorconfig is local dev, k8s, retry, telemtery, signals is prod-ready... 
- Web needs to run the checks in verbose mode.
