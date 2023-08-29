package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"bytes"
	"strconv"
	_"unsafe"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)





//==============================================================================================================================
//	 Patent Chain Structure Definitions
//==============================================================================================================================

type SmartContract struct {
}

type PO struct{
	POID					 string					`json:"poID"`
	REF_POID				 string					`json:"refPOID"`
	OriginationOrgId		 string     			`json:"originationOrgId"`
	OriginationOrgName     	 string					`json:"originationOrgName"`
	OriginationOrgType     	 string					`json:"originationOrgType"`
	DestinationOrgId 		 string					`json:"destinationOrgId"`
	DestinationOrgName		 string					`json:"destinationOrgName"`
	DestinationOrgType     	 string					`json:"destinationOrgType"`
	PONumber				 string					`json:"poNumber"`
	Product_ID				 string					`json:"productID"`
	ProductDesc				 string					`json:"productDesc"`
	RawMaterialSourcing		 string					`json:"rawmaterialSourcing"`
	ZipCode					 string					`json:"zipCode"`
	CreatedDate				 string					`json:"createdDate"`
	ExpectedDate			 string					`json:"expectedDate"`
	Specification			 string					`json:"specification"`
	UOM						 string					`json:"UOM"`
	Quantity				 uint64					`json:"quantity"`
	Comments				 string					`json:"comments"`
	CreatedBy				 string					`json:"createdBy"`
	Status					 string					`json:"status"`
	Status_Update			 string					`json:"statusUpdate"`
	TxId					 string					`json:"tx_id"`
	DOC_TYPE				 string					`json:"docType"`
	UPDATED_DATE			 string 				`json:"updatedDate"`
}


type Invoice struct{
	Invoice_ID					 string					`json:"invoiceID"`
	Invoice_Number				 string					`json:"invoiceNumber"`
	OriginationOrgId			 string     			`json:"originationOrgId"`
	OriginationOrgName 	    	 string					`json:"originationOrgName"`
	OriginationOrgType  	   	 string					`json:"originationOrgType"`
	DestinationOrgId 			 string					`json:"destinationOrgId"`
	DestinationOrgName			 string					`json:"destinationOrgName"`
	DestinationOrgType   	  	 string					`json:"destinationOrgType"`
	REF_POID					 string					`json:"refPOID"`
	PONumber					 string					`json:"poNumber"`
	RawMaterialSourcing			 string					`json:"rawmaterialSourcing"`
	Transporter_ID				 string					`json:"transporterID"`
	Transporter_Name			 string					`json:"transporterName"`
	Shipping_Doc				 string					`json:"shippingDoc"`
	Product_ID					 string					`json:"productID"`
	ProductDesc					 string					`json:"productDesc"`
	Heat_Index	 			   	 string					`json:"heatIndex"`
	Lab_Certificate		         string					`json:"labCert"`
	Lab_Certificate_Hash         string					`json:"labCertHash"`
	Lab_Certificate_Mime         string					`json:"labCertMime"`
	Lab_Certificate_Name         string					`json:"labCertName"`
	Quantity					 uint64					`json:"quantity"`
	Quantity_Accepted			 uint64					`json:"quantityAccepted"`
	UOM							 string					`json:"UOM"`
	RM_Certificate				 string					`json:"rmCert"`
	Comments					 string					`json:"comments"`
	Status						 string					`json:"status"` //Created || Accepted || Rejected || Transported_Shipped
	Status_Update				 string					`json:"statusUpdate"` // To filter records if update only comments
	TxId						 string					`json:"tx_id"`
	DOC_TYPE					 string					`json:"docType"`
	CreatedBy					 string					`json:"createdBy"`
	UpdatedBy					 string					`json:"updatedBy"`
	CreatedDate					 string					`json:"createdDate"`
	UPDATED_DATE				 string 				`json:"updatedDate"`
}

type Product struct{
	DOC_TYPE				 string					`json:"docType"`
	Product_Key				 string					`json:"productKey"`
	Product_ID				 string					`json:"productID"`
	ProductDesc				 string					`json:"productDesc"`
	OrgId					 string     			`json:"orgId"`
	OrgName					 string     			`json:"orgName"`
	OrgType  			   	 string					`json:"orgType"`
	Quantity				 uint64					`json:"quantity"`
	TxId					 string					`json:"tx_id"`
}

