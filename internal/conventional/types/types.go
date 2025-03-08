package types

type CommitType = string

const (
	BUILD    CommitType = "build"
	CHORE    CommitType = "chore"
	CI       CommitType = "ci"
	DOCS     CommitType = "docs"
	FEAT     CommitType = "feat"
	FIX      CommitType = "fix"
	PERF     CommitType = "perf"
	REFACTOR CommitType = "refactor"
	REVERT   CommitType = "revert"
	STYLE    CommitType = "style"
	TEST     CommitType = "test"
)

var ValidCommitTypes = []CommitType{BUILD, CHORE, CI, DOCS, FEAT, FIX, PERF, REFACTOR, REVERT, STYLE, TEST}
