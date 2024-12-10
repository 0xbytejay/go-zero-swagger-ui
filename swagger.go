package GoZeroSwaggerUI

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	swaggerFiles "github.com/swaggo/files"
)

type Config struct {
	URL string
}

var (
	config Config
)

func FromJson(jsonFilePath string) error {
	jsonFile, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return err
	}

	err = swaggerFiles.WriteFile(filepath.Base(jsonFilePath), jsonFile, os.ModePerm)
	if err != nil {
		return err
	}
	config.URL = filepath.Base(jsonFilePath)
	return nil
}

func FileSystem() http.FileSystem {

	swaggerInitializerJsFile, err := swaggerFiles.FS.OpenFile(swaggerFiles.CTX, "swagger-initializer.js", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return swaggerFiles.HTTP
	}

	initializerJS, _ := template.New("swagger-initializer").Parse(swaggerInitializerJSTpl)

	err = initializerJS.Execute(swaggerInitializerJsFile, config)
	if err != nil {
		return swaggerFiles.HTTP
	}

	err = swaggerInitializerJsFile.Close()
	if err != nil {
		return swaggerFiles.HTTP
	}

	return swaggerFiles.HTTP
}

const swaggerInitializerJSTpl = `window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">
  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
	window.ui = SwaggerUIBundle({
    	url: "{{.URL}}",
    	dom_id: '#swagger-ui',
    	deepLinking: true,
    	presets: [
			SwaggerUIBundle.presets.apis,
			SwaggerUIStandalonePreset
    	],
    	plugins: [
			SwaggerUIBundle.plugins.DownloadUrl
    	],
    	layout: "StandaloneLayout"
	});
  //</editor-fold>
};
`
