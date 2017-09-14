package main

import (
	"bufio"
	"google.golang.org/api/people/v1"
	"os"
)

func getContacts(service *people.Service) ([]*people.Person, error) {
	var contacts []*people.Person
	var nextPage string

	for {
		call := service.People.Connections.List("people/me")
		call = call.PersonFields("names").SortOrder("FIRST_NAME_ASCENDING").PageSize(1000)
		if nextPage != "" {
			call = call.PageToken(nextPage)
		}
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, resp.Connections...)

		nextPage = resp.NextPageToken
		if nextPage == "" {
			break
		}
	}

	return contacts, nil
}

func getNameIDMap(p []*people.Person) map[string]string {
	m := make(map[string]string)

	for _, person := range p {
		for _, name := range person.Names {
			if name.Metadata.Primary {
				m[name.DisplayName] = person.ResourceName
			}
		}
	}

	return m
}

func getDeleteRequest() ([]string, error) {
	var names []string

	f, err := os.Open("delete.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return names, nil
}
