package syslog_server

import (
	"log"
	"os"
	"time"

	"sirendaou.com/duserver/common/syslog"
)

type logFile struct {
	dir      string
	filename string
	date     time.Time
	file     *os.File
}

var (
	DATEFORMAT = "2006-01-02"
)

func NewLogFile(fileDir, fileName string) *logFile {
	t, err := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	if err != nil {
		log.Println(err)
	}
	if len(fileDir) > 0 && fileDir[len(fileDir)-1] != '/' {
		fileDir += "/"
	}
	file, err := os.OpenFile(fileDir+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 444)
	if err != nil {
		log.Println(2, err)
	}
	return &logFile{dir: fileDir, filename: fileName, date: t, file: file}
}

func (f *logFile) WriteLogMsg(logMsg *syslog.LogMsg) {
	f.check()
	f.file.WriteString(logMsg.Format())
}

func (f *logFile) WriteString(logMsg string) {
	f.check()
	f.file.WriteString(logMsg)
}

func (f *logFile) isMustRename() bool {
	t, err := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	if err != nil {
		log.Println(err)
	}
	if t.After(f.date) {
		return true
	}
	return false
}

func (f *logFile) check() {
	fn := f.dir + "/" + f.filename + "." + f.date.Format(DATEFORMAT)
	if f.isMustRename() && !isExist(fn) {
		if f.file != nil {
			f.file.Close()
		}
		err := os.Rename(f.dir+f.filename, fn)
		if err != nil {
			// WARN
		}
		f.date, _ = time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		f.file, _ = os.OpenFile(f.dir+f.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 444)
	}
}

func openFile(path string) *os.File {
	if !isExist(path) {
		file, err := os.Create(path)
		if err != nil {
			log.Println("1", err)
		}
		return file
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		log.Println(2, err)
	}
	return file
}
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
