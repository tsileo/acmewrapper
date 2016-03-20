package acmewrapper

import (
	"errors"
	"io/ioutil"
	"os"
)

var ErrNotHandled = errors.New("not handled")

func (w *AcmeWrapper) loadFile(path string) ([]byte, error) {
	//use custom load file callback?
	if w.Config.LoadFileCallback != nil {
		if b, err := w.Config.LoadFileCallback(path); err == nil {
			return b, nil
		} else if err != ErrNotHandled {
			return nil, err
		}
	}
	//default load from disk
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (w *AcmeWrapper) saveFile(path string, contents []byte) error {
	//use custom save file callback?
	if w.Config.SaveFileCallback != nil {
		if err := w.Config.SaveFileCallback(path, contents); err == nil {
			return nil
		} else if err != ErrNotHandled {
			return err
		}
	}
	//default save to disk (current user read+write only!)
	if err := ioutil.WriteFile(path, contents, 0600); err != nil {
		return err
	}
	return nil
}

func (w *AcmeWrapper) backupAndSaveFile(path string, contents []byte) error {
	//load previous file and hold in memory
	prev, err := w.loadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	//then save over previous file
	if err := w.saveFile(path, contents); err != nil {
		//if fails, attempt to restore previous file just incase it was damaged
		w.saveFile(path, prev)
		return err
	}
	//if save was successful, overwrite previous backup, with previous file
	if len(prev) > 0 {
		if err := w.saveFile(path+".bak", prev); err != nil {
			return err
		}
	}
	return nil
}
