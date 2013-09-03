package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	defer func() {
		str := recover()
		if str != nil {
			fmt.Println(str)
		}
	}()

	//godir()
	//gowalk()
	//readTerminalLog()

	fmt.Println("emv decoder V0.03 by Philsong@techtrex.com", "\nPls send suggestion to me, thanks")
	fmt.Println("\n------------------------\nemv TAG support list in below:")

	type Prompt struct {
		tip1st     string
		tip2nd     string
		tlvdatalen int
		bits       [][8]string
	}
	emvdecoder := [10]Prompt{
		{("1:AUC(support bits analysis)"), "Please enter a AUC value: ", 2, make([][8]string, 2)},
		{("2:TVR(support bits analysis)"), "Please enter a TVR value: ", 5, make([][8]string, 5)},
		{("3:TSI"), "Please enter a TSI value: ", 2, make([][8]string, 2)},
		{("4:CVR"), "Please enter a CVR value: ", 2, make([][8]string, 2)},
		{("5:AIP"), "Please enter a AIP value: ", 2, make([][8]string, 2)},

		{("6:TC"), "Please enter a TC value: ", 3, make([][8]string, 3)},
		{("7:ATC"), "Please enter a ATC value: ", 5, make([][8]string, 5)},
		{("8:CVM"), "Please enter a CVM value: ", 3, make([][8]string, 3)},
		{("9:Issuer Script Results"), "Please enter a Issuer Script Results value: ", 5, make([][8]string, 5)},
		{("10:Authorisation Response Code"), "Please enter a Authorisation Response Code value: ", 2, make([][8]string, 2)},
	}

	for _, value := range emvdecoder {
		fmt.Println(value.tip1st)
	}

	fmt.Print("Please enter a valid index in above:")
	var input string
	fmt.Scanln(&input)

	trimmed := strings.TrimSpace(input)
	item, _ := strconv.Atoi(trimmed)

	if item > 0 && item <= len(emvdecoder) {
		fmt.Println(emvdecoder[item-1].tip2nd)
	} else {
		fmt.Println("Please enter a valid index!")
		return
	}

	fmt.Scanln(&input)
	tlvdata := strings.TrimSpace(input)
	tlvdata = strings.Replace(tlvdata, " ", "", -1)
	//fmt.Printf("tlvdata str value[%s]\n", tlvdata)

	if len(tlvdata) != emvdecoder[item-1].tlvdatalen*2 {
		fmt.Printf("wrong data, must be %d bytes\n\n", emvdecoder[item-1].tlvdatalen)
		return
	}

	tlvbytes, _ := hex.DecodeString(tlvdata)
	fmt.Println("tlvbytes mem value", tlvbytes)
	fmt.Printf("tlvbytes hex value 0x%08x\n", tlvbytes)

	fmt.Println("-------------------------我是分割线--------------------------")

	switch item {
	case 1:
		emvdecoder[item-1].bits = initAUC()
	case 2:
		emvdecoder[item-1].bits = initTVR()
	case 3:
		emvdecoder[item-1].bits = initTSI()
	case 4:
		emvdecoder[item-1].bits = initCVR()
	case 5:
		emvdecoder[item-1].bits = initAIP()
	case 6:
		emvdecoder[item-1].bits = initTC()
	case 7:
		emvdecoder[item-1].bits = initATC()
	case 8:
		emvdecoder[item-1].bits = initCVR()
	case 9:
		emvdecoder[item-1].bits = initCVR()
	case 10:
		emvdecoder[item-1].bits = initCVR()
	default:
		return
	}

	for i, v := range tlvbytes {
		//fmt.Printf("BYTE[%d] base16 is 0x%02x\n", i, v)
		fmt.Printf("BYTE[%d] %08b\n", i+1, v)

		printElement(emvdecoder[item-1].bits[i], v)
	}

	fmt.Println("-------------------------the end---------------------------")
}

func printElement(tvr_elements [8]string, v uint8) {
	for j := 0; j < 8; j++ {
		var shiftNum uint32 = uint32(7 - j)
		//fmt.Printf("shift[%d]\n", shiftNum)
		var mask uint8 = 0x01 << shiftNum
		//fmt.Println("mask", mask)
		fmt.Print(j, ":  ")
		if v&mask == mask {
			fmt.Println(tvr_elements[j])
		} else {
			fmt.Println("--------------")
		}
	}
}

func readTerminalLog() {
	file, err := os.Open("test 01.txt")
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()
	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return
	}
	str := string(bs)
	writeTerminalLog(str)
	fmt.Println(str)
}

func writeTerminalLog(str string) {
	file, err := os.Create("test.txt")
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()
	file.WriteString(str)
}

func initAUC() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "Valid for domestic cash transactions "
	elements[index][1] = "Valid for international cash transactions "
	elements[index][2] = "Valid for domestic goods "
	elements[index][3] = "Valid for international goods "
	elements[index][4] = "Valid for domestic services "
	elements[index][5] = "Valid for international services "
	elements[index][6] = "Valid at ATMs "
	elements[index][7] = "Valid at terminals other than ATMs "

	index = 1
	elements[index][0] = "Domestic cashback allowed "
	elements[index][1] = "International cashback allowed "
	elements[index][2] = "Application not yet effective"
	elements[index][3] = "Requested service not allowed for card product "
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

