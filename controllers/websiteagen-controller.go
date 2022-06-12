package controllers

import (
	"log"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_sales/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/helpers"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/models"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
)

const Fieldwebsiteagen_home_redis = "LISTWEBSITEAGEN_SALES_ISBPANEL"

func Websiteagenhome(c *fiber.Ctx) error {
	var obj entities.Model_websiteagen
	var arraobj []entities.Model_websiteagen
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldwebsiteagen_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		websiteagen_id, _ := jsonparser.GetInt(value, "websiteagen_id")
		websiteagen_name, _ := jsonparser.GetString(value, "websiteagen_name")

		obj.Websiteagen_id = int(websiteagen_id)
		obj.Websiteagen_name = websiteagen_name
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_websiteagenHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldwebsiteagen_home_redis, result, 1*time.Hour)
		log.Println("WEBSITEAGEN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("WEBSITEAGEN CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     message_RD,
			"record":      arraobj,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
