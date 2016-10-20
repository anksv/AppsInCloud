package deployManager

import (
	"errors"
	"log"
	"reflect"
)

// All rule typr must implement this rule
type DeployerInterface interface {
	CreateWorkLoad(*map[string]string, interface{}) interface{}
	DeleteWorkLoad(interface{}) interface{}
	ListWorkLoad(interface{}) interface{}
	GetWorkLoadDetails(map[string]string, interface{}) *map[string]string
}

var ErrTimeoutError = errors.New("the operation has timed out")

type RedisApp struct {
	AppId    string
	Capacity string
	Master   string
	Slave    string
}

//AtomicApp is a simple app in the Hybrid appp
type AtomicApp struct {
	AppName string
	Type    string //continer,database
	SubType string //marathon,redis
	//AppId       string      //actual app id
	Application interface{} // app data marathon.Application || RedisApp
}

type HybridApp struct {
	Applications []AtomicApp
	Dependencies []Dependency //string wil be the nodename
}

// TODO: SHIV
// Need to resolve circular dependency later
type Dependency struct {
	DependentAppName string
	AppName          string            //this is the dependent app
	Environment      map[string]string //this wil be the expected o/p param from the app
	//the env map will tell the app what need to be fetched from the dependent app
	// TODO: vert simple version example DBNAME ="redis.mesos.com" PORT = "6565" etc..
}

func SearchApplication(newApp *HybridApp, name string) *AtomicApp {
	for key, app := range newApp.Applications {

		if app.AppName == name {
			return &newApp.Applications[key]
		}

	}

	log.Fatalln("SearchApplication: App ", name, " not found")
	return nil
}

func CreateWorkLoad(newApp *HybridApp) interface{} {
	//resolve dependencies

	depEnv := make(map[string]*map[string]string)

	doneMap := make(map[string]bool)

	for key, dep := range newApp.Dependencies {

		log.Println(dep.DependentAppName, "Depends on ", dep.AppName)
		//deply dep.AppName first
		//newApp.Applications[dep.AppName]  //create this app

		log.Println(reflect.TypeOf(newApp.Applications[key].Application))

		//search for the proper function

		foundApp := SearchApplication(newApp, dep.AppName)
		log.Println(reflect.TypeOf(foundApp.Application))

		appHandle, ok := foundApp.Application.(DeployerInterface)
		if ok {

			//thiss, _ := appHandle.(RedisDeployer)
			//log.Println(reflect.TypeOf(appHandle), thiss.AppId, thiss)
			result := appHandle.CreateWorkLoad(nil, appHandle) // this will create the redis instance

			if result == nil {
				log.Println("CreateWorkLoad: No result")
				return nil
			}

			envMap := appHandle.GetWorkLoadDetails(dep.Environment, result)
			depEnv[dep.DependentAppName] = envMap

			doneMap[dep.AppName] = true

			log.Println("env from redis app", depEnv[dep.AppName], envMap)
			//we need the response also
			//we need to take the relevant result and use it
		} else {
			log.Fatalln("DeployManager:CreateWorkLoad not a Deloyer type redis ", newApp.Dependencies[key].DependentAppName, dep.AppName)
		}

	}

	for _, applicationHandle := range newApp.Applications { //appName

		if _, ok := doneMap[applicationHandle.AppName]; ok {
			continue
		}

		appHandle, ok := applicationHandle.Application.(DeployerInterface)
		if ok {
			result := appHandle.CreateWorkLoad(depEnv[applicationHandle.AppName], appHandle) // this will create the redis instance

			if result == nil {
				log.Println("CreateWorkLoad: FAiled to create workload")
				return nil
			} else {
				return result
			}

			//envMap := appHandle.GetWorkLoadDetails(dep.Environment, result)
			//depEnv[appname] = envMap
			//we need the response also
			//we need to take the relevant result and use it
		} else {
			log.Println("DeployManager:CreateWorkLoad not a Deloyer type ")
		}

	}

	return nil

	//parse dependencies

}

func DeleteWorkLoad(newApp *HybridApp) {
	//resolve dependencies

	for appname, dep := range newApp.Dependencies {

		log.Println(appname, "Depends on ", dep.AppName)
		//deply dep.AppName first
		//newApp.Applications[dep.AppName]  //create this app

		appHandle, ok := newApp.Applications[appname].Application.(DeployerInterface) // TODO:SHIV
		if ok {
			appHandle.DeleteWorkLoad(appHandle) // this will create the redis instance
			//we need the response also
		}

	}

	//parse dependencies

}

/*func WaitOnApplication(name string, timeout time.Duration) error {
	if r.appExistAndRunning(name) {
		return nil
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-time.After(timeout):
			return ErrTimeoutError
		case <-ticker.C:
			if r.appExistAndRunning(name) {
				return nil
			}
		}
	}
}*/

func DeleteApp(newApp *HybridApp) {
}

func ListApp(newApp *HybridApp) {

}
