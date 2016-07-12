package level

const (
	LEVEL_DEBUG     = 1
	LEVEL_INFO      = 2
	LEVEL_NOTICE    = 3
	LEVEL_WARNING   = 4
	LEVEL_ERROR     = 5
	LEVEL_CRITICAL  = 6
	LEVEL_ALERT     = 7
	LEVEL_EMERGENCY = 8
)

var LogLevels map[int]string

func init() {
	LogLevels = make(map[int]string)

	LogLevels[LEVEL_DEBUG] = "debug"
	LogLevels[LEVEL_INFO] = "info"
	LogLevels[LEVEL_NOTICE] = "notice"
	LogLevels[LEVEL_WARNING] = "warning"
	LogLevels[LEVEL_ERROR] = "error"
	LogLevels[LEVEL_CRITICAL] = "critical"
	LogLevels[LEVEL_ALERT] = "alert"
	LogLevels[LEVEL_EMERGENCY] = "emergency"
}
