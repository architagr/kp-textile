package common

import "fmt"

func GetVendorSortKey() string {
	return VendorSortKey
}

func GetVendorContactSortKey(contactId string) string {
	return fmt.Sprintf("%s|%s", ContactSortKey, contactId)
}
