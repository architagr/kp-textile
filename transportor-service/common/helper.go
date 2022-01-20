package common

import "fmt"

func GetTransporterSortKey(vendorId string) string {
	return fmt.Sprintf("%s|%s", TransporterSortKey, vendorId)
}

func GetTransporterContactSortKey(vendorId, contactId string) string {
	return fmt.Sprintf("%s|%s|%s", ContactSortKey, vendorId, contactId)
}
