package utils

import (
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

type AppUUId struct {
	uuid.UUID
}

func GetUUID() (val string) {

	val = ""
	uuid4, err := uuid.NewV4()
	if err != nil {
		log.Println("Util: AppUUId Error:", err)
		return
	}
	log.Println("Generated UUID is ", uuid4)
	return uuid4.String()
}
