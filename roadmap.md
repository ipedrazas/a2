# Roadmap

## Disabled vs Not applicable

Right now, we don't have a way of differntiating between a check we don't want to run from a check it makes no sense to run.

Does it make sense to have a list of the checks we want to run instead?


## security:network: plan for “allowlisted egress” and make it UX-friendly

Support whitelisting, you’ll want a config that is:

- explicit
- reviewable in PRs
- not too granular (or it becomes unmaintainable)

Suggestion for allowlist structure (conceptual)

- Allow by:
	- domain (preferred): api.openai.com, github.com, registry.npmjs.org
	- protocol/port constraints (https:443)
	- optionally path prefixes for especially sensitive domains

- Avoid raw IP allowlists unless you truly hit fixed IPs; domains are more stable.


## Firecraker execution (WIP)

When running `a2` as a web service, we need to isolate the execution since it's untrusted code.

Is essence, the web service should get the github repo, send it to a queue that a firecracker executor consumes and runs it in a microvm:

- git clone
- a2 check -v --f json

