# AI Usage Report

**Date:** July 29, 2025

## Overview

This report documents the use of AI tools, models, and supporting infrastructure in the development of the `adowork` CLI project. This report itself was AI-generated, but reviewed and updated by the human author.

## AI Tools and Systems Used

### GitHub Copilot

- **Role:** Code generation, refactoring, and suggestions directly in the editor.
- **Models:**
  - GPT-4.1
  - Claude Sonnet 4
  - Gemmini 2.5 Pro

### Taskmaster

- **Purpose:** Task-driven project management, task breakdown, and workflow automation.
- **Usage:** Managed all project tasks, subtasks, and status transitions. Enabled structured, iterative development and clear tracking of progress.
- **Features Leveraged:**
  - Task expansion and complexity analysis
  - Automated status updates and reporting
  - Tag-based task contexts for workflow isolation

#### Taskmaster AI Models (from .taskmaster/config.json)

- **Main Model:**
  - Provider: Anthropic
  - Model ID: `claude-sonnet-4-20250514`
  - Max Tokens: 64,000
  - Temperature: 0.2
- **Research Model:**
  - Provider: Perplexity
  - Model ID: `sonar-pro`
  - Max Tokens: 8,700
  - Temperature: 0.1
- **Fallback Model:**
  - Provider: Anthropic
  - Model ID: `claude-3-7-sonnet-20250219`
  - Max Tokens: 120,000
  - Temperature: 0.2

### MCP (Model Context Protocol) Servers

- **MCP Endpoints Used:**
  - Local MCP servers (VSCode GitHub Copilot extension integration)
  - No remote MCP endpoints were configured; all requests routed through the local development environment
- **MCP Servers Used:**
  - **time**: Time zone conversion and current time server
  - **brave**: Brave Search API integration for web search
  - **perplexity-ask**: Perplexity AI research and question-answering server
  - **context7**: Context7 documentation and library reference server
  - **taskmaster-ai**: Taskmaster's core AI functionality server (npm package)
- **API Integration:**
  - Anthropic API (via secure input prompts)
  - Perplexity API (via secure input prompts)
  - Brave Search API (via secure input prompts)

## Other tools

### yek

- **Purpose:** A fast Rust-based tool to serialize text-based files in a repository or directory for LLM consumption
- **Usage in Project:** Generated the [azure-devops-go-api.txt](azure-devops-go-api.txt) file by serializing the Microsoft Azure DevOps Go API library repository
