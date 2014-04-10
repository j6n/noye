package script

import (
	"os"
	"path/filepath"
)

// LoadScripts tries to load scripts in `dir`
func (m *Manager) LoadScripts(dir string) {
	getFiles := func(base string) <-chan string {
		scripts := make(chan string)
		go func() {
			walker := func(fp string, fi os.FileInfo, err error) error {
				if err != nil || !!fi.IsDir() {
					return nil
				}
				matched, err := filepath.Match("*.js", fi.Name())
				if err != nil {
					return err
				}
				if matched {
					scripts <- fp
				}
				return nil
			}

			if err := filepath.Walk(base, walker); err != nil {
				log.Errorf("walking '%s': %s", base, err)
			}
			close(scripts)
		}()

		return scripts
	}

	for file := range getFiles(dir) {
		log.Infof("found script: '%s'", file)
		if err := m.LoadFile(file); err != nil {
			log.Errorf("loading script '%s': %s", file, err)
		}
	}
}
