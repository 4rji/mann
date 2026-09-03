# Snyk MCP Tools

Tools exposed by the Snyk MCP Server. Each entry is the tool name to call from an MCP client (Claude Code, Cursor, etc.), not a shell command.

## Generate AI-BOM for a Python project (CycloneDX v1.6 JSON, experimental)
`snyk_aibom`

Lists AI models, datasets and tools used. Needs internet, feature access, preview-channel CLI, and must run inside a Python project.

## Authenticate with Snyk
`snyk_auth`

## Log out and clear the stored token
`snyk_logout`

## SAST scan of source code (no execution)
`snyk_code_scan`

Apex, C/C++, Dart/Flutter, Elixir, Go, Groovy, Java/Kotlin, JavaScript, .NET, PHP, Python, Ruby, Rust, Scala, Swift/Objective-C, TypeScript, VB.NET.

## SCA scan of dependency manifests
`snyk_sca_scan`

Reads package.json, pom.xml, requirements.txt and friends; reports open-source vulns and license issues.

## Scan a container image
`snyk_container_scan`

OS packages plus application dependencies.

## Scan Infrastructure as Code for misconfigurations
`snyk_iac_scan`

Terraform, Kubernetes YAML/JSON, AWS CloudFormation, Azure Resource Manager JSON, Serverless Framework.

## Scan an existing SBOM for known vulnerabilities (experimental)
`snyk_sbom_scan`

Components must carry PackageURLs (purls).

## Trust a folder so Snyk may scan it
`snyk_trust`

## Send feedback to Snyk
`snyk_send_feedback`

## Show the installed Snyk MCP Server version
`snyk_version`

## Examples
```
snyk_trust
snyk_auth
snyk_code_scan
snyk_sca_scan
snyk_iac_scan
snyk_container_scan
```