//TVR data
func initTVR() [][8]string {
	var tvr_elements [][8]string = make([][8]string, 5)

	index := 0
	tvr_elements[index][0] = "Offline data authentication was not performed"
	tvr_elements[index][1] = "SDA failed"
	tvr_elements[index][2] = "ICC data missing"
	tvr_elements[index][3] = "Card appears on terminal exception file"
	tvr_elements[index][4] = "DDA failed"
	tvr_elements[index][5] = "CDA failed"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 1
	tvr_elements[index][0] = "ICC and terminal have different application versions"
	tvr_elements[index][1] = "Expired application "
	tvr_elements[index][2] = "Application not yet effective"
	tvr_elements[index][3] = "Requested service not allowed for card product "
	tvr_elements[index][4] = "New card "
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 2
	tvr_elements[index][0] = "Cardholder verification was not "
	tvr_elements[index][1] = "Unrecognised CVM "
	tvr_elements[index][2] = "PIN Try Limit exceeded "
	tvr_elements[index][3] = "PIN entry required and PIN pad "
	tvr_elements[index][4] = "not present or not working "
	tvr_elements[index][5] = "Online PIN entered "
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 3
	tvr_elements[index][0] = "Transaction exceeds floor limit "
	tvr_elements[index][1] = "Lower consecutive offline limit"
	tvr_elements[index][2] = "Upper consecutive offline limit "
	tvr_elements[index][3] = "Transaction selected randomly for online processing "
	tvr_elements[index][4] = "Merchant forced transaction online"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 4
	tvr_elements[index][0] = "Default TDOL used "
	tvr_elements[index][1] = "Issuer authentication failed "
	tvr_elements[index][2] = "Script processing failed before final GENERATE AC"
	tvr_elements[index][3] = "Script processing failed after final GENERATE AC"
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	return tvr_elements
}

func initTSI() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "Offline data authentication was performed  "
	elements[index][1] = "Cardholder verification was performed "
	elements[index][2] = "Card risk management was performed  "
	elements[index][3] = "Issuer authentication was performed "
	elements[index][4] = "Terminal risk management was performed "
	elements[index][5] = "Script processing was performed "
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	index = 1
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

func initCVR() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	index = 1
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

/*
Fail cardholder verification if this CVM is unsuccessful
Apply succeeding CV Rule if this CVM is unsuccessful
Fail CVM processing
Plaintext PIN verification performed by ICC
Enciphered PIN verified online
Plaintext PIN verification
performed by ICC and signature (paper)
Enciphered PIN verification performed by ICC
Enciphered PIN verification performed by ICC and signature (paper)
Values in the range 000110-011101 reserved for future use by this specification
Signature (paper)
No CVM required
Values in the range 100000-101111 reserved for use by the individual payment systems
Values in the range 110000-111110 reserved for use by the issuer
This value is not available for use
*/
func initAIP() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "RFU"
	elements[index][1] = "SDA supported "
	elements[index][2] = "DDA supported"
	elements[index][3] = "Cardholder verification is supported "
	elements[index][4] = "Terminal risk management is to be performed "
	elements[index][5] = "Issuer authentication is supported"
	elements[index][6] = "RFU"
	elements[index][7] = "CDA supported "

	index = 1
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

func initTC() [][8]string {
	var tvr_elements [][8]string = make([][8]string, 3)

	index := 0
	tvr_elements[index][0] = "Manual key entry "
	tvr_elements[index][1] = "Magnetic stripe "
	tvr_elements[index][2] = "IC with contacts"
	tvr_elements[index][3] = "RFU"
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 1
	tvr_elements[index][0] = "Plaintext PIN for ICC verification"
	tvr_elements[index][1] = "Enciphered PIN for online verification   "
	tvr_elements[index][2] = "Signature (paper)"
	tvr_elements[index][3] = "Enciphered PIN for offline verification "
	tvr_elements[index][4] = "No CVM Required  "
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 2
	tvr_elements[index][0] = "SDA"
	tvr_elements[index][1] = "DDA"
	tvr_elements[index][2] = "Card capture  "
	tvr_elements[index][3] = "RFU"
	tvr_elements[index][4] = "CDA"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	return tvr_elements
}

func initATC() [][8]string {
	var tvr_elements [][8]string = make([][8]string, 5)

	index := 0
	tvr_elements[index][0] = "Cash"
	tvr_elements[index][1] = "Goods"
	tvr_elements[index][2] = "Services"
	tvr_elements[index][3] = "Cashback"
	tvr_elements[index][4] = "Inquiry"
	tvr_elements[index][5] = "Transfer"
	tvr_elements[index][6] = "Payment"
	tvr_elements[index][7] = "Administrative"

	index = 1
	tvr_elements[index][0] = "Cash Deposit"
	tvr_elements[index][1] = "RFU"
	tvr_elements[index][2] = "RFU"
	tvr_elements[index][3] = "RFU"
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 2
	tvr_elements[index][0] = "Numeric keys "
	tvr_elements[index][1] = "Alphabetic and special characters keys  "
	tvr_elements[index][2] = "Command keys  "
	tvr_elements[index][3] = "Function keys  "
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 3
	tvr_elements[index][0] = "Print, attendant "
	tvr_elements[index][1] = "Print, cardholder "
	tvr_elements[index][2] = "Display, attendant "
	tvr_elements[index][3] = "Display, cardholder "
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "Code table 10"
	tvr_elements[index][7] = "Code table 9"

	index = 4
	tvr_elements[index][0] = "Code table 8"
	tvr_elements[index][1] = "Code table 7"
	tvr_elements[index][2] = "Code table 6"
	tvr_elements[index][3] = "Code table 5"
	tvr_elements[index][4] = "Code table 4"
	tvr_elements[index][5] = "Code table 3"
	tvr_elements[index][6] = "Code table 2"
	tvr_elements[index][7] = "Code table 1"

	return tvr_elements
}