//Common structure to maintain a product quantity a organization have (like it will be the quantity accepted by Eng Company so that it cannot be able to ship more than that) 
type Product_Quantity struct{
	OrgId					 string     			`json:"orgId"`
	OrgName					 string     			`json:"orgName"`
	Product_ID				 string					`json:"productID"`
	ProductDesc				 string					`json:"productDesc"`
	Quantity				 uint64					`json:"quantity"`
}

//Common struture to maintain the product quantity being shipped against PO (Used in createPO , createInvoice )
type Quantity struct{
	POID						 string					`json:"poID"`
	Quantity_Ordered			 uint64					`json:"quantityOrdered"`
	Quantity_Sent				 uint64					`json:"quantitySent"`
	Quantity_Accepted			 uint64					`json:"quantityAccepted"`
	Quantity_Balance			 uint64					`json:"quantityBalance"`
}

type Notification struct{
	Notification_ID						 string					`json:"notificationID"`
	Status								 string					`json:"status"` 
	Type								 string					`json:"type"` //PO //Dispatch Note
	OrgId								 string     			`json:"orgId"`
	CreatedBy							 string					`json:"createdBy"`
	TxId								 string					`json:"tx_id"`
	DOC_TYPE							 string					`json:"docType"`
	PONumber							 string					`json:"poNumber"`
	Invoice_Number						 string					`json:"invoiceNumber"`
}


/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createPO" {
		return s.createPO(APIstub, args)
	}else if function == "updatePOStatus" {
		return s.updatePOStatus(APIstub, args)
	}else if function == "queryPOHistory" {
		return s.queryPOHistory(APIstub, args)
	}else if function == "queryPOByPONumber" {
		return s.queryPOByPONumber(APIstub, args)
	}else if function == "queryPOByRefPOID" {
		return s.queryPOByRefPOID(APIstub, args)
	}else if function == "queryPOsByKey" {
		return s.queryPOsByKey(APIstub, args)
	}else if function == "queryPOsByStatus" {
		return s.queryPOsByStatus(APIstub, args)
	}else if function == "queryPOsByProductID" {
		return s.queryPOsByProductID(APIstub, args)
	}else if function == "queryPOByOrgID" {
		return s.queryPOByOrgID(APIstub, args)
	}else if function == "queryAllPOs" {
		return s.queryAllPOs(APIstub, args)
	}else if function == "queryPOHistoryWithoutMetadata" {
		return s.queryPOHistoryWithoutMetadata(APIstub, args)
	}else if function == "updatePOComments" {
		return s.updatePOComments(APIstub, args)
	}else if function == "createProduct" {
		return s.createProduct(APIstub, args)
	}else if function == "addProductStock" {
		return s.addProductStock(APIstub, args)
	}else if function == "queryProductByOrgID" {
		return s.queryProductByOrgID(APIstub, args)
	}else if function == "queryProductByID" {
		return s.queryProductByID(APIstub, args)
	}else if function == "queryAllProducts" {
		return s.queryAllProducts(APIstub, args)
	}else if function == "queryPOQuantityDetails" {
		return s.queryPOQuantityDetails(APIstub, args)
	}else if function == "createInvoice" {
		return s.createInvoice(APIstub, args)
	}else if function == "updateInvoiceStatus" {
		return s.updateInvoiceStatus(APIstub, args)
	}else if function == "updateInvoiceComments" {
		return s.updateInvoiceComments(APIstub, args)
	}else if function == "queryInvoiceByID" {
		return s.queryInvoiceByID(APIstub, args)
	}else if function == "queryInvoiceByInvoiceNumber" {
		return s.queryInvoiceByInvoiceNumber(APIstub, args)
	}else if function == "queryInvoiceByStatus" {
		return s.queryInvoiceByStatus(APIstub, args)
	}else if function == "queryInvoiceByOrgID" {
		return s.queryInvoiceByOrgID(APIstub, args)
	}else if function == "queryNotificationByOrgID" {
		return s.queryNotificationByOrgID(APIstub, args)
	}else if function == "queryNotificationByID" {
		return s.queryNotificationByID(APIstub, args)
	}else if function == "queryInvoiceByPOID" {
		return s.queryInvoiceByPOID(APIstub, args)
	}else if function == "queryInvoiceByPONumber" {
		return s.queryInvoiceByPONumber(APIstub, args)
	}else if function == "recallInvoice" {
		return s.recallInvoice(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createNotification(APIstub shim.ChaincodeStubInterface,notificationID string, status string, notificationType string, orgId string, createdBy string, poNumber string, invoiceNumber string) string {
	
	var notification Notification
	notification.Notification_ID = notificationID
	notification.Status = status
	notification.Type = notificationType
	notification.OrgId = orgId
	notification.CreatedBy = createdBy
	notification.TxId = APIstub.GetTxID()
	notification.DOC_TYPE = "Notification"
	notification.PONumber = poNumber
	notification.Invoice_Number = invoiceNumber

	key := "NOTIFICATION_"+ notificationID +"_"+ notification.TxId
	notificationAsBytes, _ := json.Marshal(notification)
	APIstub.PutState(key, notificationAsBytes)
	return "Notification Created successfully"
}


func (s *SmartContract) queryNotificationByOrgID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgId := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Notification\",\"orgId\":\"%s\"}}", orgId)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryNotificationByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	notificationID := args[0]
	notificationAsBytes, _ := APIstub.GetState(notificationID)
	return shim.Success(notificationAsBytes)
}



