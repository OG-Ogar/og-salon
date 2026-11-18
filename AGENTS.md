# Project notes for coding agents

- This is a Go 1.23 server-rendered application with TypeScript enhancements in `frontend/`.
- Run `go test ./...` for Go validation.
- Run `npm run build` to compile TypeScript into `static/app.js`.
- Coordinate task claims and handoffs through the primary checkout's `multi-agent-orchestrator/` state. Implement product changes only in an assigned worktree and claimed paths.
