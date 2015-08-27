package models

//ProjectConfig json struct
type Loops struct {
	CronFormat string `json:"cron_format"`
}
type ProjectConfig struct {
	Name  string `json:"name"`
	AppID string `json:"app_id"`
	Loops Loops  `json:"loops"`
}
