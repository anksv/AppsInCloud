package deployManager

//Redis deployer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

type RedisStatusResponse struct {
	Name     string
	Type     string
	Status   string
	Capacity int
	Master   *RedisInstanceStatus
	Slaves   []*RedisInstanceStatus
}

type RedisInstanceStatus struct {
	IP                 string
	Port               string
	MemoryCapacity     int
	MemoryUsed         int64
	Uptime             int64
	ClientsConnected   int
	LastSyncedToMaster int
}

const (
	CREATE_SUFFIX = "/v1/CREATE"
	DELETE_SUFFIX = "/v1/DELETE"
	STATUS_SUFFIX = "/v1/STATUS"
)

type RedisDeployer struct {
	AppId    string
	Capacity string
	Master   string
	Slave    string
}

var BaseURL string

func init() {
	BaseURL = beego.AppConfig.String("redisFrameWorkURL")
}
func NewRedisDeployer() *RedisDeployer {

	return &RedisDeployer{}

}

func init() {
	//we nuild the redis url here or read from the beego config

	// we work with one redis connection

}

//deps map[string]string,
func (newData *RedisDeployer) CreateWorkLoad(deps *map[string]string, data interface{}) interface{} {

	log.Println("RedisDeployer: CreateWorkLoad called", newData, data)

	//newDataa := data.(RedisDeployer)

	thiss, ok := data.(*RedisDeployer)
	this := thiss

	if !ok {
		log.Fatalln("Unable CreateWorkLoad RedisDeployer")

	}

	log.Println("RedisDeployer: CreateWorkLoad called", newData, data, thiss)

	req, err := http.NewRequest("POST", BaseURL+""+CREATE_SUFFIX+"/"+this.AppId+"/"+this.Capacity+"/"+this.Master+"/"+this.Slave, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("\n\nresponse Body:", string(body))

	//we need to wait till the redis instance is create

	res := this.WaitOnApplication(time.Second * 35)
	if res != nil {
		//we got the redis instance deployed we will pass the result to the deploymnagaer
		return res
	}

	return nil

}

func (this *RedisDeployer) appRunning() interface{} {
	//initial status check
	resultVal := this.ListWorkLoad(nil)
	redisResponse, ok := resultVal.(RedisStatusResponse)
	if ok {
		//check fields if not empty
		if redisResponse.Master != nil {
			if redisResponse.Master.IP != "" && redisResponse.Master.Port != "" {
				return redisResponse
			}
		}

	} else {
		log.Println("RedisDeployer: Unable to get the response")
	}

	return nil
}

//wait the get the status on the application deployed
func (this *RedisDeployer) WaitOnApplication(timeout time.Duration) interface{} {

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-time.After(timeout):
			return ErrTimeoutError
		case <-ticker.C:
			if res := this.appRunning(); res != nil {
				return res
			}
		}
	}
}

func (this *RedisDeployer) DeleteWorkLoad(data interface{}) interface{} {

	req, err := http.NewRequest("DELETE", BaseURL+""+DELETE_SUFFIX+"/"+this.AppId, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return nil

}

func (this *RedisDeployer) ListWorkLoad(data interface{}) interface{} {

	var localResponse RedisStatusResponse

	req, err := http.NewRequest("GET", BaseURL+""+STATUS_SUFFIX+"/"+this.AppId, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &localResponse)

	if err != nil {
		log.Println("ListWorkLoad: Unmarshall json response failed", err)
		return nil
	}

	log.Println("response Body:", string(body), localResponse)
	return localResponse

}

func (this *RedisDeployer) GetWorkLoadDetails(resultMaps map[string]string, result interface{}) *map[string]string {

	var resMap map[string]string

	resMap = make(map[string]string)

	res, ok := result.(RedisStatusResponse)
	if !ok {
		log.Println("GetWorkLoadDetails: Unable to parse the response")
		return nil
	}

	for key, value := range resultMaps {

		switch value {
		case "Master_IP":
			resMap[key] = res.Master.IP
		case "Master_Port":
			resMap[key] = res.Master.Port
		}

	}

	return &resMap

}
