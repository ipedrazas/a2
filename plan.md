# Plan

This document tracks proposed new languages for a2.

### Priority Levels
- **P0 (High)**: Critical for a2 adoption
- **P1 (Medium)**: Important for a2 adoption
- **P2 (Low)**: Nice-to-have for a2 adoption


## Languages to support in a2

1. Java [P0] ✅ 
2. Rust [P1] ✅
3. Typescript [P0] ✅
4. Ruby [P2]
5. Swift [P2]


## Tasks

1. Project-level profiles: Support .a2/profiles/ in project directory (not in this iteration)
2. Profile/target inheritance: Allow extends: cli to inherit from another profile/target
3. Validation command: a2 profiles validate / a2 targets validate to check user definitions
4. Environment-based loading: Support A2_CONFIG_DIR environment variable to override config location
5. Add a command `a2 init .` to create a `.a2.yaml`, we could have a flag to indicate if we want to initialise profiles and targets too.