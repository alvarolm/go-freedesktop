package freedesktop

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var xdgIconExtensions []string = []string{".png", ".svg", ".xpm"}

// Returns the current icon theme, if it can be found
func GetIconTheme() string {
	session := os.Getenv("DESKTOP_SESSION")
	switch session {
	case "mate":
		cmd := exec.Command("gsettings", "get", "org.mate.interface", "icon-theme")
		out, err := cmd.Output()
		if err == nil {
			return strings.Trim(string(out), "' \n")
		}
	case "gnome":
		cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme")
		out, err := cmd.Output()
		if err == nil {
			return strings.Trim(string(out), "' \n")
		}
	case "kde":
		// TODO: look up KDE icon theme
		// ???: where is this setting stored?
		// try looking more in ~/.kde4/share/config/kdeglobals
	}
	return ""
}

func AppIcon(icon string) string {
	return AppIconForSize(icon, 48)
}

func AppIconForSize(icon string, size int) (filename string) {
	theme := GetIconTheme()
	fmt.Println("usign theme:", theme)
	if theme != "" {
		filename = FindIconHelper(icon, size, theme)
		if filename != "" {
			return
		}
	}

	filename = FindIconHelper(icon, size, "hicolor")
	if filename != "" {
		return
	}

	return LookupFallbackIcon(icon)
}

func FindIconHelper(icon string, size int, theme string) (filename string) {
	filename = LookupIcon(icon, size, theme)
	if filename != "" {
		return
	}

	// TODO: look in this theme's parents
	// ???: how to get the icon theme's parents?
	return ""
}

func LookupIcon(icon string, size int, theme string) string {
	themeDirs := make([]string, 0)
	for _, dir := range xdgIcons {
		themeDir := path.Join(dir, theme)
		if _, err := os.Stat(themeDir); err != nil {
			// theme not found
			continue
		}
		themeDirs = append(themeDirs, themeDir)
	}

	for _, themeDir := range themeDirs {
		themeConf, err := ParseConfigFile(path.Join(themeDir, "index.theme"))
		if err != nil {
			// no index.theme
			continue
		}

		var lookupSubDirs []string
		for header, conf := range themeConf {
			if header != "Icon Theme" {
				themeSize := conf["Size"]
				if themeSize == strconv.Itoa(size) {
					lookupSubDirs = append(lookupSubDirs, header)
				}
			}
		}

		if len(lookupSubDirs) == 0 {
			// No sizes matchs
			continue
		}

		for _, subDir := range lookupSubDirs {
			for _, ext := range xdgIconExtensions {
				file := path.Join(themeDir, subDir, icon+ext)
				if _, err := os.Stat(file); err == nil {
					return file
				}
			}
		}
	}

	return ""
}

func LookupFallbackIcon(icon string) string {
	for _, dir := range xdgIcons {
		for _, ext := range xdgIconExtensions {
			file := path.Join(dir, icon+ext)
			if _, err := os.Stat(file); err == nil {
				return file
			}
		}
	}

	return ""
}
