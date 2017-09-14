package main

import (
	"fmt"
	"log"

	"google.golang.org/api/people/v1"
)

func main() {
	client := getOAuthClient()
	service, err := people.New(client)
	if err != nil {
		log.Fatalln("Error getting People service:", err)
	}

	deleteContacts(service)
	findIllegalGivenName(service)
}

func deleteContacts(s *people.Service) {
	contacts, err := getContacts(s)
	if err != nil {
		log.Fatalln(err)
	}

	namesMap := getNameIDMap(contacts)
	deleteRequest, err := getDeleteRequest()
	if err != nil {
		fmt.Println("No delete.txt found, skipping ...")
		return
	}

	var exists = false
	for _, req := range deleteRequest {
		if _, ok := namesMap[req]; ok {
			fmt.Print("[PRESENT] ")
			exists = true
		} else {
			fmt.Print("[   x   ] ")
		}
		fmt.Println(req)
	}

	if !exists {
		return
	}

	var confirm string
	fmt.Print("Delete names above? (Y/N) ")
	if _, err := fmt.Scan(&confirm); err != nil {
		log.Fatalln(err)
	}
	if confirm != "Y" && confirm != "y" {
		return
	}

	for _, req := range deleteRequest {
		if resourceName, ok := namesMap[req]; ok {
			_, err := s.People.DeleteContact(resourceName).Do()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println(req, "deleted")
		}
	}
}

func findIllegalGivenName(s *people.Service) {
	contacts, err := getContacts(s)
	if err != nil {
		log.Fatalln(err)
	}

	illegals := make(map[string]bool)
	illegals["Pak"] = true
	illegals["Bu"] = true
	illegals["Kak"] = true
	illegals["Om"] = true
	illegals["Tante"] = true

	for _, contact := range contacts {
		for _, name := range contact.Names {
			if !name.Metadata.Primary {
				continue
			}

			if !illegals[name.GivenName] {
				continue
			}

			fmt.Printf("Illegal: [%s] from %s\n", name.GivenName, name.DisplayName)
		}
	}
}
