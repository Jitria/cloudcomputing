package logger

import (
	"assign/common"
	"assign/config"
	"assign/types"
)

func RegistPerson(ID string) error {
	query := "INSERT INTO info (studentID ) VALUES (?)"

	_, err := config.GlobalConfig.DB.Exec(query, ID)
	if err != nil {
		common.StopProgram(err)
	}

	return nil
}

func UpdateInfo(info types.Info) error {
	query := "UPDATE info SET deploymentName=?, serviceName=?, nodePort=? WHERE studentID=?"

	_, err := config.GlobalConfig.DB.Exec(query, info.DeploymentName, info.ServiceName, info.NodePort, info.StudentID)
	if err != nil {
		common.StopProgram(err)
	}

	return nil
}

func DeletePerson(ID string) error {
	query := "DELETE FROM info WHERE studentID = ?"

	_, err := config.GlobalConfig.DB.Exec(query, ID)
	if err != nil {
		common.StopProgram(err)
		return err
	}

	return nil
}

func GetInfo(ID string) types.Info {
	query := "SELECT nodePort, deploymentName, serviceName, ip FROM info WHERE studentID = ?"

	var info types.Info
	var err error
	err = config.GlobalConfig.DB.QueryRow(query, ID).Scan(&info.NodePort, &info.DeploymentName, &info.ServiceName, &info.Ip)
	if err != nil {
		common.StopProgram(err)
	}

	return info
}

func IsExist(ID string) bool {
	query := "SELECT COUNT(*) FROM info WHERE studentID = ?"

	var count int
	var err error
	err = config.GlobalConfig.DB.QueryRow(query, ID).Scan(&count)
	if err != nil {
		common.StopProgram(err)
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}
