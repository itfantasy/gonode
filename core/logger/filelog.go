package logger

import (
	"fmt"
	"os"

	"github.com/itfantasy/gonode/utils/io"
)

type FileLogWriter struct {
	filename string
	file     *os.File
	logchan  chan *LogInfo
}

func (f *FileLogWriter) LogWrite(info *LogInfo) {
	if info == nil {
		return
	}
	info.Println()
	f.logchan <- info
}

func (f *FileLogWriter) Close() {
	f.logchan <- nil
}

func (f *FileLogWriter) dispose() {
	f.file.Sync()
	f.file.Close()
	close(f.logchan)
}

func NewFileLogWriter(filename string) (*FileLogWriter, error) {
	f := new(FileLogWriter)
	if !io.FileExists(filename) {
		dir := io.FetchDirByFilePath(filename)
		io.MakeDir(dir)
	}
	f.filename = filename
	fd, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	f.file = fd
	f.logchan = make(chan *LogInfo, 1024)
	go func() {
		defer f.dispose()
		for info := range f.logchan {
			if info == nil {
				break
			}
			_, err := fmt.Fprint(f.file, info.FormatString())
			if err != nil {
				fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", f.filename, err)
				return
			}
		}
	}()
	return f, nil
}
