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

- We got profiles wrong. We want to have profiles like `desktop, API, library, cli` not prod, dev or poc. Those are targets more than profiles. So, we need to rename Profiles to Targets, and define profiles properly. The idea is that we can do `a2 check --profile desktop` and run the checks associated with desktop, which means that we will have to have a way of defining what a `desktop` profile should be.
- Make profiles configurable. Users should be able to define what an `API` profile looks like. Let's define the default profiles as external configuration files that the users can modify if they want to to match their expectations.