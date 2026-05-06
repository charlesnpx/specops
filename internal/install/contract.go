package install

type Operation string

const (
	OperationPlan      Operation = "plan"
	OperationInstall   Operation = "install"
	OperationUninstall Operation = "uninstall"
)

type Target string

const (
	TargetClaude Target = "claude"
	TargetCodex  Target = "codex"
	TargetTools  Target = "tools"
	TargetAll    Target = "all"
)

type Options struct {
	Operation   Operation
	Target      Target
	InstallRoot string
	Version     string
}

type Report struct {
	Schema    int                     `json:"schema"`
	Name      string                  `json:"name"`
	Version   string                  `json:"version"`
	Operation Operation               `json:"operation"`
	Kind      string                  `json:"kind"`
	Targets   map[string]TargetReport `json:"targets"`
	Warnings  []string                `json:"warnings"`
}

type TargetReport struct {
	Files []FileReport `json:"files"`
}

type FileReport struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256,omitempty"`
}

type filePlan struct {
	Target     Target
	Path       string
	Content    []byte
	Executable bool
	ToolBinary bool
}
