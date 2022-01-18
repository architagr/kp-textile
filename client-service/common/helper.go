package common

import "fmt"

func GetClientSortKey(clientId string) string {
	return fmt.Sprintf("%s|%s", ClientSortKey, clientId)
}

func GetClientContactSortKey(clientId, contactId string) string {
	return fmt.Sprintf("%s|%s|%s", ContactSortKey, clientId, contactId)
}
