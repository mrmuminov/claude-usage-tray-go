package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// Platform-specific functions (declared in install_{platform}.go):
//   installDir() string
//   binaryName() string
//   setupAutostart(binPath string) error
//   removeAutostart() error
//   autostartConfigured() bool

// isInstalled checks if the binary is already installed.
func isInstalled() bool {
	destPath := filepath.Join(installDir(), binaryName())
	_, err := os.Stat(destPath)
	return err == nil
}

// Install copies the current binary to the install location and sets up autostart.
func Install() error {
	destDir := installDir()
	destPath := filepath.Join(destDir, binaryName())

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", destDir, err)
	}

	srcPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("determine current executable: %w", err)
	}
	srcPath, err = filepath.EvalSymlinks(srcPath)
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}

	srcAbs, _ := filepath.Abs(srcPath)
	destAbs, _ := filepath.Abs(destPath)
	if srcAbs != destAbs {
		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("copy binary to %s: %w", destPath, err)
		}
		fmt.Printf("Binary installed to: %s\n", destPath)
	} else {
		fmt.Printf("Binary already at: %s\n", destPath)
	}

	if err := setupAutostart(destPath); err != nil {
		return fmt.Errorf("setup autostart: %w", err)
	}
	fmt.Println("Autostart configured.")

	if runtime.GOOS == "linux" {
		fmt.Printf("\nHint: Ensure %s is in your PATH.\n", destDir)
		fmt.Printf("  export PATH=\"%s:$PATH\"\n", destDir)
	}

	return nil
}

// Uninstall removes autostart config and the installed binary.
func Uninstall() error {
	destDir := installDir()
	destPath := filepath.Join(destDir, binaryName())

	if err := removeAutostart(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: remove autostart: %v\n", err)
	} else {
		fmt.Println("Autostart removed.")
	}

	srcPath, _ := os.Executable()
	srcPath, _ = filepath.EvalSymlinks(srcPath)
	srcAbs, _ := filepath.Abs(srcPath)
	destAbs, _ := filepath.Abs(destPath)

	if srcAbs == destAbs {
		fmt.Println("Note: Cannot remove running binary. Delete manually:")
		fmt.Printf("  rm %s\n", destPath)
	} else {
		if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove binary %s: %w", destPath, err)
		}
		fmt.Printf("Binary removed: %s\n", destPath)
	}

	return nil
}

// PrintStatus shows current installation state.
func PrintStatus() {
	destDir := installDir()
	destPath := filepath.Join(destDir, binaryName())

	fmt.Printf("claude-usage-tray-go %s\n", Version)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Install path: %s\n", destPath)

	if _, err := os.Stat(destPath); err == nil {
		fmt.Println("Binary installed: yes")
	} else {
		fmt.Println("Binary installed: no")
	}

	if autostartConfigured() {
		fmt.Println("Autostart: configured")
	} else {
		fmt.Println("Autostart: not configured")
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmpPath := dst + ".tmp"
	out, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, dst)
}
