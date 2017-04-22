package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
)

func main() {
}

func getClient() *api.Client {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Printf("%v", err)
	}

	return c
}

func readFromSrc(c api.Client, p string) *api.Secret {
	s, err := c.Logical().Read(p)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return s
}

func listFromSrc(c *api.Client, p string) *api.Secret {
	l, err := c.Logical().List(p)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return l
}

func writeToDest(c api.Client, p string, s api.Secret) {
	_, err := c.Logical().Write(p, s.Data)
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func mv(source string, dest string) (string, error) {
	var msg string
	var sep string

	client := getClient()
	listSecret := listFromSrc(client, source)

	if listSecret == nil {
		srcSecret := readFromSrc(*client, source)

		if srcSecret != nil {
			writeToDest(*client, dest, *srcSecret)
			msg = fmt.Sprintf("Moved %s to %s\n", source, dest)
		}

	} else {
		list := listSecret.Data["keys"]
		if keys, ok := list.([]interface{}); ok {
			for _, key := range keys {
				keyStr := fmt.Sprintf("%v", key)

				//if strings.HasSuffix(keyStr, "/") {
				if strings.HasSuffix(source, "/") {
					sep = ""
				} else {
					sep = "/"
				}

				absSource := source + sep + keyStr
				absDest := dest + sep + keyStr

				keyMsg, _ := mv(absSource, absDest)
				msg = keyMsg + msg
			}
		} else {
			fmt.Printf("list not a []interface{}: %v\n", list)
		}
	}

	return msg, nil
}

/*
## Failures
1. create a client
2. authenticate
3. read the contents of secret/re/foo
	- if permission denied
		- Log "Could not move from secret/re/foo. You are not permitted to access the source path."
	- else
		- save the data from the contents
4. write the data to secret/teams/re/foo
	- if permission denied
		- Log "Could not move to secret/teams/re/foo. You are not permitted to access the destination path."
	- else
		- Log "
		- exit 0
*/
