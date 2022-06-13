package entities

import "encoding/json"

type Model_crm struct {
	Crm_idcrmsales  int    `json:"crm_idcrmsales"`
	Crm_idusersales int    `json:"crm_idusersales"`
	Crm_phone       string `json:"crm_phone"`
	Crm_name        string `json:"crm_name"`
	Crm_create      string `json:"crm_create"`
	Crm_update      string `json:"crm_update"`
}
type Model_crmsales_simple struct {
	Crmsales_username     string `json:"crmsales_username"`
	Crmsales_nameemployee string `json:"crmsales_nameemployee"`
}
type Model_crmsales struct {
	Crmsales_id           int    `json:"crmsales_id"`
	Crmsales_phone        string `json:"crmsales_phone"`
	Crmsales_namamember   string `json:"crmsales_namamember"`
	Crmsales_username     string `json:"crmsales_username"`
	Crmsales_nameemployee string `json:"crmsales_nameemployee"`
	Crmsales_create       string `json:"crmsales_create"`
	Crmsales_update       string `json:"crmsales_update"`
}
type Model_crmisbtv struct {
	Crmisbtv_username  string `json:"crmisbtv_username"`
	Crmisbtv_name      string `json:"crmisbtv_name"`
	Crmisbtv_coderef   string `json:"crmisbtv_coderef"`
	Crmisbtv_point     int    `json:"crmisbtv_point"`
	Crmisbtv_status    string `json:"crmisbtv_status"`
	Crmisbtv_lastlogin string `json:"crmisbtv_lastlogin"`
	Crmisbtv_create    string `json:"crmisbtv_create"`
	Crmisbtv_update    string `json:"crmisbtv_update"`
}
type Model_crmduniafilm struct {
	Crmduniafilm_username string `json:"crmduniafilm_username"`
	Crmduniafilm_name     string `json:"crmduniafilm_name"`
}

type Controller_crm struct {
	Crm_status string `json:"crm_status"`
}
type Controller_crmsales struct {
	Crmsales_phone string `json:"crmsales_phone"`
}
type Controller_crmisbtv struct {
	Crmisbtv_search string `json:"crmisbtv_search"`
	Crmisbtv_page   int    `json:"crmisbtv_page"`
}
type Controller_crmsave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Crm_idusersales int    `json:"crm_idusersales" validate:"required"`
	Crm_idcrmsales  int    `json:"crm_idcrmsales" validate:"required"`
	Crm_idwebagen   int    `json:"crm_idcwebagen"`
	Crm_status      string `json:"crm_status" validate:"required"`
	Crm_status_dua  string `json:"crm_status_dua" `
	Crm_note        string `json:"crm_note" `
	Crm_phone       string `json:"crm_phone" `
	Crm_iduseragen  string `json:"crm_iduseragen" `
	Crm_deposit     int    `json:"crm_deposit"`
}

type Controller_crmstatussave struct {
	Page       string `json:"page" validate:"required"`
	Crm_page   int    `json:"crm_page"`
	Crm_id     int    `json:"crm_id"`
	Crm_status string `json:"crm_status" validate:"required"`
}
type Controller_crmsalessave struct {
	Page              string `json:"page" validate:"required"`
	Search            string `json:"search" `
	Crm_page          int    `json:"crm_page"`
	Crmsales_phone    string `json:"crmsales_phone" validate:"required"`
	Crmsales_username string `json:"crmsales_username" validate:"required"`
}
type Controller_crmsavesource struct {
	Page       string          `json:"page" validate:"required"`
	Sdata      string          `json:"sdata" validate:"required"`
	Crm_page   int             `json:"crm_page"`
	Crm_source string          `json:"crm_source" `
	Crm_data   json.RawMessage `json:"crm_data" validate:"required"`
}
type Controller_crmsalesdelete struct {
	Page           string `json:"page" validate:"required"`
	Search         string `json:"search" `
	Crm_page       int    `json:"crm_page"`
	Crmsales_id    int    `json:"crmsales_id" validate:"required"`
	Crmsales_phone string `json:"crmsales_phone" validate:"required"`
}
