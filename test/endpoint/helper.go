package endpoint

import (
	"encoding/json"

	"bitbucket.org/kiloops/api/models"
	"github.com/jinzhu/gorm"
)

type empty struct{}

var emptyJSON, _ = json.Marshal(empty{})

func saveUser() models.User {
	user := models.User{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
	models.Gdb.InTx(func(txn *gorm.DB) {
		if txn.Save(&user).Error != nil {
			panic("error creating the user")
		}
	})
	return user
}

func addSSH(user *models.User) models.SSH {
	sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCadJM0DdJotRnSWW7coFcCxMW1cCZIJqkyW3wMQoUOU2VHuLExh44tDpSAiz2EEeFlqJ5hI67ZI+3bSx7puKSr44l78H/Kb8UDLidAUao7JZoo0thq7bAVesGr+8aligmULvxH3sQqstI9yNcifJ56jHUVTB14PslBmhA56pmGOva0ojmdt9l2aBy4LxQBDc5Js+AcPlfC2zXE7rtaiafB/M3992V+7CEisbAv7CpsI3SPdpW2p4mfR1zMVpf4Jt6lQJW6Sr53/bzAP4/Tif3fgbZhoSL8qnnLi3556gWi90FwFhCoqqDR/lN3sxJQx5NxxCF8mbNgpmS5qDptFyF michel.ingesoft@gmail.com`
	ssh := models.SSH{PublicKey: sshKey, UserID: user.ID}
	models.Gdb.InTx(func(txn *gorm.DB) {
		if txn.Save(&ssh).Error != nil {
			panic("error creating the ssh")
		}
	})
	return ssh
}
