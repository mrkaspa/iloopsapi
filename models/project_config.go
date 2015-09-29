package models

//ProjectConfig json struct
type Loops struct {
	CronFormat string `json:"cron_format"`
}
type ProjectConfig struct {
	Name       string `json:"name"`
	AppID      string `json:"app_id"`
	MainScript string `json:"main_script"`
	Loops      Loops  `json:"loops"`
}
