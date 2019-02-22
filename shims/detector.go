package shims

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Installer interface {
	InstallOnlyVersion(depName string, installDir string) error
	InstallCNBS(orderFile string, installDir string) error
}

type DefaultDetector struct {
	V3LifecycleDir string

	AppDir string

	V3BuildpacksDir string

	OrderMetadata string
	GroupMetadata string
	PlanMetadata  string

	Installer Installer
}

func (d DefaultDetector) Detect() error {
	if err := d.Installer.InstallCNBS(d.OrderMetadata, d.V3BuildpacksDir); err != nil {
		return err
	}

	return d.RunLifecycleDetect()
}

func (d DefaultDetector) RunLifecycleDetect() error {
	if err := d.Installer.InstallOnlyVersion(V3_DETECTOR_DEP, d.V3LifecycleDir); err != nil {
		return err
	}

	cmd := exec.Command(
		filepath.Join(d.V3LifecycleDir, V3_DETECTOR_DEP),
		"-app", d.AppDir,
		"-buildpacks", d.V3BuildpacksDir,
		"-order", d.OrderMetadata,
		"-group", d.GroupMetadata,
		"-plan", d.PlanMetadata,
	)
	cmd.Env = append(os.Environ(), "PACK_STACK_ID=org.cloudfoundry.stacks."+os.Getenv("CF_STACK"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr,"OUTPUT!!!!!!!!!!!!!!!!!")
	fmt.Fprintln(os.Stderr,string(output))
	return nil
}
