package common

import "fmt"

const (
	SORTKEY_BAILDETAILS_STOCK = "Stock"
	SORTKEY_BAILDETAILS_SOLD  = "Sold"
)

const (
	STATUS_PURCHASE_STOCK = "Stock"
	STATUS_PURCHASE_SOLD  = "Sold"
)

// const (
// 	SORTKEY_BAILDETAILS_PURCHASE   = "Bale|Purchase"
// 	SORTKEY_BAILDETAILS_SALES      = "Bale|Sales"
// 	SORTKEY_BAILDETAILS_OUTOFSTOCK = "Bale|OutOfStock"
// 	SORTKEY_BAILDETAILS_DELETED    = "Bale|Deleted"
// )

// const (
// 	SORTKEY_BAILINFO   = "Info"
// 	SORTKEY_BAILDELETE = "Deleted"
// )

// const (
// 	SORTKEY_INVENTORY_PURCHASE = "Inventory|Purchase"
// 	SORTKEY_INVENTORY_SALES    = "Inventory|Sales"
// 	SORTKEY_INVENTORY_DELETE   = "Inventory|Deleted"
// )

// func GetInventoryPurchanseSortKey(billNo string) string {
// 	return fmt.Sprintf("%s|%s", SORTKEY_INVENTORY_PURCHASE, billNo)
// }
// func GetInventorySalesSortKey(billNo string) string {
// 	return fmt.Sprintf("%s|%s", SORTKEY_INVENTORY_SALES, billNo)
// }

// func GetInventoryDeleteSortKey(billNo string) string {
// 	return fmt.Sprintf("%s|%s", SORTKEY_INVENTORY_DELETE, billNo)
// }

// func GetBaleInfoSortKey(baleNo string) string {
// 	return fmt.Sprintf("%s|%s", SORTKEY_BAILINFO, baleNo)
// }
// func GetBaleInfoDeleteSortKey(baleNo string) string {
// 	return fmt.Sprintf("%s|%s", SORTKEY_BAILDELETE, baleNo)
// }

// func GetBaleDetailPurchanseSortKey(quality, baleNo string) string {
// 	return fmt.Sprintf("%s|%s|%s", SORTKEY_BAILDETAILS_PURCHASE, quality, baleNo)
// }
// func GetBaleDetailSalesSortKey(quality, baleNo, salesBillNumber string) string {
// 	return fmt.Sprintf("%s|%s|%s|%s", SORTKEY_BAILDETAILS_SALES, quality, baleNo, salesBillNumber)
// }
// func GetBaleDetailOutOfStockSortKey(quality, baleNo string) string {
// 	return fmt.Sprintf("%s|%s|%s", SORTKEY_BAILDETAILS_OUTOFSTOCK, quality, baleNo)
// }
// func GetBaleDetailDeleteSortKey(quality, baleNo string) string {
// 	return fmt.Sprintf("%s|%s|%s", SORTKEY_BAILDETAILS_DELETED, quality, baleNo)
// }

func GetPurchaseSortKey(productId, qualityId, purchaseId string) string {
	return fmt.Sprintf("%s|%s|%s", productId, qualityId, purchaseId)
}

func GetSalesSortKey(productId, qualityId, purchaseId string) string {
	return fmt.Sprintf("%s|%s|%s", productId, qualityId, purchaseId)
}

func GetInStockBaleSortKey(productId, qualityId, baleNo string) string {
	return fmt.Sprintf("%s|%s|%s|%s", SORTKEY_BAILDETAILS_STOCK, productId, qualityId, baleNo)
}
func GetSoldBaleSortKey(productId, qualityId, baleNo string) string {
	return fmt.Sprintf("%s|%s|%s|%s", SORTKEY_BAILDETAILS_SOLD, productId, qualityId, baleNo)
}