//TODOS
/*
1. Create Product Done 
2. Raise Stock (Increase Quantity) Done
*/

func (s *SmartContract) createProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//Parsing request
	ob := args[0]
    product := Product{}
	err := json.Unmarshal([]byte(ob), &product)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	//Manufacturer Check
	if product.OrgType != "Manufacturer" {
		return shim.Error("Only manufacturer can create Product Catlogue")
	}

	product.DOC_TYPE = "Product";
	product.TxId = APIstub.GetTxID()
	product.Product_Key = product.Product_ID + "_" +  product.OrgId
	productAsBytes, _ := json.Marshal(product)
	APIstub.PutState(product.Product_Key, productAsBytes)
	return shim.Success([]byte("Product Created successfully"))
}


func (s *SmartContract) addProductStock(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//Parsing request
	ob := args[0]
    product := Product{}
	err := json.Unmarshal([]byte(ob), &product)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	
	//Manufacturer Check
	if product.OrgType != "Manufacturer" {
		return shim.Error("Only manufacturer can update the Quantity of Product")
	}

	//Fetch Product Details
	productkey := product.Product_ID + "_" + product.OrgId
	productAsBytes, _ := APIstub.GetState(productkey)
	var productObject Product
	if len(productAsBytes) != 0 {
		err := json.Unmarshal(productAsBytes, &productObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}else{
		return shim.Error("Not a valid organization to add Stock as product is owned by different organization")
	}

	//Update Product Quantity to Ledger
	productObject.TxId = APIstub.GetTxID()
	productObject.Quantity = productObject.Quantity + product.Quantity
	productAsBytesNew, _ := json.Marshal(productObject)
	APIstub.PutState(productObject.Product_Key, productAsBytesNew)
	return shim.Success([]byte("Product Created successfully"))
}

func (s *SmartContract) queryProductByOrgID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgId := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Product\",\"orgId\":\"%s\"}}", orgId)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryProductByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productID := args[0]
	orgType := "Manufacturer"
	queryString := "{\"selector\":{\"docType\":\"Product\",\"$and\": [{\"orgType\":\"" + orgType +"\"},{\"productID\":\""+productID+"\"}]}}"
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryAllProducts(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgType := "Manufacturer"
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Product\",\"orgType\":\"%s\"}}", orgType)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


func (s *SmartContract) queryPOQuantityDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	poID := args[0]
	poQuantityKey := "PO"+poID
	poAsBytes, _ := APIstub.GetState(poQuantityKey)
	return shim.Success(poAsBytes)
}

/*
	1. Create Invoice by shipping the Quantity Done
	2. Create Quantity Ob for tracking Quantity Raise, accepted, rejected Done
	3  PO raised Quantity Check Done 
	4  Org Quantity Check to check Org Holdings Pending
 */
