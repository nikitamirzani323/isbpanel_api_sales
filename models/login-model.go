package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_sales/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_sales/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Login_Model(username, password, ipaddress string) (bool, string, error) {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	tglnow, _ := goment.New()
	var password_db, iddepartement_db, nmdepart_db string
	sql_select := `
			SELECT
			A.password, A.iddepartement, B.nmdepartement     
			FROM ` + configs.DB_tbl_mst_employee + ` as A 
			JOIN ` + configs.DB_tbl_mst_departement + ` as B ON B.iddepartement = A.iddepartement  
			WHERE A.username  = $1
			AND A.iddepartement = 'SLS' 
			AND A.statusemployee = 'Y' 
		`

	log.Println(sql_select, username)
	row := con.QueryRowContext(ctx, sql_select, username)
	switch e := row.Scan(&password_db, &iddepartement_db, &nmdepart_db); e {
	case sql.ErrNoRows:
		return false, "", errors.New("Username and Password Not Found")
	case nil:
		flag = true
	default:
		return false, "", errors.New("Username and Password Not Found")
	}

	hashpass := helpers.HashPasswordMD5(password)
	log.Println("Password : " + hashpass)
	log.Println("Hash : " + password_db)
	if hashpass != password_db {
		return false, "", nil
	}

	if flag {
		sql_update := `
			UPDATE ` + configs.DB_tbl_mst_employee + ` 
			SET lastlogin=$1, ipaddress=$2 ,
			updateemployee=$3,  updatedateemployee=$4   
			WHERE username = $5 
			AND statusemployee = 'Y' 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_employee, "UPDATE",
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			ipaddress,
			username, tglnow.Format("YYYY-MM-DD HH:mm:ss"), username)

		if flag_update {
			flag = true
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	return true, iddepartement_db, nil
}
func UpdatePassword_Model(username, password string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	hashpass := helpers.HashPasswordMD5(password)
	sql_update := `
		UPDATE ` + configs.DB_tbl_mst_employee + ` 
		SET password=$1, 
		updateemployee=$2,  updatedateemployee=$3   
		WHERE username = $4 
		AND iddepartement = 'SLS'  
		AND statusemployee = 'Y' 
	`
	flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_employee, "UPDATE",
		hashpass, username, tglnow.Format("YYYY-MM-DD HH:mm:ss"), username)

	if flag_update {
		msg = "Succes"
		log.Println(msg_update)
	} else {
		log.Println(msg_update)
	}
	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()
	return res, nil
}
