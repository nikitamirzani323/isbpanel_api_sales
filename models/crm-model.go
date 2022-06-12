package models

import (
	"context"
	"log"
	"strconv"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_sales/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_crm(username, status string) (helpers.Response, error) {
	var obj entities.Model_crm
	var arraobj []entities.Model_crm
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		A.idcrmsales, A.phone, B.nama, B.idusersales,
		A.createcrmsales, to_char(COALESCE(A.createdatecrmsales,NOW()), 'YYYY-MM-DD HH24:MI:SS'),
		A.updatecrmsales, to_char(COALESCE(A.updatedatecrmsales,NOW()) , 'YYYY-MM-DD HH24:MI:SS')
		FROM ` + configs.DB_tbl_trx_crmsales + `  as A 
		JOIN ` + configs.DB_tbl_trx_usersales + `  as B ON B.phone = A.phone
		WHERE A.username=$1 
		AND B.statususersales=$2 
		ORDER BY A.createdatecrmsales ASC    
	`

	row, err := con.QueryContext(ctx, sql_select, username, status)
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
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_crm(admin, status, status_dua, phone, username, note, sData string, idusersales, idcrmsales, idwebagen int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	sql_update_parent := `
			UPDATE 
			` + configs.DB_tbl_trx_usersales + `  
			SET statususersales=$1, 
			updateusersales=$2, updatedateusersales=$3 
			WHERE idusersales =$4 
		`

	flag_update, msg_update := Exec_SQL(sql_update_parent, configs.DB_tbl_trx_usersales, "UPDATE",
		status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idusersales)

	if flag_update {
		log.Println(msg_update)

		sql_update_detail := `
			UPDATE 
			` + configs.DB_tbl_trx_crmsales + `  
			SET statuscrmsales_satu=$1, statuscrmsales_dua=$2, notecrmsales=$3, 
			updatecrmsales=$4, updatedatecrmsales=$5  
			WHERE idcrmsales=$6  
		`

		flag_update_detail, msg_update_detail := Exec_SQL(sql_update_detail, configs.DB_tbl_trx_crmsales, "UPDATE",
			status, status_dua, note, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcrmsales)

		if flag_update_detail {
			log.Println(msg_update_detail)

			sql_insert := `
				insert into
				` + configs.DB_tbl_trx_usersales_log + ` (
					idusersaleslog , idusersales, idcrmsales, statuslog, 
					createusersaleslog, createdateusersaleslog 
				) values (
					$1, $2, $3, $4, 
					$5, $6 
				)
			`
			field_column := configs.DB_tbl_trx_usersales_log + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_usersales_log, "INSERT",
				tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idusersales, idcrmsales, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
				log.Println(msg_insert)

				if status_dua == "DEPOSIT" { // DEPOSIT
					sql_insert := `
						insert into
						` + configs.DB_tbl_trx_usersales_deposit + ` (
							idusersalesdeposit , idcrmsales, idwebagen, phone, username, 
							createusersalesdeposit, createdateusersalesdeposit  
						) values (
							$1, $2, $3, $4, $5,
							$6, $7 
						)
					`
					field_column := configs.DB_tbl_trx_usersales_deposit + tglnow.Format("YYYY")
					idrecord_counter := Get_counter(field_column)
					flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_usersales_deposit, "INSERT",
						tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idcrmsales, idwebagen, phone, username,
						admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

					if flag_insert {
						msg = "Succes"
						log.Println(msg_insert)
					} else {
						log.Println(msg_insert)
					}
				}
			} else {
				log.Println(msg_insert)
			}
		} else {
			log.Println(msg_update_detail)
		}
	} else {
		log.Println(msg_update)
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
