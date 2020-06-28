package logger

type ConsoleLogWriter struct {
}

func (c *ConsoleLogWriter) LogWrite(info *LogInfo) {
	if info == nil {
		return
	}
	info.Println()
}

func (c *ConsoleLogWriter) Close() {

}

func NewConsoleLogWriter() *ConsoleLogWriter {
	c := new(ConsoleLogWriter)
	return c
}
