package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/stretchr/objx"
)

func checktoken(token string) string {
	client := &http.Client{}
	var link = "https://discordapp.com/api/v7/users/@me"
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return string(bodyString)
}

func main() {
	contents, err := ioutil.ReadFile("tokens.txt")
	if err != nil {
		log.Fatal(err)
	}
	list := string(contents)
	tokens := strings.Split(list, "\r\n")
	if len(tokens) < 2 {
		color.Red("You need 1 or more tokens to run this!")
	} else {
		color.Red(`Printing format: 

	Token: valid, id, username, email, verified, phone`)
		for _, token := range tokens {
			resp := checktoken(token)
			var dresp map[string]interface{}
			err := json.Unmarshal([]byte(resp), &dresp)
			if err != nil {
				log.Fatal(err)
			}
			if dresp["message"] != "401: Unauthorized" {
				if dresp["email"] == nil {
					dresp["email"] = "none"
				}
				o := objx.New(dresp)
				verified := o.Get("verified")
				color.Green("%s: Valid, %s, %s#%s, %s, %s, %s", token, dresp["id"], dresp["username"], dresp["discriminator"], dresp["email"], verified, dresp["phone"])
			} else {
				color.Red("%s: INVALID", token)
			}
		}
	}
}
