package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Module struct {
	Name           string
	NameLower      string
	NameCamel      string
	NameSnake      string
	Timestamp      string
	PackagePath    string
}

func main() {
	moduleName := flag.String("name", "", "Name of the module to generate")
	migrationsOnly := flag.Bool("migrations", false, "Generate only migrations")
	flag.Parse()

	if *moduleName == "" {
		fmt.Println("Please provide a module name using -name flag")
		fmt.Println("Usage:")
		fmt.Println("  Complete module: go run cmd/generator/main.go -name ModuleName")
		fmt.Println("  Only migrations: go run cmd/generator/main.go -name ModuleName -migrations")
		os.Exit(1)
	}

	module := Module{
		Name:        *moduleName,
		NameLower:   strings.ToLower(*moduleName),
		NameCamel:   toCamelCase(*moduleName),
		NameSnake:   toSnakeCase(*moduleName),
		Timestamp:   time.Now().Format("20060102150405"),
		PackagePath: "personal-api",
	}

	// Define all possible templates
	allTemplates := map[string]string{
		fmt.Sprintf("migrations/%s_create_%s_table.up.sql", module.Timestamp, module.NameSnake):   migrationUpTemplate,
		fmt.Sprintf("migrations/%s_create_%s_table.down.sql", module.Timestamp, module.NameSnake): migrationDownTemplate,
		fmt.Sprintf("internal/entity/%s.go", module.NameLower):                                    entityTemplate,
		fmt.Sprintf("internal/repository/%s_repository.go", module.NameLower):                     repositoryTemplate,
		fmt.Sprintf("internal/service/%s_service.go", module.NameLower):                          serviceTemplate,
		fmt.Sprintf("api/v1/handler/%s_handler.go", module.NameLower):                           handlerTemplate,
		fmt.Sprintf("docs/setup_%s.md", module.NameLower):                                       setupInstructionsTemplate,
	}

	// Select templates based on flags
	templates := make(map[string]string)
	if *migrationsOnly {
		// Only include migration templates
		for filename, tmpl := range allTemplates {
			if strings.HasPrefix(filename, "migrations/") {
				templates[filename] = tmpl
			}
		}
		fmt.Printf("Generating migrations for %s module...\n", module.Name)
	} else {
		// Include all templates
		templates = allTemplates
		fmt.Printf("Generating complete %s module...\n", module.Name)
	}

	// Generate files
	for filename, tmpl := range templates {
		if err := generateFile(filename, tmpl, module); err != nil {
			fmt.Printf("Error generating %s: %v\n", filename, err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", filename)
	}

	if *migrationsOnly {
		fmt.Printf("\nMigrations generated successfully!\n")
		fmt.Println("\nTo run migrations:")
		fmt.Println("migrate -path migrations -database \"postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable\" up")
	} else {
		fmt.Printf("\nModule generated successfully!\n")
		fmt.Println("\nNext steps:")
		fmt.Println("1. Review the generated files")
		fmt.Println("2. Follow the setup instructions in docs/setup_" + module.NameLower + ".md")
		fmt.Println("3. Run migrations using:")
		fmt.Println("   migrate -path migrations -database \"postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable\" up")
	}

	// Update routes
	routesPath := "api/v1/routes/routes.go"
	if err := appendRoutes(routesPath, generateRoutes(module)); err != nil {
		fmt.Printf("Error updating routes: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s module!\n", module.Name)
}

func generateFile(filename, tmpl string, data Module) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	return t.Execute(f, data)
}

func toCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}

func toSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.ToLower(strings.Join(strings.Fields(s), "_"))
}

func appendRoutes(filename, routes string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read routes file: %v", err)
	}

	// Find the closing brace of the SetupRoutes function
	idx := strings.LastIndex(string(content), "}")
	if idx == -1 {
		return fmt.Errorf("could not find closing brace in routes file")
	}

	// Insert the new routes before the closing brace
	newContent := string(content[:idx]) + routes + string(content[idx:])
	return os.WriteFile(filename, []byte(newContent), 0644)
}

