package models

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_sales/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/helpers"
	"github.com/gofiber/fiber/v2"
)

func Fetch_crm(username, status string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_crm
	var arraobj []entities.Model_crm
	var res helpers.Responsemovie
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 150
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(A.idcrmsales) as totalmember  "
	sql_selectcount += "FROM " + configs.DB_tbl_trx_crmsales + " as A  "
	sql_selectcount += "JOIN " + configs.DB_tbl_trx_usersales + " as B ON B.phone = A.phone  "
	sql_selectcount += "WHERE A.username='" + username + "' "
	sql_selectcount += "AND statususersales='" + status + "' "

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "A.idcrmsales, A.phone,  "
	sql_select += "B.nama, B.idusersales,  "
	sql_select += "createcrmsales, to_char(COALESCE(createdatecrmsales,NOW()), 'YYYY-MM-DD HH24:MI:SS'), "
	sql_select += "updatecrmsales, to_char(COALESCE(updatedatecrmsales,NOW()) , 'YYYY-MM-DD HH24:MI:SS')"
	sql_select += "FROM " + configs.DB_tbl_trx_crmsales + " as A  "
	sql_select += "JOIN " + configs.DB_tbl_trx_usersales + " as B ON B.phone = A.phone  "
	sql_select += "WHERE username='" + strings.ToLower(username) + "' "
	sql_select += "AND statususersales='" + status + "' "
	sql_select += "ORDER BY createcrmsales DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcrmsales_db, idusersales_db                                                      int
			phone_db, nama_db                                                                  string
			createcrmsales_db, createdatecrmsales_db, updatecrmsales_db, updatedatecrmsales_db string
		)

		err = row.Scan(
			&idcrmsales_db, &phone_db, &nama_db, &idusersales_db,
			&createcrmsales_db, &createdatecrmsales_db, &updatecrmsales_db, &updatedatecrmsales_db)
		helpers.ErrorCheck(err)

		create := ""
		update := ""
		if createcrmsales_db != "" {
			create = createcrmsales_db + ", " + createdatecrmsales_db
		}
		if updatecrmsales_db != "" {
			update = updatecrmsales_db + ", " + updatedatecrmsales_db
		}
		obj.Crm_idcrmsales = idcrmsales_db
		obj.Crm_idusersales = idusersales_db
		obj.Crm_name = nama_db
		obj.Crm_phone = phone_db
		obj.Crm_create = create
		obj.Crm_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
