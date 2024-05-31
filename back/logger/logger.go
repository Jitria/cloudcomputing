package logger

import (
	"assign/common"
	"assign/config"
	"assign/types"
	"database/sql"
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
	query := "UPDATE info SET podName=?, serviceName=?, nodePort=? WHERE studentID=?"

	_, err := config.GlobalConfig.DB.Exec(query, info.PodName, info.ServiceName, info.NodePort, info.StudentID)
	if err != nil {
		common.StopProgram(err)
	}

	return nil
}

func UpdateInfo(info types.Info, attr string) error {
	var query string
	var err error

	switch attr {
	case "podName":
		query = "UPDATE info SET podName=? WHERE studentID=?"
		_, err = config.GlobalConfig.DB.Exec(query, info.PodName, info.StudentID)
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
	query := "SELECT nodePort, podName, serviceName FROM info WHERE studentID = ?"

	var info types.Info
	var err error
	err = config.GlobalConfig.DB.QueryRow(query, ID).Scan(&info.NodePort, &info.PodName, &info.ServiceName)
	if err != nil {
		common.StopProgram(err)
	}

	return info
}

func GetpodNames() []string {
	query := "SELECT DISTINCT podName FROM info"

	rows, err := config.GlobalConfig.DB.Query(query)
	if err != nil {
		common.StopProgram(err)
	}
	defer rows.Close()

	var podNames []string
	for rows.Next() {
		var podName string
		err := rows.Scan(&podName)
		if err != nil {
			common.StopProgram(err)
		}
		podNames = append(podNames, podName)
	}
	if err := rows.Err(); err != nil {
		common.StopProgram(err)
	}

	return podNames
}

func GetserviceNames() []string {
	query := "SELECT DISTINCT serviceName FROM info"

	rows, err := config.GlobalConfig.DB.Query(query)
	if err != nil {
		common.StopProgram(err)
	}
	defer rows.Close()

	var serviceNames []string
	for rows.Next() {
		var serviceName string
		err := rows.Scan(&serviceName)
		if err != nil {
			common.StopProgram(err)
		}
		serviceNames = append(serviceNames, serviceName)
	}
	if err := rows.Err(); err != nil {
		common.StopProgram(err)
	}

	return serviceNames
}

func GetStudentID(info types.Info, attr string) (string, error) {
	var query string
	var err error

	switch attr {
	case "podName":
		query = "SELECT studentID FROM info WHERE podName=?"
		err = config.GlobalConfig.DB.QueryRow(query, info.PodName).Scan(&info.StudentID)
	case "serviceName":
		query = "SELECT studentID FROM info WHERE serviceName=?"
		err = config.GlobalConfig.DB.QueryRow(query, info.ServiceName).Scan(&info.StudentID)
	case "nodePort":
		query = "SELECT studentID FROM info WHERE nodePort=?"
		err = config.GlobalConfig.DB.QueryRow(query, info.NodePort).Scan(&info.StudentID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		common.StopProgram(err)
	}
	return info.StudentID, nil
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
