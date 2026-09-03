# Snyk MCP Tools

Tools exposed by the Snyk MCP Server. These are **not** shell commands: they are tool names an
MCP client (Claude Code, Cursor, VS Code, Windsurf) calls for you. You ask in plain language
("scan this repo for vulns") and the client picks the tool; the names below are what shows up
in the tool log, and what you name explicitly when you want a specific scan.

Setup: install the Snyk CLI, then register the server in your MCP client config with the
command `snyk mcp -t stdio` (older releases also need `--experimental`). Every tool takes an
absolute path to the folder or file to scan, plus an optional Snyk org id. First run of any
scan: `snyk_trust` the folder, then `snyk_auth`.

Typical flow: trust -> auth -> scan -> fix -> re-scan.

## Start the Snyk MCP server over stdio (shell command, goes in the client config)
`snyk mcp -t stdio`

The one real CLI command here. Put it as the `command` of an MCP server entry in your client
config. Use `-t sse` instead if the client wants an HTTP/SSE transport.

## Trust a folder so Snyk may scan it
`snyk_trust`

Run this first on any new repo. Snyk refuses to scan untrusted paths, so a scan that returns
"path not trusted" means you skipped this. Give it the absolute path of the project root, not
a subfolder.

## Authenticate with Snyk
`snyk_auth`

Opens a browser to log in, or uses an API token if one is configured. Needed once per machine;
the token is cached locally. Every scan tool fails with an auth error until this succeeds.

## Log out and clear the stored token
`snyk_logout`

Deletes the locally stored auth token. Use it to switch accounts or before handing over a
shared machine. Re-run `snyk_auth` afterwards.

## SAST scan of source code (no execution)
`snyk_code_scan`

Static analysis of your own first-party code — injection, path traversal, hardcoded secrets,
weak crypto. Point it at the repo root to scan everything, or at one file/subfolder to check
just what you touched. Fast, so it is the one to re-run after applying a fix. Covers Apex,
C/C++, Dart/Flutter, Elixir, Go, Groovy, Java/Kotlin, JavaScript, .NET, PHP, Python, Ruby,
Rust, Scala, Swift/Objective-C, TypeScript, VB.NET.

## SCA scan of dependency manifests
`snyk_sca_scan`

Third-party dependencies, not your code. Reads the manifests it finds (package.json, pom.xml,
requirements.txt, go.mod, Gemfile.lock, ...) and queries the Snyk vuln DB. Returns the
vulnerable package, the path that pulls it in, and the version that fixes it — ask the client
to bump those versions and re-run. A monorepo needs the all-projects behavior to catch
manifests in subfolders; a lockfile present means much more precise results.

## Scan a container image
`snyk_container_scan`

Give it an image reference (`nginx:1.25`, `myapp:latest`, a registry path). Reports vulns in
the base-image OS packages plus app dependencies baked into the layers. The useful output is
the base-image upgrade suggestion — usually one FROM line fixes most findings. Pass the
Dockerfile path too when you have it, so it can point at the offending instruction. The image
must exist locally or be pullable with your current docker credentials.

## Scan Infrastructure as Code for misconfigurations
`snyk_iac_scan`

Config, not code and not dependencies: public S3 buckets, wide-open security groups,
privileged containers, missing encryption. Point it at the directory holding your manifests.
Reads Terraform (.tf and plan JSON), Kubernetes YAML/JSON, AWS CloudFormation, Azure Resource
Manager JSON, and Serverless Framework. Each finding names the file, the resource, and the
line to change.

## Scan an existing SBOM for known vulnerabilities (experimental)
`snyk_sbom_scan`

For when you have an SBOM but not the source — a vendor artifact, a build output, an air-gapped
handoff. Give it the path to a CycloneDX or SPDX document. Every component must carry a
PackageURL (purl); components without one are silently skipped, so an empty result usually
means a purl-less SBOM, not a clean one.

## Generate AI-BOM for a Python project (CycloneDX v1.6 JSON, experimental)
`snyk_aibom`

Inventories the AI side of a Python codebase: models, datasets, frameworks and other AI
components, emitted as CycloneDX v1.6 JSON. Run it from inside the Python project. Needs
internet, the feature enabled on your account, and the preview-channel CLI — on a stable CLI
it just errors out.

## Send feedback to Snyk
`snyk_send_feedback`

Sends a free-text message to Snyk from inside the session. For false positives or missing
coverage. Do not paste code or secrets into it.

## Show the installed Snyk MCP Server version
`snyk_version`

Prints the server version. First thing to check when a tool is missing or a flag behaves
differently than the docs say — the experimental tools move fast.

## Examples
```
snyk mcp -t stdio
snyk_trust
snyk_auth
snyk_code_scan
snyk_sca_scan
snyk_iac_scan
snyk_container_scan
snyk_sbom_scan
snyk_aibom
snyk_version
snyk_logout
snyk_send_feedback
```
