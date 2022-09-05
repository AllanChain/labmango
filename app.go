package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

type LabManConfig struct {
	LabDir          string `json:"labDir" yaml:"labDir"`
	JupyterTemplate string `json:"jupyterTemplate" yaml:"jupyterTemplate"`
	LyXTemplate     string `json:"lyxTemplate" yaml:"lyxTemplate"`
}

// App struct
type App struct {
	ctx    context.Context
	config LabManConfig
	cmd    *exec.Cmd
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
		LabDir:          path.Join(homedir, "labs"),
		JupyterTemplate: path.Join(homedir, "Templates", "lab-data.ipynb"),
		LyXTemplate:     path.Join(homedir, "Templates", "lab.lyx"),
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
	a.signalConfigUpdate()
}

func (a *App) signalConfigUpdate() {
	runtime.EventsEmit(a.ctx, "config-update", a.config)
	a.RefreshLabs()
}

func (a *App) PromptError(msg string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Message: msg,
		Type:    runtime.ErrorDialog,
	})
}

func (a *App) RefreshLabs() {
	labs := a.ListLabs()
	runtime.EventsEmit(a.ctx, "labs-refresh", labs)
	log.Println("Labs refreshed.")
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

func (a *App) ChangeJupyterTemplate() {
	start_dir := path.Dir(a.config.JupyterTemplate)
	if _, err := os.Stat(start_dir); err != nil {
		start_dir = getHomeDir()
	}
	jupyter_template, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: start_dir,
	})
	if jupyter_template != "" {
		a.config.JupyterTemplate = jupyter_template
		a.SyncConfig()
	}
}
func (a *App) ChangeLyXTemplate() {
	start_dir := path.Dir(a.config.LyXTemplate)
	if _, err := os.Stat(start_dir); err != nil {
		start_dir = getHomeDir()
	}
	lyx_template, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: start_dir,
	})
	if lyx_template != "" {
		a.config.LyXTemplate = lyx_template
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
	sort.Slice(files, func(i, j int) bool {
		info_i, err := files[i].Info()
		if err != nil {
			log.Panic(err)
		}
		info_j, err := files[j].Info()
		if err != nil {
			log.Panic(err)
		}
		return info_i.ModTime().Unix() > info_j.ModTime().Unix()
	})
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		filename := file.Name()
		if strings.HasPrefix(filename, ".") {
			continue
		}
		labs = append(labs, file.Name())
	}
	return labs
}

func (a *App) ExploreLab(lab string) error {
	return open.Start(path.Join(a.config.LabDir, lab))
}

func (a *App) CreateLab(lab string) {
	lab_dir := path.Join(a.config.LabDir, lab)
	err := os.MkdirAll(lab_dir, 0755)
	if err != nil {
		a.PromptError(fmt.Sprintf("Unable to create %s: %s", lab_dir, err))
		return
	}
	for _, subdir := range []string{"report", "data"} {
		dir := path.Join(lab_dir, subdir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			a.PromptError(fmt.Sprintf("Unable to create %s: %s", dir, err))
			return
		}
	}
	for _, subdir := range []string{"images", "figs", "tables"} {
		dir := path.Join(lab_dir, "report", subdir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			a.PromptError(fmt.Sprintf("Unable to create %s: %s", dir, err))
			return
		}
	}
	if _, err := os.Stat(a.config.JupyterTemplate); err != nil {
		a.PromptError(fmt.Sprintf("Jupyter template %s not found", a.config.JupyterTemplate))
		return
	}
	template_data, err := os.ReadFile(a.config.JupyterTemplate)
	if err != nil {
		a.PromptError(fmt.Sprintf("Unable to read Jupyter template %s", a.config.JupyterTemplate))
		return
	}
	jupyter_filename := path.Join(lab_dir, "data", "data.ipynb")
	err = os.WriteFile(jupyter_filename, template_data, 0644)
	if err != nil {
		a.PromptError(fmt.Sprintf("Unable to write Jupyter notebook %s", jupyter_filename))
		return
	}

	if _, err := os.Stat(a.config.LyXTemplate); err != nil {
		a.PromptError(fmt.Sprintf("LyX template %s not found", a.config.LyXTemplate))
		return
	}
	template_data, err = os.ReadFile(a.config.LyXTemplate)
	if err != nil {
		a.PromptError(fmt.Sprintf("Unable to read LyX template %s", a.config.LyXTemplate))
		return
	}
	template_data = bytes.Replace(template_data, []byte("title"), []byte(lab), 1)
	lyx_filename := path.Join(lab_dir, "report", lab+".lyx")
	err = os.WriteFile(lyx_filename, template_data, 0644)
	if err != nil {
		a.PromptError(fmt.Sprintf("Unable to write LyX file %s", lyx_filename))
		return
	}

	a.RefreshLabs()
}

func (a *App) DeleteLab(lab string) error {
	response, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Message:       "Are you sure you want to delete?",
		Buttons:       []string{"OK", "Cancel"},
		DefaultButton: "Cancel",
		CancelButton:  "Cancel",
		Type:          runtime.QuestionDialog,
		Title:         "Confirm Delete",
	})
	if response != "Yes" {
		return nil
	}
	lab_dir := path.Join(a.config.LabDir, lab)
	if _, err := os.Stat(lab_dir); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: fmt.Sprintf("%s not found!", lab),
		})
		return err
	}
	recycle_dir := path.Join(a.config.LabDir, ".recycle")
	if _, err := os.Stat(recycle_dir); err != nil {
		err := os.MkdirAll(recycle_dir, 0755)
		if err != nil {
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Message: "Error making",
			})
		}
	}
	now := time.Now()
	new_lab_name := fmt.Sprintf(
		"%s-%04d%02d%02d_%02d%02d%02d",
		lab,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
	)
	err := os.Rename(lab_dir, path.Join(recycle_dir, new_lab_name))
	if err != nil {
		log.Println(err)
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Error moving",
		})
	}
	a.RefreshLabs()
	return nil
}

func (a *App) LaunchLab(lab string) {
	cmd := exec.Command("jupyter-lab", "--no-browser")
	cmd.Dir = a.config.LabDir
	a.cmd = cmd
	a.signalLabStatus()
	stdoutPipe, _ := cmd.StderrPipe()
	stdoutReader := bufio.NewReader(stdoutPipe)
	cmd.Start()

	go func() {
		bufReader := bufio.NewReader(stdoutReader)
		nextLineIsURL := false
		processing := true
		for {
			output, _, err := bufReader.ReadLine()
			if err != nil {
				break
			}
			if !processing {
				continue
			}
			if nextLineIsURL {
				labURL := string(output[bytes.LastIndex(output, []byte(" "))+1:])
				if lab != "" {
					labURL = strings.Replace(labURL, "?token", "/tree/"+lab+"/data?token", 1)
				}
				open.Start(labURL)
				processing = false
			} else if bytes.HasSuffix(output, []byte("is running at:")) {
				nextLineIsURL = true
			}
			log.Printf("Line: %s", output)
		}
		a.cmd = nil
		log.Println("Command end")
		a.signalLabStatus()
	}()
}

func (a *App) KillLab() {
	if a.cmd == nil {
		a.PromptError("Jupyter Lab is not running.")
	}
	a.cmd.Process.Kill()
}

func (a *App) signalLabStatus() {
	runtime.EventsEmit(a.ctx, "jlab-running", a.cmd != nil)
}

func (a *App) EditReport(lab string) {
	open.Start(path.Join(a.config.LabDir, lab, "report", lab+".lyx"))
}