func generateRoutes(m Module) string {
	return fmt.Sprintf(`
	// %s routes
	v1.GET("/%s", %sHandler.GetAll)
	v1.GET("/%s/:id", %sHandler.Get%s)
	v1.POST("/%s", %sHandler.Create)
	v1.PUT("/%s/:id", %sHandler.Update)
	v1.DELETE("/%s/:id", %sHandler.Delete)
`,
		m.Name,
		m.NameLower, strings.ToLower(m.Name),
		m.NameLower, strings.ToLower(m.Name), m.Name,
		m.NameLower, strings.ToLower(m.Name),
		m.NameLower, strings.ToLower(m.Name),
		m.NameLower, strings.ToLower(m.Name))
}

const migrationUpTemplate = `CREATE TABLE IF NOT EXISTS {{.NameSnake}}s (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);`

const migrationDownTemplate = `DROP TABLE IF EXISTS {{.NameSnake}}s;`

const entityTemplate = `package entity

import (
	"gorm.io/gorm"
	"time"
)

type {{.Name}} struct {
	ID          uint           ` + "`gorm:\"primarykey\"`" + `
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt ` + "`gorm:\"index\"`" + `
	Name        string         ` + "`gorm:\"not null\"`" + `
	Description string         ` + "`gorm:\"type:text\"`" + `
	Status      string         ` + "`gorm:\"type:varchar(50);default:active\"`" + `
}

func ({{.Name}}) TableName() string {
	return "{{.NameSnake}}s"
}
`

const repositoryTemplate = `package repository

import (
	"{{.PackagePath}}/internal/entity"
	"gorm.io/gorm"
)

type {{.Name}}Repository interface {
	Create({{.NameLower}} *entity.{{.Name}}) error
	Update({{.NameLower}} *entity.{{.Name}}) error
	Delete(id uint) error
	GetByID(id uint) (*entity.{{.Name}}, error)
	GetAll() ([]entity.{{.Name}}, error)
}

type {{.NameLower}}RepositoryImpl struct {
	db *gorm.DB
}

func New{{.Name}}Repository(db *gorm.DB) {{.Name}}Repository {
	return &{{.NameLower}}RepositoryImpl{db: db}
}

func (r *{{.NameLower}}RepositoryImpl) Create({{.NameLower}} *entity.{{.Name}}) error {
	return r.db.Create({{.NameLower}}).Error
}

func (r *{{.NameLower}}RepositoryImpl) Update({{.NameLower}} *entity.{{.Name}}) error {
	return r.db.Save({{.NameLower}}).Error
}

func (r *{{.NameLower}}RepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.{{.Name}}{}, id).Error
}

func (r *{{.NameLower}}RepositoryImpl) GetByID(id uint) (*entity.{{.Name}}, error) {
	var {{.NameLower}} entity.{{.Name}}
	err := r.db.First(&{{.NameLower}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{.NameLower}}, nil
}

func (r *{{.NameLower}}RepositoryImpl) GetAll() ([]entity.{{.Name}}, error) {
	var {{.NameLower}}s []entity.{{.Name}}
	err := r.db.Find(&{{.NameLower}}s).Error
	return {{.NameLower}}s, err
}
`

