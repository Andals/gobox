package log

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

var logLevels map[int]string

func init() {
	logLevels = make(map[int]string)

	logLevels[LEVEL_DEBUG] = "debug"
	logLevels[LEVEL_INFO] = "info"
	logLevels[LEVEL_NOTICE] = "notice"
	logLevels[LEVEL_WARNING] = "warning"
	logLevels[LEVEL_ERROR] = "error"
	logLevels[LEVEL_CRITICAL] = "critical"
	logLevels[LEVEL_ALERT] = "alert"
	logLevels[LEVEL_EMERGENCY] = "emergency"
}
