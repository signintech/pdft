package core

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

//ErrPathIsNotDir path is not dir
var ErrPathIsNotDir = errors.New("path is not dir")

//ErrPathIsNotEmpty path is not empty
var ErrPathIsNotEmpty = errors.New("path is not empty")

//CreateSubCmd create sub cmd
type CreateSubCmd struct {
	ProjectPath string
}

//Create create
func (c CreateSubCmd) Create() error {

	path := c.ProjectPath
	err := c.checkDir(path)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = c.createTmplJSON(path)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func (c CreateSubCmd) checkDir(path string) error {

	isExist, err := c.isDirExists(path)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if !isExist {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return errors.Wrap(err, "")
		}
	} else {
		isdir, err := c.isDir(path)
		if err != nil {
			return errors.Wrap(err, "")
		}
		if !isdir {
			return ErrPathIsNotDir
		}
	}

	isEmpty, err := c.isDirEmpty(path)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if !isEmpty {
		return ErrPathIsNotEmpty
	}

	return nil
}

func (c CreateSubCmd) isDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (c CreateSubCmd) isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	return info.IsDir(), nil
}

func (c CreateSubCmd) isDirEmpty(path string) (bool, error) {

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func (c CreateSubCmd) createTmplJSON(path string) error {
	return nil
}
