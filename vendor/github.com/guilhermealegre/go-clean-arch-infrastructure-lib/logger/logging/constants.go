package logging

// Tracer caller skip steps
const (
	DefaultCallerSkip = 2 // Number of frames to skip to get the function caller
)

// Outputs types
const (
	OutputConsole = "console"
	OutputFile    = "file"
	OutputRabbit  = "rabbit"
	OutputSQS     = "sqs"
)

// Types
const (
	BackendType  Type = "backend"
	FrontendType Type = "frontend"
)
