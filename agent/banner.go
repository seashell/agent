package agent

import (
	"fmt"

	"github.com/seashell/agent/version"
)

// Banner is a banner to be displayed when
// the Seashell agent is started
var Banner = fmt.Sprintf(`
-------------------------------------------------
                            _              _   _ 
 ___    ___    __ _   ___  | |__     ___  | | | |
/ __|  / _ \  / _' | / __| | '_ \   / _ \ | | | |
\__ \ |  __/ | (_| | \__ \ | | | | |  __/ | | | |
|___/  \___|  \__'_| |___/ |_| |_|  \___| |_| |_|
																	
					    {{ .AnsiColor.Cyan }}%s{{ .AnsiColor.Default }}
-------------------------------------------------

`, version.GetVersion().VersionNumber())
