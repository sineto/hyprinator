package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

var getenv = os.Getenv

func GetSocketPath(name string) (string, error) {
	sig, err := instanceSignature()
	if err != nil {
		return "", err
	}

	rt, err := runtimeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(rt, "hypr", sig, name), nil
}

func instanceSignature() (string, error) {
	sig := getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if sig == "" {
		return "", fmt.Errorf("HYPRLAND_INSTANCE_SIGNATURE not set. Ensure Hyprland is running")
	}

	return sig, nil
}

func runtimeDir() (string, error) {
	dir := getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		return "", fmt.Errorf("XDG_RUNTIME_DIR not set")
	}

	return dir, nil
}
