package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type openAPI struct {
	OpenAPI    string                   `json:"openapi"`
	Info       map[string]interface{}   `json:"info"`
	Servers    []map[string]interface{} `json:"servers"`
	Paths      map[string]interface{}   `json:"paths"`
	Components map[string]interface{}   `json:"components"`
}

func mergeSwaggerDocs(docsDir string) (*openAPI, error) {
	mergedDocs := &openAPI{
		OpenAPI:    "3.0.0",
		Info:       make(map[string]interface{}),
		Paths:      make(map[string]interface{}),
		Components: make(map[string]interface{}),
		Servers: []map[string]interface{}{
			{
				"url":         "http://127.0.0.1:5005",
				"description": "localhost",
			},
			{
				"url":         "ws://127.0.0.1:5005",
				"description": "webSocket",
			},
			{
				"url":         "https://example.com",
				"description": "Main API Server",
			},
		},
	}

	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			//serviceName := filepath.Base(path)
			//serviceName = serviceName[:len(serviceName)-len(filepath.Ext(serviceName))]

			var doc openAPI
			err = yaml.Unmarshal(fileContent, &doc)
			if err != nil {
				return fmt.Errorf("error unmarshalling file %s: %w", path, err)
			}

			for pathKey, pathValue := range doc.Paths {
				if _, exists := mergedDocs.Paths[pathKey]; !exists {
					mergedDocs.Paths[pathKey] = pathValue
				} else {
					mergedPaths := mergedDocs.Paths[pathKey].(map[string]interface{})
					newPaths := pathValue.(map[string]interface{})
					for opKey, opValue := range newPaths {
						mergedPaths[opKey] = opValue
					}
				}
			}

			for compKey, compValue := range doc.Components {
				if _, exists := mergedDocs.Components[compKey]; !exists {
					mergedDocs.Components[compKey] = compValue
				} else {
					mergedComponents := mergedDocs.Components[compKey].(map[string]interface{})
					newComponents := compValue.(map[string]interface{})
					for subCompKey, subCompValue := range newComponents {
						mergedComponents[subCompKey] = subCompValue
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error merging Swagger docs: %w", err)
	}

	mergedDocs.Info["title"] = "API Documentation"
	mergedDocs.Info["version"] = "1.0.0"

	return mergedDocs, nil
}

func swaggerRoute(rs *RouterService) {

	docsDir := "./docs/swagger_docs"
	app := rs.app
	app.Get("/swagger", func(c *fiber.Ctx) error {
		mergedDocs, err := mergeSwaggerDocs(docsDir)
		if err != nil {
			log.Printf("Error merging Swagger docs: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error generating Swagger documentation")
		}

		mergedYAML, err := yaml.Marshal(mergedDocs)
		if err != nil {
			log.Printf("Error marshalling merged Swagger docs to YAML: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error converting Swagger documentation to YAML")
		}

		c.Set("Content-Type", "application/yaml")
		return c.Send(mergedYAML)
	})

	app.Get("/swagger-ui/*", swagger.New(swagger.Config{

		URL: "/swagger",

		CustomScript: template.JS(`
              setTimeout(function() {
					const token = localStorage.getItem("jwt_token");
					if (token) {
						ui.authActions.authorize({
							BearerAuth: {
								name: "BearerAuth",
								schema: {
									type: "http",
									in: "header",
									scheme: "bearer",
									bearerFormat: "JWT"
								},
								value: token
							}
						});
					}

					ui.getSystem().authActions.authorize = (auth) => {
						if (auth.BearerAuth && auth.BearerAuth.value) {
							const token = localStorage.getItem("jwt_token");

							const tokenValue = auth.BearerAuth.value.replace("Bearer ", "");
							localStorage.setItem("jwt_token", tokenValue);
 							

						}
 						window.location.reload()
					 
					};
				},2000)
        `),
	}))

}
