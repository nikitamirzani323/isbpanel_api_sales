package routers

import (
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/controllers"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Init() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		// c.Set("Content-Security-Policy", "frame-ancestors 'none'")
		// c.Set("X-XSS-Protection", "1; mode=block")
		// c.Set("X-Content-Type-Options", "nosniff")
		// c.Set("X-Download-Options", "noopen")
		// c.Set("Strict-Transport-Security", "max-age=5184000")
		// c.Set("X-Frame-Options", "SAMEORIGIN")
		// c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Get("/ipaddress", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      "data",
			"BASEURL":     c.BaseURL(),
			"HOSTNAME":    c.Hostname(),
			"IP":          c.IP(),
			"IPS":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomain":   c.Subdomains(),
		})
	})
	app.Get("/dashboard", monitor.New())

	app.Post("/api/login", controllers.CheckLogin)
	app.Post("/api/valid", middleware.JWTProtected(), controllers.Home)
	app.Post("/api/alladmin", middleware.JWTProtected(), controllers.Adminhome)
	app.Post("/api/detailadmin", middleware.JWTProtected(), controllers.AdminDetail)
	app.Post("/api/saveadmin", middleware.JWTProtected(), controllers.AdminSave)

	app.Post("/api/alladminrule", middleware.JWTProtected(), controllers.Adminrulehome)
	app.Post("/api/saveadminrule", middleware.JWTProtected(), controllers.AdminruleSave)

	app.Post("/api/domain", middleware.JWTProtected(), controllers.Domainhome)
	app.Post("/api/domainsave", middleware.JWTProtected(), controllers.DomainSave)

	app.Post("/api/game", middleware.JWTProtected(), controllers.Gamehome)
	app.Post("/api/gamesave", middleware.JWTProtected(), controllers.Gamesave)

	app.Post("/api/crm", middleware.JWTProtected(), controllers.Crmhome)
	app.Post("/api/crmsales", middleware.JWTProtected(), controllers.Crmsales)
	app.Post("/api/crmisbtv", middleware.JWTProtected(), controllers.Crmisbtvhome)
	app.Post("/api/crmduniafilm", middleware.JWTProtected(), controllers.Crmduniafilm)
	app.Post("/api/crmsave", middleware.JWTProtected(), controllers.CrmSave)
	app.Post("/api/crmsavestatus", middleware.JWTProtected(), controllers.CrmSavestatus)
	app.Post("/api/crmsalessave", middleware.JWTProtected(), controllers.CrmSalesSave)
	app.Post("/api/crmsalesdelete", middleware.JWTProtected(), controllers.CrmSalesdelete)
	app.Post("/api/crmsavesource", middleware.JWTProtected(), controllers.CrmSavesource)

	app.Post("/api/departement", middleware.JWTProtected(), controllers.Departementhome)
	app.Post("/api/departementsave", middleware.JWTProtected(), controllers.DepartementSave)
	app.Post("/api/employee", middleware.JWTProtected(), controllers.Employeehome)
	app.Post("/api/employeebydepart", middleware.JWTProtected(), controllers.EmployeeByDepart)
	app.Post("/api/employeesave", middleware.JWTProtected(), controllers.EmployeeSave)
	return app
}
