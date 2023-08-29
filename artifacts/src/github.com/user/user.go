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

type User struct{
	UserId					 string					`json:"userId"`
	UserName				 string					`json:"userName"`
	Password				 string					`json:"password"`
	OrgId					 string     			`json:"orgId"`
	OrgName     			 string					`json:"orgName"`
	OrganizationType 	   	 string					`json:"organizationType"`
	Email 					 string					`json:"email"`
	Phone					 string					`json:"phone"`
	Status					 string					`json:"status"`
	Role					 string					`json:"role"`
	TxId					 string					`json:"tx_id"`
	DOC_TYPE				 string					`json:"docType"`
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
	if function == "createUser" {
		return s.createUser(APIstub, args)
	}else if function == "updateUser" {
		return s.updateUser(APIstub, args)
	}else if function == "queryUserByID" {
		return s.queryUserByID(APIstub, args)
	}else if function == "queryUserByOrganizationID" {
		return s.queryUserByOrganizationID(APIstub, args)
	}else if function == "queryUserByOrganizationName" {
		return s.queryUserByOrganizationName(APIstub, args)
	}else if function == "queryUserByRole" {
		return s.queryUserByRole(APIstub, args)
	}else if function == "queryAllUsers" {
		return s.queryAllUsers(APIstub, args)
	}else if function == "changeUserPassword" {
		return s.changeUserPassword(APIstub, args)
	}else if function == "resetUserPassword" {
		return s.resetUserPassword(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}


func (s *SmartContract) createUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    user := User{}
	err := json.Unmarshal([]byte(ob), &user)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	user.DOC_TYPE = "User";
	user.TxId = APIstub.GetTxID()
	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(user.UserId, userAsBytes)
	return shim.Success([]byte("User Created successfully"))
}

func (s *SmartContract) updateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    user := User{}
	err := json.Unmarshal([]byte(ob), &user)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}

	userObAsBytes, _ := APIstub.GetState(user.UserId)
	if len(userObAsBytes) == 0 {
		return shim.Error("Failed to get user with this Id")
	}

	user.DOC_TYPE = "User";
	user.TxId = APIstub.GetTxID()
	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(user.UserId, userAsBytes)
	return shim.Success([]byte("User updated successfully"))
}


func (s *SmartContract) queryUserByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	UserId := args[0]
	poAsBytes, _ := APIstub.GetState(UserId)
	return shim.Success(poAsBytes)
}

func (s *SmartContract) queryUserByOrganizationID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgId := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"User\",\"orgId\":\"%s\"}}", orgId)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryUserByOrganizationName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	orgName := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"User\",\"orgName\":\"%s\"}}", orgName)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryUserByRole(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	role := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"User\",\"role\":\"%s\"}}", role)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryAllUsers(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"User\"}}")
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

func (s *SmartContract) changeUserPassword(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    user := User{}
	err := json.Unmarshal([]byte(ob), &user)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	user.DOC_TYPE = "User";
	user.TxId = APIstub.GetTxID()
	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(user.UserId, userAsBytes)
	return shim.Success([]byte("Password Changed successfully"))
}


func (s *SmartContract) resetUserPassword(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	ob := args[0]
    user := User{}
	err := json.Unmarshal([]byte(ob), &user)
	if err != nil {
		fmt.Println(err.Error()) 
		return shim.Error("Unmashalling Error")
	}
	user.DOC_TYPE = "User";
	user.TxId = APIstub.GetTxID()
	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(user.UserId, userAsBytes)
	return shim.Success([]byte("Password reset successfully"))
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

