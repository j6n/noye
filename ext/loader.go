package ext

import (
	"os"
	"path/filepath"
)

// LoadScripts tries to load scripts in `dir`
func (m *Manager) LoadScripts(dir string) {
	scripts := getFiles(dir)
	for script := range scripts {
		log.Infof("found script: '%s'\n", script)
		if err := m.Load(script); err != nil {
			log.Errorf("loading script '%s': '%s'\n", script, err)
		}
	}
}

func getFiles(base string) <-chan string {
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
			log.Errorf("Walking '%s': %s\n", base, err)
		}
		close(scripts)
	}()

	return scripts
}
