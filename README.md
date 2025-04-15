# Snap

Snap is an opinionated, fun, and community-focused version control system written in Go. It supports core version control features like initializing a repo, committing files, and branching, but with some key differences.

## Core Features

- Written in Go
- Command-line interface (CLI) similar to Git
- All commits must start with a Snapmoji (e.g., :sparkles:, ✨)
- Built-in issue tracking (create, assign, close issues)
- Gamification system where users earn points for contributions
- Local-first design (like Git)

## Commands

### Core Version Control

- `snap init` – Initialize a Snap repository
- `snap add <file>` – Stage files
- `snap commit -m "<message>"` – Commit with Snapmoji validation
- `snap status` – Show working tree status
- `snap log` – View commit history

### Issue Tracking

- `snap issue new -t "<title>" -d "<description>"` – Create a new issue
- `snap issue list` – List all open issues
- `snap issue show <id>` – Show issue details
- `snap issue close <id>` – Close an issue
- `snap issue assign <id> <assignee>` – Assign an issue to a user

### Gamification

- `snap me` – Show user stats and contribution points
- `snap leaderboard` – Show top contributors in the repo

### Fun Commands

- `snap boom <message>` – Shortcut for quick commit (like git commit -am)
- `snap crackle` – Stylized commit log view
- `snap pop` – Undo last commit
- `snap vibe` – Show mood of repo based on recent commits/snapmojis

### Web Interface

- `snap web` – Launch local Snap dashboard
- `snap web --port 8888` – Custom port
- `snap web --open` – Automatically open browser

Once the web interface is running, you can access these features:
- Home – Repository overview and stats
- Commits – Browse all commits
- Issues – View and manage issues
- Users – See contributor stats
- Quest – View your assigned issues

## Getting Started

1. Build the application:
   ```
   go build
   ```

2. Initialize a new repository:
   ```
   ./snap init
   ```

3. Add files to the staging area:
   ```
   ./snap add <file>
   ```

4. Commit changes (must start with a Snapmoji):
   ```
   ./snap commit -m "✨ Initial commit" -a "username"
   ```

5. Check the status of your repository:
   ```
   ./snap status
   ```

## Snapmoji Support

All commits in Snap must start with a Snapmoji. Here are some commonly used Snapmojis:

- ✨ `:sparkles:` – Introduce new features
- 🐛 `:bug:` – Fix a bug
- 📚 `:books:` – Add or update documentation
- ♻️ `:recycle:` – Refactor code
- 🔧 `:wrench:` – Add or update configuration files
- ✅ `:white_check_mark:` – Add, update, or pass tests
- 🚀 `:rocket:` – Deploy stuff
- 💄 `:lipstick:` – Add or update the UI and style files
- 🔥 `:fire:` – Remove code or files
- 🚑 `:ambulance:` – Critical hotfix

## Project Structure

- `cmd/` - CLI commands
- `pkg/` - Core functionality
  - `pkg/repository/` - Repository management
  - `pkg/snapmoji/` - Snapmoji validation
  - `pkg/issue/` - Issue tracking
  - `pkg/user/` - User stats and gamification
  - `pkg/storage/` - Storage utilities

## Future Enhancements

- `snap review` – Request a review on a branch
- `snap clone`, `snap push`, `snap pull` – Remote repository support
- Branch management
- Merge functionality
- More advanced gamification features

## License

This project is open source and available under the MIT License.
