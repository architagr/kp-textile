package common

import "fmt"

const (
	SORTKEY_BAILDETAILS_PURCHASE   = "Bail|Purchase"
	SORTKEY_BAILDETAILS_SALES      = "Bail|Sales"
	SORTKEY_BAILDETAILS_OUTOFSTOCK = "Bail|OutOfStock"
	SORTKEY_BAILDETAILS_DELETED    = "Bail|Deleted"
)

const (
	SORTKEY_BAILINFO   = "Info"
	SORTKEY_BAILDELETE = "Deleted"
)

const (
	SORTKEY_INVENTORY_PURCHASE = "Inventory|Purchase"
	SORTKEY_INVENTORY_SALES    = "Inventory|Sales"
	SORTKEY_INVENTORY_DELETE   = "Inventory|Deleted"
)

func GetInventoryPurchanseSortKey(billNo string) string {
	return fmt.Sprintf("%s|%s", SORTKEY_INVENTORY_PURCHASE, billNo)
}
func GetInventorySalesSortKey(billNo string) string {
	return fmt.Sprintf("%s|%s", SORTKEY_INVENTORY_SALES, billNo)
}

func GetInventoryDeleteSortKey(billNo string) string {
	return fmt.Sprintf("%s|%s", SORTKEY_INVENTORY_DELETE, billNo)
}

func GetBailInfoSortKey(bailNo string) string {
	return fmt.Sprintf("%s|%s", SORTKEY_BAILINFO, bailNo)
}
func GetBailInfoDeleteSortKey(bailNo string) string {
	return fmt.Sprintf("%s|%s", SORTKEY_BAILDELETE, bailNo)
}

func GetBailDetailPurchanseSortKey(quality, bailNo string) string {
	return fmt.Sprintf("%s|%s|%s", SORTKEY_BAILDETAILS_PURCHASE, quality, bailNo)
}
func GetBailDetailSalesSortKey(quality, bailNo, salesBillNumber string) string {
	return fmt.Sprintf("%s|%s|%s|%s", SORTKEY_BAILDETAILS_SALES, quality, bailNo, salesBillNumber)
}
func GetBailDetailOutOfStockSortKey(quality, bailNo string) string {
	return fmt.Sprintf("%s|%s|%s", SORTKEY_BAILDETAILS_OUTOFSTOCK, quality, bailNo)
}
func GetBailDetailDeleteSortKey(quality, bailNo string) string {
	return fmt.Sprintf("%s|%s|%s", SORTKEY_BAILDETAILS_DELETED, quality, bailNo)
}
