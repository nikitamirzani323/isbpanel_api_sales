package controllers

import (
	"log"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_sales/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/helpers"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const Fieldcrm_home_redis = "LISTCRM_SALES_ISBPANEL"

func Crmhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crm)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}

	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")
	log.Println(client_admin)

	var obj entities.Model_crm
	var arraobj []entities.Model_crm
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcrm_home_redis + "_" + client_admin + "_" + client.Crm_status)
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		crm_idcrmsales, _ := jsonparser.GetInt(value, "crm_idcrmsales")
		crm_idusersales, _ := jsonparser.GetInt(value, "crm_idusersales")
		crm_phone, _ := jsonparser.GetString(value, "crm_phone")
		crm_name, _ := jsonparser.GetString(value, "crm_name")
		crm_create, _ := jsonparser.GetString(value, "crm_create")
		crm_update, _ := jsonparser.GetString(value, "crm_update")

		obj.Crm_idcrmsales = int(crm_idcrmsales)
		obj.Crm_idusersales = int(crm_idusersales)
		obj.Crm_phone = crm_phone
		obj.Crm_name = crm_name
		obj.Crm_create = crm_create
		obj.Crm_update = crm_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_crm(client_admin, client.Crm_status)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcrm_home_redis+"_"+client_admin+"_"+client.Crm_status, result, 60*time.Minute)
		log.Println("CRM SALES MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CRM SALES  CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
func CrmSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	//admin, status, status_dua, phone, note, sData string, idusersales, idcrmsales, idwebagen, deposit int
	result, err := models.Save_crm(
		client_admin,
		client.Crm_status, client.Crm_status_dua, client.Crm_phone, client.Crm_note,
		client.Sdata, client.Crm_idusersales, client.Crm_idcrmsales, client.Crm_idwebagen, client.Crm_deposit)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_crm(client_admin, "PROCESS")
	return c.JSON(result)
}
func _deleteredis_crm(admin, status string) {
	val_master := helpers.DeleteRedis(Fieldcrm_home_redis + "_" + admin + "_" + status)
	log.Printf("Redis Delete BACKEND CRM : %d", val_master)
}