func (s *SmartContract) createInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    invoice := Invoice{}
	err := json.Unmarshal([]byte(ob), &invoice)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	invoice.DOC_TYPE = "INVOICE";
	invoice.TxId = APIstub.GetTxID()
	invoice.Quantity_Accepted = 0;
	
	//Fetch PO Ordered Quantity
	poQuantityKey := "PO"+invoice.REF_POID
	quantityObAsBytes, _ := APIstub.GetState(poQuantityKey)
	var poQuantityObject Quantity
	if len(quantityObAsBytes) != 0 {
		err := json.Unmarshal(quantityObAsBytes, &poQuantityObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}

	//Fetch Product Details
	//4  Org Quantity Check to check Org Holdings Pending
	productkey := invoice.Product_ID + "_" + invoice.OriginationOrgId
	productAsBytes, _ := APIstub.GetState(productkey)
	var productObject Product
	if len(productAsBytes) != 0 {
		err := json.Unmarshal(productAsBytes, &productObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
		if(productObject.Quantity < invoice.Quantity){
			return shim.Error("You dont have required Product Quantity which you are shipping")
		}
	}
	//Update Product Stock once dispatched 
	productObject.Quantity = productObject.Quantity - invoice.Quantity
	productAsBytesNew, _ := json.Marshal(productObject)
	APIstub.PutState(productkey, productAsBytesNew)
	

	//Check If Manfacturer shipped more than the Ordered Quantity
	if invoice.Quantity > poQuantityObject.Quantity_Balance {
		return shim.Error("You cannot send products more than the PO Quantity")
	}

	poQuantityObject.Quantity_Sent = poQuantityObject.Quantity_Sent + invoice.Quantity
	quantityObAsBytesNew, _ := json.Marshal(poQuantityObject)
	APIstub.PutState(poQuantityObject.POID, quantityObAsBytesNew)


	invoiceAsBytes, _ := json.Marshal(invoice)
	APIstub.PutState(invoice.Invoice_ID, invoiceAsBytes)

	//Create Notification
	//notificationID , status , notificationType , orgId , createdBy , poNumber , invoiceNumber 
	s.createNotification(APIstub,invoice.Invoice_ID, "DISPATCH_NOTE_CREATED", "DISPATCH_NOTE", invoice.DestinationOrgId, invoice.OriginationOrgId, "-1",invoice.Invoice_Number)

	return shim.Success([]byte("Invoice Placed successfully"))
}

/*	
 *  Update the Shipping Status
 *  1. Status Update Accept/Reject done
 *	2. Quantity update Update to check the accepted and rejected Quantity done
 *	3. Product Quantity to check Org holding of particular product Quantity  done
 *  input :=  UpdatedBy , invoiceID , Status, Comments , quantity(accepted) 
 */

func (s *SmartContract) updateInvoiceStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//Parsing Request
	ob := args[0]
    invoice := Invoice{}
	err := json.Unmarshal([]byte(ob), &invoice)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	//Fetch Invoice object
	invoiceObAsBytes, _ := APIstub.GetState(invoice.Invoice_ID)
	var invoiceObject Invoice
	if len(invoiceObAsBytes) != 0 {
		err := json.Unmarshal(invoiceObAsBytes, &invoiceObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}

	if invoice.Status == "Accepted" ||  invoice.Status == "Rejected"{

		//Check whether the invoice is accepted/rejected by destinationorg 
		if(invoiceObject.DestinationOrgId  != invoice.UpdatedBy){
			return shim.Error("Only Destination Org can accept/reject the Dispatch Note You donot have rights to accept or reject")
		}


		// 2. Quantity update  to check the accepted and rejected Quantity
		//Fetch PO Ordered Quantity
		poQuantityKey := "PO"+invoiceObject.REF_POID
		quantityObAsBytes, _ := APIstub.GetState(poQuantityKey)
		var quantityOb Quantity
		if len(quantityObAsBytes) != 0 {
			err := json.Unmarshal(quantityObAsBytes, &quantityOb)
			if err != nil {
				return shim.Error("Unmashalling Error")
			}
		}

		if invoice.Status == "Accepted" {
			//To Make relation when recalling the Invoice 
			// So to we can deduct only accepted quantity
			invoiceObject.Quantity_Accepted = invoiceObject.Quantity_Accepted + invoice.Quantity;

			//Update Quantity Object
			quantityOb.Quantity_Accepted= quantityOb.Quantity_Accepted + invoice.Quantity
			quantityOb.Quantity_Balance = quantityOb.Quantity_Balance - quantityOb.Quantity_Accepted
			
			quantityObAsBytes, _ := json.Marshal(quantityOb)
			APIstub.PutState(poQuantityKey, quantityObAsBytes)

				//3. Product Quantity to check Org holding of particular product Quantity 
				//Fetch Product Details
				productkey := invoiceObject.Product_ID + "_" + invoice.UpdatedBy
				productAsBytes, _ := APIstub.GetState(productkey)
				var productObject Product
				if len(productAsBytes) != 0 {
					err := json.Unmarshal(productAsBytes, &productObject)
					if err != nil {
						return shim.Error("Unmashalling Error")
					}
					productObject.Quantity = productObject.Quantity + invoice.Quantity
					productObject.TxId = APIstub.GetTxID()

					productObAsBytes, _ := json.Marshal(productObject)
					APIstub.PutState(productkey, productObAsBytes)
				}else{
					product := Product{}
					product.Product_Key = productkey
					product.Product_ID = invoiceObject.Product_ID
					product.ProductDesc = invoiceObject.ProductDesc
					product.OrgId = invoice.UpdatedBy
					product.Quantity = invoice.Quantity
					product.DOC_TYPE = "Product";
					product.TxId = APIstub.GetTxID()
					productAsBytesNew, _ := json.Marshal(product)
					APIstub.PutState(product.Product_Key, productAsBytesNew)
				}	
		}else if invoice.Status == "Rejected"{
				//Fetch Product Details
					//4  Org Quantity Check to check Org Holdings Pending
					productkey := invoiceObject.Product_ID + "_" + invoiceObject.OriginationOrgId
					productAsBytes, _ := APIstub.GetState(productkey)
					var productObject Product
					if len(productAsBytes) != 0 {
						err := json.Unmarshal(productAsBytes, &productObject)
						if err != nil {
							return shim.Error("Unmashalling Error")
						}
					}
					//Update Product Stock once dispatched 
					productObject.Quantity = productObject.Quantity + invoice.Quantity
					productAsBytesNew, _ := json.Marshal(productObject)
					APIstub.PutState(productkey, productAsBytesNew)
	
		}

	}

	if invoice.Status == "Shipped" {
		//Check whether the invoice is accepted/rejected by destinationorg 
		if(invoiceObject.Transporter_ID != invoice.UpdatedBy){
			return shim.Error("Only Transprter can mark Dispatch Note as Shipped. You do not have rights to mark shipped" + invoiceObject.Transporter_ID +" || " + invoice.Transporter_ID + "||" )
		}
	}	

	//1. Status Update Accept/Reject
	invoiceObject.Status = invoice.Status
	invoiceObject.Comments = invoice.Comments
	invoiceObject.UpdatedBy = invoice.UpdatedBy
	invoiceObject.UPDATED_DATE = invoice.UPDATED_DATE
	invoiceAsBytes, _ := json.Marshal(invoiceObject)
	APIstub.PutState(invoice.Invoice_ID, invoiceAsBytes)

	//Create Notification
	//notificationID , status , notificationType , orgId , createdBy , poNumber , invoiceNumber 
	s.createNotification(APIstub,invoice.Invoice_ID, "DISPATCH_NOTE_"+invoice.Status, "DISPATCH_NOTE",  invoiceObject.OriginationOrgId,invoice.UpdatedBy, "-1",invoiceObject.Invoice_Number)


	return shim.Success([]byte("Invoice Updated successfully"))
}

func (s *SmartContract) updateInvoiceComments(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
		//Parsing Request
		ob := args[0]
		invoice := Invoice{}
		err := json.Unmarshal([]byte(ob), &invoice)
		if err != nil {
			fmt.Println(err.Error()) 
			return shim.Error("Unmashalling Error")
		}
	
		//Fetch Invoice object
		invoiceObAsBytes, _ := APIstub.GetState(invoice.Invoice_ID)
		var invoiceObject Invoice
		if len(invoiceObAsBytes) != 0 {
			err := json.Unmarshal(invoiceObAsBytes, &invoiceObject)
			if err != nil {
				return shim.Error("Unmashalling Error")
			}
		}

    //1. Status Update Accept/Reject
	invoiceObject.Comments = invoice.Comments
	invoiceObject.Status_Update = "false"
	invoiceObject.UpdatedBy = invoice.UpdatedBy
	invoiceAsBytes, _ := json.Marshal(invoiceObject)
	APIstub.PutState(invoice.Invoice_ID, invoiceAsBytes)

	return shim.Success([]byte("Invoice Comments Updated successfully"))

}

func (s *SmartContract) queryInvoiceByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	invoiceID := args[0]
	invoiceAsBytes, _ := APIstub.GetState(invoiceID)
	return shim.Success(invoiceAsBytes)
}

