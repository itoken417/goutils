package chrome

import (
	"github.com/sclevine/agouti"
	"log"
)

var driver = agouti.ChromeDriver(
	agouti.ChromeOptions("prefs", map[string]interface{}{
		"download.default_directory":         "./tmp",
		"download.prompt_for_download":       false,
		"download.directory_upgrade":         true,
		"plugins.plugins_disabled":           "Chrome PDF Viewer",
		"plugins.always_open_pdf_externally": true,
	}),
	agouti.ChromeOptions("args", []string{
		"--headless",
		"--allow-insecure-localhost",
		"--disable-gpu",
		"--homepage=about:blank",
		"--no-first-run",
		"--no-default-browser-check",
		"--no-sandbox",
		"--window-size=1280,800",
		"--disable-popup-blocking",
		"--disable-dev-shm-usage",
	}),
	agouti.ChromeOptions(
		"binary", "/usr/bin/google-chrome-stable",
	),
)

func Init() {
	if err := driver.Start(); err != nil {
		log.Println(err)
	}
}

func GetDriver() *agouti.WebDriver {
	return driver
}

func GetNewPage(ssl_check bool) *agouti.Page {
	if ssl_check {
		page, err := driver.NewPage()
		if err != nil {
			log.Println(err)
		}
		return page
	} else {
		capabilities := agouti.NewCapabilities().With("acceptInsecureCerts")
		page, err := driver.NewPage(agouti.Desired(capabilities))
		if err != nil {
			log.Println(err)
		}
		return page
	}
}

func Destroy() {
	defer driver.Stop()
}
