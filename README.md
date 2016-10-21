# AppDeployment
We will using this code to deploy apps using DCOS, marathon on mesos orchestration Engine

1. Change the marathon url and redis url in the app.conf file.
2. Build the code using go build
3. Run the Binary ./paas-apiserver
4. Pass server will be running on port 8081
5. Create Redis Instance using below command
   curl -X "POST" http://10.11.12.150:5656/v1/CREATE/TestInstance/100/1/2
   Request Accepted, Instance will be created.
6. Check the Status
   curl -v http://10.11.12.150:5656/v1/STATUS/TestInstance
7. Try to connect and deploy an app with the JSON using POSTMASTER
 