func (s *SmartContract) queryInvoiceByInvoiceNumber(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	invoiceNumber := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"INVOICE\",\"invoiceNumber\":\"%s\"}}", invoiceNumber)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


func (s *SmartContract) queryInvoiceByPONumber(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	poNumber := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"INVOICE\",\"poNumber\":\"%s\"}}", poNumber)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


func (s *SmartContract) queryInvoiceByPOID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	poID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"INVOICE\",\"refPOID\":\"%s\"}}", poID)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryInvoiceByStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	status := args[0]
	orgID := args[1]
	queryString := "{\"selector\":{\"docType\":\"INVOICE\",\"status\":\""+status+"\",\"$or\": [{\"originationOrgId\":\"" + orgID +"\"},{\"destinationOrgId\":\""+orgID+"\"},{\"transporterID\":\""+orgID+"\"}]}}"
	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"INVOICE\",\"status\":\"%s\"}}", status)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryInvoiceByOrgID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgID := args[0]
	queryString := "{\"selector\":{\"docType\":\"INVOICE\",\"$or\": [{\"originationOrgId\":\"" + orgID +"\"},{\"destinationOrgId\":\""+orgID+"\"},{\"transporterID\":\""+orgID+"\"}]}}"
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

/*
 *  Recall Invoice 
 *  1] Product Quantity Source added by Recall quantity in organization holding
 *	2] Product Quantity Accepted by Destination will be deducted from orgnization Holdings
 * 	3] Update Qunatity Object 
 */


