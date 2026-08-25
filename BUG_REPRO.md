# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	homemaker-followup/cmd/followup	[no test files]
ok  	homemaker-followup/internal/dashboard	0.002s
ok  	homemaker-followup/internal/domain	0.002s
ok  	homemaker-followup/internal/excel	0.009s
--- FAIL: TestFollowupHeaderError (0.00s)
    loader_test.go:20: expected startup format error
FAIL
FAIL	homemaker-followup/internal/followup	0.015s
ok  	homemaker-followup/internal/httpapi	0.009s
ok  	homemaker-followup/internal/reminders	0.002s
ok  	homemaker-followup/internal/reporting	0.003s
ok  	homemaker-followup/internal/storage	0.031s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/followup): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/followup): exit `0`
- Frontend build (web): exit `0`
