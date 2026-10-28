# bjira

> [!WARNING]
> **Work in Progress**

A lightweight command-line interface for Jira Cloud, written in Go.

Manage and track Jira issues, sprints, and tasks directly from your terminal without using the web interface.

## Build

Requires Go installed.

```bash
go build -o bin/bjira ./src
```

## Configuration

Create a configuration file:

- **Linux:** `~/.config/bjira/cfg`
- **Windows:** `cfg` next to `bjira.exe`

```env
JIRA_URL=https://your-company.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-api-token
```

Generate your API token at [Atlassian API Tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

**Keep the configuration file private. Never commit API tokens to Git.**

## Usage

```bash
# Example usage: 
bjira ls -sprint="MW Sprint 33" -user="all" -status="To Do"
```