func (s *SmartContract) recallInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//Parsing Request
	ob := args[0]
    invoice := Invoice{}
	err := json.Unmarshal([]byte(ob), &invoice)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	//Fetch Invoice object
	invoiceObAsBytes, _ := APIstub.GetState(invoice.Invoice_ID)
	var invoiceObject Invoice
	if len(invoiceObAsBytes) != 0 {
		err := json.Unmarshal(invoiceObAsBytes, &invoiceObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}

		//Check whether the invoice is accepted/rejected by destinationorg 
		if(invoiceObject.OriginationOrgId  != invoice.UpdatedBy){
			return shim.Error("Only Origination Org can recall the Dispatch Note You donot have rights recall")
		}

		// 2. Quantity update  to check the accepted and rejected Quantity
		//Fetch PO Ordered Quantity
		poQuantityKey := "PO"+invoiceObject.REF_POID
		quantityObAsBytes, _ := APIstub.GetState(poQuantityKey)
		var quantityOb Quantity
		if len(quantityObAsBytes) != 0 {
			err := json.Unmarshal(quantityObAsBytes, &quantityOb)
			if err != nil {
				return shim.Error("Unmashalling Error")
			}
		}

			// Change Quantity Object
			quantityOb.Quantity_Sent= quantityOb.Quantity_Sent  - invoiceObject.Quantity 
			if invoiceObject.Quantity_Accepted > 0 {
				quantityOb.Quantity_Accepted= quantityOb.Quantity_Accepted - invoiceObject.Quantity_Accepted
				quantityOb.Quantity_Balance= quantityOb.Quantity_Balance  + invoiceObject.Quantity_Accepted 
			}
			quantityObAsBytesNew, _ := json.Marshal(quantityOb)
			APIstub.PutState(poQuantityKey, quantityObAsBytesNew)


			//3. Product Quantity to update Org holding of particular product Quantity Source (invoice quantity added in product quantity)
				productkey := invoiceObject.Product_ID + "_" + invoice.UpdatedBy
				productAsBytes, _ := APIstub.GetState(productkey)
				var productObject Product
				if len(productAsBytes) != 0 {
					err := json.Unmarshal(productAsBytes, &productObject)
					if err != nil {
						return shim.Error("Unmashalling Error")
					}
					productObject.Quantity = productObject.Quantity + invoiceObject.Quantity
					productObject.TxId = APIstub.GetTxID()

					productObAsBytes, _ := json.Marshal(productObject)
					APIstub.PutState(productkey, productObAsBytes)
				}


				//3. Product Quantity to update Org holding of particular product Quantity Destination (invoice accepted quantity removed from destination holding quantity)
				productkey1 := invoiceObject.Product_ID + "_" + invoiceObject.DestinationOrgId
				productAsBytes1, _ := APIstub.GetState(productkey1)
				var productObject1 Product
				if len(productAsBytes1) != 0 {
					err := json.Unmarshal(productAsBytes1, &productObject1)
					if err != nil {
						return shim.Error("Unmashalling Error")
					}
					productObject1.Quantity = productObject1.Quantity - invoiceObject.Quantity_Accepted
					productObject1.TxId = APIstub.GetTxID()

					productObAsBytes1, _ := json.Marshal(productObject1)
					APIstub.PutState(productkey1, productObAsBytes1)
				}

	//1. Status Update Accept/Reject
	invoiceObject.Status = invoice.Status
	invoiceObject.Comments = invoice.Comments
	invoiceObject.UpdatedBy = invoice.UpdatedBy
	invoiceObject.UPDATED_DATE = invoice.UPDATED_DATE
	invoiceAsBytes, _ := json.Marshal(invoiceObject)
	APIstub.PutState(invoice.Invoice_ID, invoiceAsBytes)

	//Create Notification
	//notificationID , status , notificationType , orgId , createdBy , poNumber , invoiceNumber 
	s.createNotification(APIstub,"Recalled_"+invoice.Invoice_ID, "DISPATCH_NOTE_"+invoice.Status, "DISPATCH_NOTE",  invoiceObject.OriginationOrgId,invoice.UpdatedBy, "-1",invoiceObject.Invoice_Number)

	return shim.Success([]byte("Invoice Recalled successfully"))
}


