package logger

import (
	"assign/common"
	"assign/config"
	"assign/types"
)

func ClearDB() error {
	query := "DELETE FROM info"

	_, err := config.GlobalConfig.DB.Exec(query)
	if err != nil {
		common.StopProgram(err)
		return err
	}

	return nil
}

func RegistPerson(ID string) error {
	query := "INSERT INTO info (studentID ) VALUES (?)"

	_, err := config.GlobalConfig.DB.Exec(query, ID)
	if err != nil {
		common.StopProgram(err)
	}

	return nil
}

func RegistInfo(info types.Info) error {
	query := "UPDATE info SET deploymentName=?, serviceName=?, nodePort=? WHERE studentID=?"

	_, err := config.GlobalConfig.DB.Exec(query, info.DeploymentName, info.ServiceName, info.NodePort, info.StudentID)
	if err != nil {
		common.StopProgram(err)
	}

	return nil
}

func UpdateInfo(info types.Info, attr string) error {
	var query string
	var err error

	switch attr {
	case "deploymentName":
		query = "UPDATE info SET deploymentName=? WHERE studentID=?"
		_, err = config.GlobalConfig.DB.Exec(query, info.DeploymentName, info.StudentID)
	case "serviceName":
		query = "UPDATE info SET serviceName=? WHERE studentID=?"
		_, err = config.GlobalConfig.DB.Exec(query, info.ServiceName, info.StudentID)
	case "nodePort":
		query = "UPDATE info SET nodePort=? WHERE studentID=?"
		_, err = config.GlobalConfig.DB.Exec(query, info.NodePort, info.StudentID)
	}

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
	query := "SELECT nodePort, deploymentName, serviceName FROM info WHERE studentID = ?"

	var info types.Info
	var err error
	err = config.GlobalConfig.DB.QueryRow(query, ID).Scan(&info.NodePort, &info.DeploymentName, &info.ServiceName)
	if err != nil {
		common.StopProgram(err)
	}

	return info
}

func GetStudentID(info types.Info, attr string) string {
	var query string
	var err error

	switch attr {
	case "deploymentName":
		query = "SELECT studentID FROM info WHERE deploymentName=?"
		err = config.GlobalConfig.DB.QueryRow(query, info.DeploymentName).Scan(&info.StudentID)
	case "serviceName":
		query = "SELECT studentID FROM info WHERE serviceName=?"
		err = config.GlobalConfig.DB.QueryRow(query, info.ServiceName).Scan(&info.StudentID)
	case "nodePort":
		query = "SELECT studentID FROM info WHERE nodePort=?"
		err = config.GlobalConfig.DB.QueryRow(query, info.NodePort).Scan(&info.StudentID)
	}

	if err != nil {
		common.StopProgram(err)
	}
	return info.StudentID
}

func IsIDExist(ID string) bool {
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
