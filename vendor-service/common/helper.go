package common

import "fmt"

func GetVendorSortKey(vendorId string) string {
	return fmt.Sprintf("%s|%s", VendorSortKey, vendorId)
}

func GetVendorContactSortKey(vendorId, contactId string) string {
	return fmt.Sprintf("%s|%s|%s", ContactSortKey, vendorId, contactId)
}
