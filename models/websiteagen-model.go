package models

import (
	"context"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_sales/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/helpers"
	"github.com/gofiber/fiber/v2"
)

func Fetch_websiteagenHome() (helpers.Response, error) {
	var obj entities.Model_websiteagen
	var arraobj []entities.Model_websiteagen
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idwebagen , nmwebagen  
			FROM ` + configs.DB_tbl_mst_websiteagen + `  
			ORDER BY nmwebagen ASC    
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idwebagen_db int
			nmwebagen_db string
		)

		err = row.Scan(&idwebagen_db, &nmwebagen_db)

		helpers.ErrorCheck(err)

		obj.Websiteagen_id = idwebagen_db
		obj.Websiteagen_name = nmwebagen_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