func (s *SmartContract) createPO(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    po := PO{}
	err := json.Unmarshal([]byte(ob), &po)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	po.DOC_TYPE = "PO";
	po.TxId = APIstub.GetTxID()

	poAsBytes, _ := json.Marshal(po)
	APIstub.PutState(po.POID, poAsBytes)

	//Create Notification
	//notificationID , status , notificationType , orgId , createdBy , poNumber , invoiceNumber 
	s.createNotification(APIstub,po.POID, "PO_CREATION", "PO", po.DestinationOrgId, po.OriginationOrgId, po.PONumber,"-1")

	return shim.Success([]byte("PO Placed successfully"))
}

func (s *SmartContract) updatePOStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    po := PO{}
	err := json.Unmarshal([]byte(ob), &po)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	poListAsBytes, _ := APIstub.GetState(po.POID)
	var poObject PO
	if len(poListAsBytes) != 0 {
		err := json.Unmarshal(poListAsBytes, &poObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}


	if(poObject.DestinationOrgId  != po.CreatedBy){
		return shim.Error("Only Destination Org can accept/reject the Purchase Order You donot have rights to accept or reject")
	}

		//Make Entry for ordered Quantity
		quantityOb := Quantity{}
		quantityOb.Quantity_Ordered = poObject.Quantity
		quantityOb.Quantity_Sent    = 0
		quantityOb.Quantity_Accepted= 0
		quantityOb.Quantity_Balance = poObject.Quantity
		quantityOb.POID = "PO"+po.POID
		quantityObAsBytes, _ := json.Marshal(quantityOb)
		APIstub.PutState(quantityOb.POID, quantityObAsBytes)

	poObject.Status = po.Status
	poObject.Comments = po.Comments
	poObject.UPDATED_DATE = po.UPDATED_DATE

	poAsBytesNew, _ := json.Marshal(poObject)
	APIstub.PutState(poObject.POID, poAsBytesNew)

	//Create Notification
	//notificationID , status , notificationType , orgId , createdBy , poNumber , invoiceNumber 
	s.createNotification(APIstub,po.POID, po.Status, "PO", poObject.OriginationOrgId, po.CreatedBy, poObject.PONumber,"-1")

	return shim.Success([]byte("PO Updated successfully"))
}

