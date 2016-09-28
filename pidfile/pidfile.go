package pidfile

import (
	"errors"
	"os"
	"strconv"
)

const (
	PID_FILE_MIN_ID          = 0
	PID_FILE_TMP_PATH_SUFFIX = ".tmp"
)

type PidFile struct {
	*Pid
	*File

	IsTmp      bool
	OriginPath string
}

func CreatePidFile(path string) (*PidFile, error) {
	file := NewFile(path)
	isTmp := false

	pid, err := ReadPidFromFile(file)
	if err == nil && pid.ProcessExist() {
		file = NewFile(path + PID_FILE_TMP_PATH_SUFFIX)
		isTmp = true
	}

	pid = NewPid(os.Getpid())
	if err := WritePidToFile(file, pid); err != nil {
		return nil, err
	}

	return &PidFile{
		Pid:        pid,
		File:       file,
		IsTmp:      isTmp,
		OriginPath: path,
	}, nil
}

func ClearPidFile(pidfile *PidFile) error {
	pid, err := ReadPidFromFile(NewFile(pidfile.Path))
	if err == nil && pidfile.Id == pid.Id {
		if err = pidfile.Remove(); err != nil {
			return err
		}
	}

	if !pidfile.IsTmp {
		tfile := NewFile(pidfile.Path + PID_FILE_TMP_PATH_SUFFIX)
		tpid, err := ReadPidFromFile(tfile)
		if err == nil && tpid.ProcessExist() {
			if err = tfile.Rename(pidfile.Path); err != nil {
				return err
			}
		}
	} else {
		ofile := NewFile(pidfile.OriginPath)
		opid, err := ReadPidFromFile(ofile)
		if err == nil && opid.Id == pidfile.Id {
			if err = ofile.Remove(); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadPidFromFile(file *File) (*Pid, error) {
	fb, err := file.Read()
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(string(fb))
	if err != nil || id <= PID_FILE_MIN_ID {
		return nil, errors.New("pid file data error")
	}

	return NewPid(id), nil
}

func WritePidToFile(file *File, pid *Pid) error {
	fb := []byte(strconv.Itoa(pid.Id))
	return file.Write(fb)
}
