package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	baseURL       = "https://ghost.yusiqo.com/pkgs/"
	requestPHPURL = "https://ghost.yusiqo.com/request.php"
	version       = "0.01"
	versionurl    = "https://raw.githubusercontent.com/yusiqo/ghost/refs/heads/main/version"
	latesturl     = "https://github.com/yusiqo/ghost/releases/latest/download/ghost"
)

type Package struct {
	Name         string   `json:"name"`
	Command      string   `json:"command"`
	Requirements []string `json:"requirements"`
}

func fetchPackage(name string) (*Package, error) {
	url := baseURL + name + ".json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Package not found: %s. Trying alternative package managers...\n", name)
		return nil, fmt.Errorf("Package not found")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var pkg Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}
	return &pkg, nil
}

func executeCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to execute command: %v", err)
	}
	return nil
}

func tryPackageManager(name string) error {
	if isCommandAvailable("yay") {
		fmt.Printf("Trying '%s' with yay...\n", name)
		cmd := exec.Command("yay", "-S", name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to install with yay: %v", err)
		}
		fmt.Printf("'%s' successfully installed with yay.\n", name)
		return nil
	}

	if isCommandAvailable("apt") {
		fmt.Printf("Trying '%s' with apt...\n", name)
		cmd := exec.Command("sudo", "apt", "install", "-y", name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to install with apt: %v", err)
		}
		fmt.Printf("'%s' successfully installed with apt.\n", name)
		return nil
	}

	return fmt.Errorf("Neither yay nor apt is available")
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func installRequirements(requirements []string) error {
	for _, req := range requirements {
		if !isCommandAvailable(req) {
			fmt.Printf("Requirement '%s' not found. Installing...\n", req)
			if err := tryPackageManager(req); err != nil {
				return fmt.Errorf("Failed to install requirement '%s': %v", req, err)
			}
		}
	}
	return nil
}

func reportPackageRequest(name string) error {
	data := map[string]string{"name": name}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON creation error: %v", err)
	}

	resp, err := http.Post(requestPHPURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP POST error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request.php returned an error: %s - %s", resp.Status, string(body))
	}

	return nil
}

func checkForUpdate() error {
	resp, err := http.Get(versionurl)
	if err != nil {
		return fmt.Errorf("Failed to check for updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to check for updates: HTTP %s", resp.Status)
	}

	latestVersion, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read version info: %v", err)
	}

	if string(latestVersion) != version {
		fmt.Printf("A new version is available: %s. Updating...\n", string(latestVersion))
		updateCommand := fmt.Sprintf("sudo curl -L %s -o /usr/local/bin/ghost && sudo chmod +x /usr/local/bin/ghost", latesturl)
		if err := executeCommand(updateCommand); err != nil {
			return fmt.Errorf("Error during update: %v", err)
		}
		fmt.Println("Update completed.")
	} else {
		fmt.Println("Ghost is already up to date.")
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ghost <command> [args]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "update":
		if err := checkForUpdate(); err != nil {
			fmt.Println("Update error:", err)
			os.Exit(1)
		}
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ghost install <package-name>")
			os.Exit(1)
		}
		pkgName := os.Args[2]
		pkg, err := fetchPackage(pkgName)
		if err != nil {
			fmt.Println("Error:", err)
			if tryErr := tryPackageManager(pkgName); tryErr != nil {
				fmt.Println("Alternative package managers failed:", tryErr)
				if reportErr := reportPackageRequest(pkgName); reportErr != nil {
					fmt.Println("Error reporting package request:", reportErr)
				}
			}
			os.Exit(1)
		}

		if err := installRequirements(pkg.Requirements); err != nil {
			fmt.Println("Error installing requirements:", err)
			os.Exit(1)
		}

		fmt.Printf("Package found: %s \nExecuting command: %s\n", pkg.Name, pkg.Command)
		if err := executeCommand(pkg.Command); err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}
		fmt.Println("Package installed successfully.")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