func (s *SmartContract) updatePOComments(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    po := PO{}
	err := json.Unmarshal([]byte(ob), &po)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	poListAsBytes, _ := APIstub.GetState(po.POID)
	var poObject PO
	if len(poListAsBytes) != 0 {
		err := json.Unmarshal(poListAsBytes, &poObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}

	if((poObject.OriginationOrgId != po.CreatedBy) && (poObject.DestinationOrgId != po.CreatedBy)){
		return shim.Error("Only Origin/Destination Orgs can update the Purchase Order comments. You do not have rights to accept or reject. " + poObject.OriginationOrgId + " || " + poObject.DestinationOrgId + " || " + po.CreatedBy)
	}

	poObject.Comments = po.Comments
	poObject.Status_Update = "false"
	poObject.CreatedBy =  po.CreatedBy

	poAsBytesNew, _ := json.Marshal(poObject)
	APIstub.PutState(poObject.POID, poAsBytesNew)

	return shim.Success([]byte("PO Updated successfully"))
}


func (s *SmartContract) rejectPO(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    po := PO{}
	err := json.Unmarshal([]byte(ob), &po)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	poListAsBytes, _ := APIstub.GetState(po.POID)
	var poObject PO
	if len(poListAsBytes) != 0 {
		err := json.Unmarshal(poListAsBytes, &poObject)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}


	if(poObject.DestinationOrgId  != po.CreatedBy){
		return shim.Error("Only Destination Org can accept/reject the Purchase Order You donot have rights to accept or reject")
	}

		//Make Entry for ordered Quantity
		quantityOb := Quantity{}
		quantityOb.Quantity_Ordered = poObject.Quantity
		quantityOb.Quantity_Sent    = 0
		quantityOb.Quantity_Accepted= 0
		quantityOb.Quantity_Balance = poObject.Quantity
		quantityOb.POID = "PO"+po.POID
		quantityObAsBytes, _ := json.Marshal(quantityOb)
		APIstub.PutState(quantityOb.POID, quantityObAsBytes)

	poObject.Status = po.Status
	poObject.Comments = po.Comments
	poObject.UPDATED_DATE = po.UPDATED_DATE

	poAsBytesNew, _ := json.Marshal(poObject)
	APIstub.PutState(poObject.POID, poAsBytesNew)

	//Create Notification
	//notificationID , status , notificationType , orgId , createdBy , poNumber , invoiceNumber 
	s.createNotification(APIstub,po.POID, po.Status, "PO", poObject.OriginationOrgId, po.CreatedBy, poObject.PONumber,"-1")

	return shim.Success([]byte("PO Updated successfully"))
}

func (s *SmartContract) queryPOByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	poID := args[0]
	poAsBytes, _ := APIstub.GetState(poID)
	return shim.Success(poAsBytes)
}

func (s *SmartContract) queryPOByPONumber(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	poNumber := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"PO\",\"poNumber\":\"%s\"}}", poNumber)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryPOByRefPOID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	poID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"PO\",\"refPOID\":\"%s\"}}", poID)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryPOsByKey(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	owner := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"PO\",\"owner\":\"%s\"}}", owner)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryPOsByStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	status := args[0]
	orgID := args[1]
	queryString := "{\"selector\":{\"docType\":\"PO\",\"status\":\""+status+"\",\"$or\": [{\"originationOrgId\":\"" + orgID +"\"},{\"destinationOrgId\":\""+orgID+"\"}]}}"
	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"PO\",\"status\":\"%s\"}}", status)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryPOsByProductID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"PO\",\"productID\":\"%s\"}}", productID)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


func (s *SmartContract) queryPOByOrgID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgID := args[0]
	queryString := "{\"selector\":{\"docType\":\"PO\",\"$or\": [{\"originationOrgId\":\"" + orgID +"\"},{\"destinationOrgId\":\""+orgID+"\"}]}}"
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryAllPOs(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"PO\"}}")
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


// ===========================================================================================
// getHistoryForRecord returns the histotical state transitions for a given key of a record
// ===========================================================================================
func (t *SmartContract) queryPOHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recordKey := args[0]

	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON vehiclePart)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}



// ===========================================================================================
// getHistoryForRecord returns the histotical state transitions for a given key of a record
// ===========================================================================================
func (t *SmartContract) queryPOHistoryWithoutMetadata(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recordKey := args[0]

	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(response.Value))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

