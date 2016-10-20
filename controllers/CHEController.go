package controllers

import (
	"encoding/json"
	"log"

	"../deployManager"

	"github.com/astaxie/beego"
)

type CHEController struct {
	beego.Controller
}

//we init the marathon client and ther configs here
func init() {

}

func (this *CHEController) CreateApp() {

	var newhybridApp deployManager.HybridApp

	var dummyhybridApp deployManager.HybridApp
	log.Println(string(this.Ctx.Input.RequestBody))

	//Unmarshall the hybridapp and send it to the DeployManager

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &newhybridApp)

	if err != nil {
		this.Ctx.WriteString("CreateWorkLoad: Unable to read the json " + err.Error())
		return
	}

	//dummyhybridApp.Applications = make(map[string]deployManager.AtomicApp)
	//dummyhybridApp.Dependencies = make(map[string]deployManager.Dependency)

	//mapref := dummyhybridApp.Applications
	/*
	   	var tmp = m["foo"]
	   tmp.x = 4
	   m["foo"] = tmp*/

	//maratha := &deployManager.MarathonDeployer{}
	//baremetal := &deployManager.RedisDeployer{}

	var maratha deployManager.MarathonDeployer
	var baremetal deployManager.RedisDeployer

	log.Println("Marathon and baremetal before ", maratha, baremetal, "\n\n\n\n")

	for key, val := range newhybridApp.Applications {

		if val.Type == "Container" {
			//var xx = mapref[key]
			//xx.Application = &maratha
			//mapref[key] = xx
			newhybridApp.Applications[key].Application = &deployManager.MarathonDeployer{}

		} else {
			//var xx = mapref[key]
			//xx.Application = &baremetal
			//mapref[key] = xx
			newhybridApp.Applications[key].Application = &deployManager.RedisDeployer{}
		}
	}

	log.Println(dummyhybridApp, "\n\n\n", newhybridApp, "First \n\n\n")
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &newhybridApp)

	if err != nil {
		this.Ctx.WriteString("CreateWorkLoad: Unable to read the json " + err.Error())
		return
	}

	log.Println(dummyhybridApp, "\n\n\n", newhybridApp, "Second\n\n\n")
	log.Println("Marathon and baremetal After ", maratha, baremetal, "\n\n\n\n")

	/*mapref = dummyhybridApp.Applications
	for key, val := range newhybridApp.Applications {

		if val.Type == "Container" {
			var xx = mapref[key]
			xx.Application = &maratha
			mapref[key] = xx
		} else {
			var xx = mapref[key]
			xx.Application = &baremetal
			mapref[key] = xx
		}
	}*/

	//log.Fatalln(dummyhybridApp, "\n\n\n", newhybridApp, "Second\n\n\n")

	res := deployManager.CreateWorkLoad(&newhybridApp)
	this.Data["json"] = res.(string)
	this.ServeJSON()

}

func (this *CHEController) DeleteApp() {

	var newhybridApp deployManager.HybridApp

	log.Println(string(this.Ctx.Input.RequestBody))

	//Unmarshall the hybridapp and send it to the DeployManager

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &newhybridApp)

	if err != nil {
		this.Ctx.WriteString("CreateWorkLoad: Unable to read the json" + err.Error())
		return
	}

	deployManager.DeleteWorkLoad(&newhybridApp)

}

func (this *CHEController) ListAPP() {

	//var temp url.Values

	////applications, err := marathonClient.Applications(temp)
	//if err != nil {
	//	log.Fatalf("Failed to list applications")
	//}
	//	this.Data["json"] = applications
	//	this.ServeJSON()

	/*log.Printf("Found %d applications running", len(applications.Apps))
	for _, application := range applications.Apps {
		log.Printf("Application: %s", application)
		details, err := marathonClient.Application(application.ID)
		if err != nil {
			log.Printf("Application: list error ", err)
			return
		}
		if details.Tasks != nil && len(details.Tasks) > 0 {
			for _, task := range details.Tasks {
				log.Printf("task: %s\n", task)
				this.Ctx.WriteString("Task " + task.SlaveID + "\n")
			}
			// check the health of the application
			health, err := marathonClient.ApplicationOK(details.ID)
			if err != nil {
				log.Println("Unable to do the Health Check for marathon", err)
				return
			}
			log.Printf("Application: %s, healthy: %t\n", details.ID, health)
			this.Ctx.WriteString(details.ID)

		}

	}*/
}
