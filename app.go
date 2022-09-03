package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

type LabManConfig struct {
	LabDir string `json:"labDir"`
}

// App struct
type App struct {
	ctx    context.Context
	config LabManConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadConfig()
}

func (a *App) defaultConfig() LabManConfig {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return LabManConfig{
		LabDir: path.Join(homedir, "labs"),
	}
}

func (a *App) loadConfig() {
	config_path := a.getConfigPath()
	log.Printf("Reading config file: %s", config_path)
	config_data, err := os.ReadFile(config_path)
	a.config = a.defaultConfig()
	if err == nil {
		yaml.Unmarshal(config_data, &a.config)
	} else {
		a.SyncConfig()
	}
}

func (a *App) CurrentConfig() *LabManConfig {
	return &a.config
}

func (a *App) SyncConfig() {
	config_path := a.getConfigPath()
	config_data, err := yaml.Marshal(&a.config)
	if err != nil {
		return
	}
	os.WriteFile(config_path, config_data, 0666)
	log.Println("Config file written.")
	runtime.EventsEmit(a.ctx, "config-update", a.config)
}

func (a *App) ChangeLabDir() {
	log.Println("ChangeLabDir")
	var start_dir string
	if _, err := os.Stat(a.config.LabDir); err != nil {
		start_dir = getHomeDir()
	} else {
		start_dir = a.config.LabDir
	}
	labdir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: start_dir,
	})
	if err != nil {
		log.Println(err)
	}
	log.Printf("chosen %s", labdir)
	if labdir != "" {
		a.config.LabDir = labdir
		a.SyncConfig()
	}
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homedir
}

func (a *App) getConfigPath() string {
	homedir := getHomeDir()
	return path.Join(homedir, ".config", "labman.yaml")
}

func (a *App) ListLabs() []string {
	labs := []string{}
	if _, err := os.Stat(a.config.LabDir); err != nil {
		return labs
	}
	files, err := os.ReadDir(a.config.LabDir)
	if err != nil {
		return labs
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		labs = append(labs, file.Name())
	}
	return labs
}

func (a *App) ExploreLab(lab string) error {
	return open.Start(path.Join(a.config.LabDir, lab))
}