const serviceTemplate = `package service

import (
	"{{.PackagePath}}/internal/entity"
	"{{.PackagePath}}/internal/repository"
)

type {{.Name}}Service interface {
	Create({{.NameLower}} *entity.{{.Name}}) error
	Update({{.NameLower}} *entity.{{.Name}}) error
	Delete(id uint) error
	GetByID(id uint) (*entity.{{.Name}}, error)
	GetAll() ([]entity.{{.Name}}, error)
}

type {{.NameLower}}ServiceImpl struct {
	repo repository.{{.Name}}Repository
}

func New{{.Name}}Service(repo repository.{{.Name}}Repository) {{.Name}}Service {
	return &{{.NameLower}}ServiceImpl{repo: repo}
}

func (s *{{.NameLower}}ServiceImpl) Create({{.NameLower}} *entity.{{.Name}}) error {
	return s.repo.Create({{.NameLower}})
}

func (s *{{.NameLower}}ServiceImpl) Update({{.NameLower}} *entity.{{.Name}}) error {
	return s.repo.Update({{.NameLower}})
}

func (s *{{.NameLower}}ServiceImpl) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *{{.NameLower}}ServiceImpl) GetByID(id uint) (*entity.{{.Name}}, error) {
	return s.repo.GetByID(id)
}

func (s *{{.NameLower}}ServiceImpl) GetAll() ([]entity.{{.Name}}, error) {
	return s.repo.GetAll()
}
`

const handlerTemplate = `package handler

import (
	"net/http"
	"strconv"
	"{{.PackagePath}}/internal/entity"
	"{{.PackagePath}}/internal/service"
	"{{.PackagePath}}/pkg/response"

	"github.com/gin-gonic/gin"
)

type {{.Name}}Handler struct {
	service service.{{.Name}}Service
}

func New{{.Name}}Handler(service service.{{.Name}}Service) *{{.Name}}Handler {
	return &{{.Name}}Handler{service: service}
}

func (h *{{.Name}}Handler) GetAll(c *gin.Context) {
	{{.NameLower}}s, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("{{.Name}}s retrieved successfully", {{.NameLower}}s))
}

func (h *{{.Name}}Handler) Get{{.Name}}(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid ID"))
		return
	}

	{{.NameLower}}, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("{{.Name}} retrieved successfully", {{.NameLower}}))
}

func (h *{{.Name}}Handler) Create(c *gin.Context) {
	var {{.NameLower}} entity.{{.Name}}
	if err := c.ShouldBindJSON(&{{.NameLower}}); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if err := h.service.Create(&{{.NameLower}}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success("{{.Name}} created successfully", {{.NameLower}}))
}

func (h *{{.Name}}Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid ID"))
		return
	}

	var {{.NameLower}} entity.{{.Name}}
	if err := c.ShouldBindJSON(&{{.NameLower}}); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	{{.NameLower}}.ID = uint(id)
	if err := h.service.Update(&{{.NameLower}}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("{{.Name}} updated successfully", {{.NameLower}}))
}

func (h *{{.Name}}Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid ID"))
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("{{.Name}} deleted successfully", nil))
}
`

const setupInstructionsTemplate = `# Setup Instructions for {{.Name}} Module

## 1. Add Entity to Auto-Migration
In cmd/api/main.go, add the following to the AutoMigrate call:

` + "```go" + `
db.AutoMigrate(
	// ... existing entities
	&entity.{{.Name}}{},
)
` + "```" + `

## 2. Initialize Repository
Add the following to the repository initialization section:

` + "```go" + `
{{.NameLower}}Repo := repository.New{{.Name}}Repository(db)
` + "```" + `

## 3. Initialize Service
Add the following to the service initialization section:

` + "```go" + `
{{.NameLower}}Service := service.New{{.Name}}Service({{.NameLower}}Repo)
` + "```" + `

## 4. Initialize Handler
Add the following to the handler initialization section:

` + "```go" + `
{{.NameLower}}Handler := handler.New{{.Name}}Handler({{.NameLower}}Service)
` + "```" + `

## 5. Update Route Setup
Update the SetupRoutes function call to include the new handler:

` + "```go" + `
routes.SetupRoutes(
	r,
	authHandler,
	postHandler,
	{{.NameLower}}Handler, // Add this line
	personalInfoHandler,
	contactInfoHandler,
	socialLinkHandler,
	postImageHandler,
)
` + "```" + `
`
