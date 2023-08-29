package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"bytes"
	_"strconv"
	_"unsafe"
	_"strings"
	_"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)


//==============================================================================================================================
//	 Patent Chain Structure Definitions
//==============================================================================================================================

type SmartContract struct {
}

type ORG struct{
	OrgId					 string     			`json:"orgId"`
	OrgName     			 string					`json:"orgName"`
	OrgAdminFirstName 		 string					`json:"orgAdminFirstName"`
	OrgAdminLastName  		 string					`json:"orgAdminLastName"`
	OrgAdminId				 string					`json:"orgAdminEmailId"`
	Phone					 string					`json:"phone"`
	Address					 string					`json:"address"`
	CountryOfInc			 string					`json:"countryOfInc"`
	StateOfInc		    	 string					`json:"stateOfInc"`
	ZipCode 		    	 string					`json:"zipCode"`
	BuisnessType 		   	 string					`json:"buisnessType"`
	OrganizationType 	   	 string					`json:"organizationType"`
	Role					 string					`json:"role"`
	Status					 string					`json:"status"`
	TxId					 string					`json:"tx_id"`
	DOC_TYPE				 string					`json:"docType"`
	Approved_By				 string     			`json:"approvedBy"`
}

type ORGS struct {
	APPROVED_ORGS 					[]ORG					`json:"approved_orgs"`
	OrgId					    	 string     			`json:"orgId"`
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
	fmt.Println(function) 
	fmt.Println(args) 

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createOrg" {
		return s.createOrg(APIstub, args)
	}else if function == "queryOrgByID" {
		return s.queryOrgByID(APIstub, args)
	}else if function == "queryOrgByOrganizationType" {
		return s.queryOrgByOrganizationType(APIstub, args)
	}else if function == "queryOrgByOrganizationName" {
		return s.queryOrgByOrganizationName(APIstub, args)
	}else if function == "queryAllOrganisations" {
		return s.queryAllOrganisations(APIstub, args)
	}else if function == "approve" {
		return s.approve(APIstub, args)
	}else if function == "queryApprovedOrgs" {
		return s.queryApprovedOrgs(APIstub, args)
	}else if function == "queryApprovedOrgsByRole" {
		return s.queryApprovedOrgsByRole(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) approve(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    org := ORG{}
	err := json.Unmarshal([]byte(ob), &org)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

		//Fetch Product Details
		approvedOrgId :=  org.OrgId
		approvedOrgIdAsBytes, _ := APIstub.GetState(approvedOrgId)
		var orgDetails ORG
		if len(approvedOrgIdAsBytes) != 0 {
			err := json.Unmarshal(approvedOrgIdAsBytes, &orgDetails)
			if err != nil {
				return shim.Error("Unmashalling Error")
			}
		}

	//Fetch Approved Orgs Array for an organization
	key := "APPROVED_ORGS_FOR_"+org.Approved_By
	approvedOrgsAsBytes, _ := APIstub.GetState(key)
	var orgs ORGS
	if len(approvedOrgsAsBytes) != 0 {
		err := json.Unmarshal(approvedOrgsAsBytes, &orgs)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
		orgs.APPROVED_ORGS = append(orgs.APPROVED_ORGS,orgDetails)
	}else{
		orgs.OrgId = org.Approved_By
		orgs.APPROVED_ORGS = append(orgs.APPROVED_ORGS,orgDetails)
	}
	
	approvedOrgsAsBytesNew, _ := json.Marshal(orgs)
	APIstub.PutState(key, approvedOrgsAsBytesNew)
	return shim.Success([]byte("Org Aprroved successfully"))
}

func (s *SmartContract) queryApprovedOrgs(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	OrgId := args[0]
	key := "APPROVED_ORGS_FOR_"+OrgId
	
	var orgs ORGS
	orgsAsBytes, _ := APIstub.GetState(key)
	if len(orgsAsBytes) != 0 {
		err := json.Unmarshal(orgsAsBytes, &orgs)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}
	approvedOrgsAsBytesNew,_ := json.Marshal(orgs.APPROVED_ORGS)
	return shim.Success(approvedOrgsAsBytesNew)
}

func (s *SmartContract) queryApprovedOrgsByRole(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    org := ORG{}
	err := json.Unmarshal([]byte(ob), &org)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	key := "APPROVED_ORGS_FOR_"+org.OrgId
	var orgs ORGS
	orgsAsBytes, _ := APIstub.GetState(key)
	if len(orgsAsBytes) != 0 {
		err := json.Unmarshal(orgsAsBytes, &orgs)
		if err != nil {
			return shim.Error("Unmashalling Error")
		}
	}

	var result ORGS
	for i,_ := range orgs.APPROVED_ORGS{
		if orgs.APPROVED_ORGS[i].OrganizationType == org.OrganizationType{
			result.APPROVED_ORGS = append(result.APPROVED_ORGS,orgs.APPROVED_ORGS[i])
		}
	}

	approvedOrgsAsBytesNew,_ := json.Marshal(result.APPROVED_ORGS)
	return shim.Success(approvedOrgsAsBytesNew)
}

func (s *SmartContract) createOrg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    org := ORG{}
	err := json.Unmarshal([]byte(ob), &org)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	org.DOC_TYPE = "Organization";
	org.TxId = APIstub.GetTxID()
	orgAsBytes, _ := json.Marshal(org)
	APIstub.PutState(org.OrgId, orgAsBytes)
	return shim.Success([]byte("Org Created successfully"))
}


func (s *SmartContract) queryOrgByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	OrgId := args[0]
	poAsBytes, _ := APIstub.GetState(OrgId)
	return shim.Success(poAsBytes)
}

func (s *SmartContract) queryOrgByOrganizationType(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	organizationType := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Organization\",\"organizationType\":\"%s\"}}", organizationType)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryOrgByOrganizationName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgName := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Organization\",\"orgName\":\"%s\"}}", orgName)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryAllOrganisations(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"Organization\"}}")
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



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

