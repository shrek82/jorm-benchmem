---
name: "update-jorm"
description: "Updates the jorm dependency to the latest version. Invoke when user asks to update jorm, provides a jorm go get command, or mentions jorm version upgrades."
---

# Update Jorm Skill

This skill automates the process of updating the `github.com/shrek82/jorm` dependency to its latest version.

## Capabilities

- Detects when the user wants to update `jorm`.
- Identifies the latest available version of `jorm`.
- Updates the `go.mod` file to use the latest version.
- runs `go mod tidy` to ensure consistency.

## Usage

When the user asks to "update jorm" or shows a command like `go get -u github.com/shrek82/jorm@v1.0.0-alpha.6`, follow these steps:

1. **Check Latest Version**:
   - Run `go list -m -versions github.com/shrek82/jorm` to see available versions if needed, or simply trust `@latest`.

2. **Execute Update**:
   - Run the following commands:
     ```bash
     go get -u github.com/shrek82/jorm@latest
     go mod tidy
     ```

3. **Verify and Report**:
   - Check the `go.mod` file or the output of `go get` to confirm the new version.
   - Inform the user which version was installed (e.g., "Updated jorm from v1.0.0-alpha.6 to v1.0.0-alpha.7").

## Example Interaction

**User**: "Help me update jorm" or "go get ... jorm@v1.0.0-alpha.6"
**Agent**: "I will update jorm to the latest version for you."
(Runs `go get -u github.com/shrek82/jorm@latest` && `go mod tidy`)
**Agent**: "Successfully updated jorm to v1.0.0-alpha.7."
