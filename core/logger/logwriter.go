package logger

type LogWriter interface {
	LogWrite(info *LogInfo)
	Close()
}
