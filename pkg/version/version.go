package version

import (
	"fmt"
	"runtime"
)

// Informações de versão do aplicativo
var (
	// Version é a versão do aplicativo
	Version = "1.3.2"

	// BuildDate é a data de compilação (injetada durante o build via -ldflags)
	BuildDate = "development"

	// GitCommit é o hash do commit Git (injetada durante o build via -ldflags)
	GitCommit = "unknown"
)

// GetVersion retorna as informações de versão formatadas
func GetVersion() string {
	return fmt.Sprintf("iCRMSenderEmail v%s", Version)
}

// GetFullVersion retorna informações completas de versão
func GetFullVersion() string {
	return fmt.Sprintf(
		"iCRMSenderEmail v%s\n"+
			"Build Date: %s\n"+
			"Git Commit: %s\n"+
			"Go Version: %s\n"+
			"OS/Arch: %s/%s",
		Version,
		BuildDate,
		GitCommit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// PrintVersion imprime informações de versão no stdout
func PrintVersion() {
	fmt.Println(GetFullVersion())
}
