package common

import "fmt"

func GetTransporterSortKey() string {
	return TransporterSortKey
}

func GetTransporterContactSortKey(contactId string) string {
	return fmt.Sprintf("%s|%s", ContactSortKey, contactId)
}
