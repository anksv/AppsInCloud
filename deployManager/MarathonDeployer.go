package deployManager

import (
	"log"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	marathon "github.com/gambol99/go-marathon"
)

var marathonClient marathon.Marathon

type MarathonDeployer struct {
	//BaseURL string //will be read from app.conf or env
	Marathon marathon.Application
}

//we init the marathon client and ther configs here
func init() {

	var err error
	config := marathon.NewDefaultConfig()
	config.URL = beego.AppConfig.String("marathonURL")
	log.Println(config.URL, config)
	marathonClient, err = marathon.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create a client for marathon, error: %s", err)
	}

}

func NewMarathonDeployer() *MarathonDeployer {

	return &MarathonDeployer{}

}

func (this *MarathonDeployer) CreateWorkLoad(deps *map[string]string, data interface{}) interface{} {

	//newData, ok := data.(marathon.Application)

	newData := this.Marathon

	/*if !ok {
		log.Println("Unable CreateWorkLoad Marathon deployer")
	}*/

	log.Println(newData.Env, "Before\n\n\n")
	newData.Env = deps
	log.Println(newData.Env, "After\n\n\n")

	newAPP, err := marathonClient.CreateApplication(&newData)
	if err != nil {
		log.Printf("Failed to create application: %s, error: %s", newData, err)
		log.Println("Error in Create Application \n" + err.Error())
		return ""

	} else {
		log.Printf("Created the application: %s", newData)
		//log.Println(newAPP.String())
		//this.this["json"] = newAPP
		//this.ServeJSON()
	}

	//var x *[]marathon.PortMapping

	var returnResult string
	var ipList []string
	var flag bool
	//var err error
	if newAPP.Container != nil && newAPP.Container.Docker != nil && newAPP.Container.Docker.PortMappings != nil {
		x := *newAPP.Container.Docker.PortMappings
		log.Println("x is ", x)
		//ipList, err := marathonClient.TaskEndpoints(newData.ID, x[0].ContainerPort, false)
		ipList, err = marathonClient.TaskEndpoints(newAPP.ID, x[0].ServicePort, false)


		myGOTO:

		if err != nil {

	                log.Println("Error in TaskEnpoints", err, x[0].ServicePort,x[0].ContainerPort)
                        log.Println("Error in fetching TaskEndpoints \n" + err.Error())


			if !flag{
			flag=true
			ipList, err = marathonClient.TaskEndpoints(newData.ID, x[0].ContainerPort, false)
			goto myGOTO
			}else{

			log.Println("Error in TaskEnpoints", err)
			log.Println("Error in fetching TaskEndpoints \n" + err.Error())
				return ""
			}

		} else {
			log.Println(ipList, "MarathonDeployer: IP EndPoints",err)
			//ipList.
			//this.ServeJSON()

			for _, val := range ipList {

				returnResult = returnResult + "" + val

			}

			returnResult = "HybridApp Up and Running at " + returnResult

		}

		err = marathonClient.WaitOnApplication(newAPP.ID, (time.Second * 2))
		if err == nil {
			application, err := marathonClient.Application(newAPP.ID) //application
			if err != nil {
				log.Println(err, "Unbale to get the application")
				return nil
			} else {
				log.Println(application, " marathonCLient app\n\n\n")
				//application. +":"+strconv.Itoa(application.Container.Docker.PortMappings["servicePort"])

			}
			return returnResult
		}

		//this.this["json"] = application
		//this.ServeJSON()
	}
	//sote the this here

	//this.Ctx.Output.Body([]byte("this will be printed"))
	//this.Ctx.Output.Body(this.Ctx.Input.RequestBody)

	return returnResult

}

func (this *MarathonDeployer) DeleteWorkLoad(data interface{}) interface{} {

	newData := this.Marathon

	if depID, err := marathonClient.DeleteApplication(newData.ID, true); err != nil {
		//this.Ctx.WriteString("failed to Delete " + this.ID + " " + err.Error())
		log.Printf("Failed to delete application: %s, error: %s", newData.ID, err)
	} else {
		log.Printf("Deleted the application: %s\n", newData.ID)
		log.Println(depID)
		//this.Ctx.WriteString("Deleted " + depID.DeploymentID)
		//this.this["json"] = depID
		//this.ServeJSON()
	}
	return nil

}

func (this *MarathonDeployer) ListWorkLoad(data interface{}) interface{} {

	newData := this.Marathon
	var temp url.Values

	applications, err := marathonClient.Applications(temp)
	if err != nil {
		log.Fatalf("Failed to list applications")
	}
	//this.this["json"] = applications
	//this.ServeJSON()

	log.Printf("Found %d applications running", len(applications.Apps))
	for _, application := range applications.Apps {
		if application.ID != newData.ID {
			continue
		}
		log.Printf("Application: %s", application)
		details, err := marathonClient.Application(application.ID)
		if err != nil {
			log.Printf("Application: list error ", err)
			return nil
		}
		if details.Tasks != nil && len(details.Tasks) > 0 {
			for _, task := range details.Tasks {
				log.Printf("task: %s\n", task)
				//this.Ctx.WriteString("Task " + task.SlaveID + "\n")
			}
			// check the health of the application
			health, err := marathonClient.ApplicationOK(details.ID)
			if err != nil {
				log.Println("Unable to do the Health Check for marathon", err)
				return nil
			}
			log.Printf("Application: %s, healthy: %t\n", details.ID, health)
			//this.Ctx.WriteString(details.ID)

		}

	}
	return nil

}

func (this *MarathonDeployer) GetWorkLoadDetails(map[string]string, interface{}) *map[string]string {
	return nil
}
