package controller

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var documentServiceCtr *DocumentServiceController

type DocumentServiceController struct {
}

func InitDocumentServiceController() (*DocumentServiceController, *commonModels.ErrorDetail) {
	if documentServiceCtr == nil {

		documentServiceCtr = &DocumentServiceController{}
	}
	return documentServiceCtr, nil
}

func (ctrl *DocumentServiceController) GetChallan(context *gin.Context) {
	var data commonModels.InventoryDto
	if err := context.ShouldBindJSON(&data); err == nil {

		// var data commonModels.InventoryDto = commonModels.InventoryDto{
		// 	BranchId:         "branchId",
		// 	InventorySortKey: "Inventory|Deleted|sales01|Inventory|Sales|1643142960",
		// 	PurchaseDate:     time.Date(2022, 12, 25, 0, 0, 0, 0, time.Local),
		// 	SalesDate:        time.Date(2022, 12, 25, 0, 0, 0, 0, time.Local),
		// 	BillNo:           "sales01",
		// 	LrNo:             "123",
		// 	ChallanNo:        "123",
		// 	HsnCode:          "sdsdf",
		// 	BailDetails: []commonModels.BailDetailsDto{
		// 		{
		// 			BailNo:         "bail01",
		// 			BilledQuantity: 1000,
		// 			Quality:        "4729adb5-7432-11ec-a804-0800275114e0",
		// 		},
		// 	},
		// 	TransporterId: "daac2ed7-7a3f-11ec-bda7-0800275114e0",
		// }
		var total int32 = 0
		var bailInfo = ""
		for i, val := range data.BailDetails {
			total = total + val.BilledQuantity
			bailInfo = bailInfo + `<tr>
		<td>` + fmt.Sprintf("%d", (i+1)) + `</td>
		<td>` + val.BailNo + `</td>
		<td>1</td>
		<td>` + fmt.Sprintf("%d", val.BilledQuantity) + `</td>
	</tr>`
		}
		var htmlBody = `<html>
	<head>
		<style>
        @media print {
            @page {
                size: A4;
                margin: 20px;
            }
            .page{
                page-break-after: always;
				border:1px solid black;
				height:98vh;
            }
			.page .header{
				padding:5px;
				border-bottom:1px solid black;
				height:10vh;
				display: flex;
				flex-direction: column;
				justify-content: space-between;
			}
			.page .header .top-heading{
				display: flex;
				justify-content: space-between;			
			}
			.page .header .center-heading{
				display: flex;
				justify-content: center;
				text-align: center;
				font-size:22px;
			} 
			.page .header .end-heading{
				display: flex;
				justify-content: center;
				text-align: center;
			}
			.page .challan-info{
				padding:10px;
				border-bottom:1px solid black;
				height:8vh;
				display: flex;
				flex-direction: column;
				justify-content: space-between;
			}
			.page .challan-info .top{
				display: flex;
				justify-content: space-between;	
				border-bottom:1px dashed  black;
			}
			.page .challan-info .center{
				
				border-bottom:1px dashed  black;
			}
			.page .challan-info .end{
				
				border-bottom:1px dashed  black;
			}
			
			.page .main-content{
				display: flex;
				flex-direction: row;
				height:77vh
			}
			.page-new .main-content{
				display: flex;
				flex-direction: row;
				height:98vh
			}
			.page .main-content .bale-info{
				border-right:1px solid  black;
				 width: 70%;
			}
			
			.page .main-content .bale-info, .other-info{
				align-self: stretch;
				
			}
			.page .main-content .other-info{
				display: flex;
				flex-direction: column;
				justify-content: stretch;
				width:30%;
			}
			.page .main-content .other-info .quality{
				padding:10px;
				vertical-align: text-top;
				height:30vh;
				border-bottom:1px solid black;
				
			}
			.page .main-content .other-info .remarks{
				padding:10px;
				vertical-align: text-top;
				height:30vh;
				border-bottom:1px solid black;
			}
			.page .main-content .other-info .note{
				padding:10px;
				vertical-align: text-top;
			}
			.page .main-content .bale-info table{
			margin-right:-1px;
				width:100%;
			}
			.page .main-content .bale-info table{
				 border-collapse: collapse;
			}
			.page .main-content .bale-info table, th, td {
				border: 1px solid black;
			}
			.page .main-content .bale-info table, td {
				text-align:center;
			}
        }
		</style>
		<title>` + data.BillNo + `</title>
	</head>

	<body>
		<div class="page">
			<div class="header">
				<div class="top-heading">
					<div class="left">All Subject to mumbai Juridiction</div>
					<div class="center">|| Shree Ganeshay Nanaha ||</div>
					<div class="center">Mob: +91 9322 226 411</div>
				</div>
				<div class="center-heading">
					<b>G.K. Syntext Pvt. Ltd.</b>
				</div>
				<div class="end-heading">
				Swadeshi Market, Room No. K 2nd Floor, Kalbadevi Road, Mumbai - 400 002.
					<br/>
				Email: agawalgaurav1993@gmail.com
				</div>
			</div>
			<div class="challan-info">
				<div class="top">
					<div class="left">Challan Number : ` + data.ChallanNo + `</div>
					<div class="right">Date:  ` + data.SalesDate.Format("02 Jan 2006") + `</div>
				</div>
				<div class="center">
				Transport: ` + data.TransporterId + `
				</div>
				<div class="end">
				L.R. No.: ` + data.LrNo + `
				</div>
			</div>
			<div class="main-content">
				<div class="bale-info">
					<table>
						<thead>
							<th>
								S. No.
							</th>
							<th>
								Bale No.
							</th>
							<th>
								Pcs.
							</th>
							<th>
								Meters
							</th>
						</thead>
						<tbody>
							` + bailInfo + `
						</tbody>
						<tfoot>
							<th>Total</th>
							<th></th>
							<th></th>
							<th>` + fmt.Sprintf("%d", total) + `</th>
						</tfoot>
					</table>
				</div>
				<div class="other-info">
					<div class="quality">
						<u>Quality :</u>
					</div>
					<div class="remarks">
						<u>Other Remarks :</u>
					</div>
					<div class="note">
						<u>Note :</u>
					</div>
				</div>
			</div>
		</div>
		<script>
			window.print();
		</script>


	</body>

</html>`

		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlBody))
	} else {
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<html><body>An error has occured</body></html>"))
	}
}
