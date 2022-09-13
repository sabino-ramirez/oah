package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sabino-ramirez/oah/models"
)

func GetProjectRequisitions(client *models.Client, target interface{}) (int, error) {
	url := fmt.Sprintf("https://lab-services.ovation.io/api/v3/project_templates/%d/requisitions?startDate=01-01-2020&endDate=01-01-2021", client.ProjectTemplateId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("new req error:", err)
	}
	req.Header.Add("Authorization", client.Bearer)

	res, err := client.Http.Do(req)
	if err != nil {
		log.Println("response error:", err)
	}

	defer res.Body.Close()

	return res.StatusCode, json.NewDecoder(res.Body).Decode(target)
}

func GetProjectTemplates(client *models.Client, target interface{}) (int, error) {
	url := fmt.Sprintf("https://lab-services.ovation.io/api/v3/project_templates?organizationId=%d", client.OrganizationId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("new req error:", err)
	}
	req.Header.Add("Authorization", client.Bearer)

	res, err := client.Http.Do(req)
	if err != nil {
		log.Println("response error:", err)
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		fmt.Printf("The response is 404 which means project templates for orgId: %d doesn't exist. \nTry different one?\n", client.OrganizationId)
		// body, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		//   return err
		// }
		// log.Println(string([]byte(body)))
	}

	return res.StatusCode, json.NewDecoder(res.Body).Decode(target)
}
