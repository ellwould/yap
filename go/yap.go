/*
    YAP (Yet Another PBX)
    Copyright (C) 2025 Elliot Michael Keavney

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see https://github.com/ellwould/yap/blob/main/LICENSE
*/

package main

import (
	//Commented out until source code for database is wrote 
	//"database/sql"
	"fmt"
	"github.com/ellwould/csvcell"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
)

//Constant for directory path that contains the files yap-start.html and yap-end.html
const dirHTML string = "/etc/yap/html-css"

//Constant for fileStartHTML file
const fileStartHTML string = "yap-start.html"

//Constant for fileEndHTML file
const fileEndHTML string = "yap-end.html"

//Function to retrive HTTP email header
func emailHeaderHTTP(r *http.Request) (email string) {
	email = r.Header.Get("X-Email")
	return email
}

//Function for the header
func header(w http.ResponseWriter) {
	fmt.Fprintf(w, "<div class=\"header\">")
	fmt.Fprintf(w, "  <h1>")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"generalButton headerButton logoutButton\">Logout</a>")
	fmt.Fprintf(w, "    <a href=\"/\" class=\"generalButton headerButton homeButton\">Home</a>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap/blob/main/LICENSE\" class=\"generalButton headerButton LicenseButton\">License</a>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    ✱ YAP (Yet Another PBX) ✱")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

//Function for the footer
func footer(w http.ResponseWriter) {
	fmt.Fprintf(w, "<div class=\"footer\">")
	fmt.Fprintf(w, "  <h2>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap\" class=\"generalButton footerButton\">YAP Source Code</a>")
	fmt.Fprintf(w, "    <a href=\"https://ell.today\" class=\"generalButton footerButton\">Other Software</a>")
	fmt.Fprintf(w, "  </h2>")
	fmt.Fprintf(w, "</div>")
}

//Function for the home page (/) - shows user infomation
func userInfomation(w http.ResponseWriter) {

	//Table showing the user thier name, email, account type and when the account was created
	fmt.Fprintf(w, "<div>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"table\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>&nbsp Name &nbsp</th>")
	fmt.Fprintf(w, "          <th>&nbsp Email &nbsp</th>")
	fmt.Fprintf(w, "          <th>&nbsp Account Type &nbsp</th>")
	fmt.Fprintf(w, "          <th>&nbsp Account Created &nbsp</th>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "          <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "          <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "          <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "      </table>")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><button onclick=\"toggleAccountDetail() \"class=\"generalButton\">Show / Hide More Account Details</button></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	//Account detail tables
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "<div id=\"accountDetailDiv\" style=\"display:none\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp Group Name &nbsp <br> and ID</th>")
	fmt.Fprintf(w, "    <th>&nbsp PBX Name &nbsp <br> and ID</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "  </tr>")	
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp Group Site Address &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Group Site Email &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Group Site Phone Number &nbsp</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp Group Invoice Address &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Group Invoice Email &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Group Invoice Phone Number &nbsp</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp Account Type Description &nbsp</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <td>&nbsp &nbsp</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</div>")
	//Java Script embedded to toggle account detail div
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function toggleAccountDetail() {")
	fmt.Fprintf(w, "  var x = document.getElementById(\"accountDetailDiv\");")
	fmt.Fprintf(w, "  if (x.style.display === \"none\") {")
	fmt.Fprintf(w, "    x.style.display = \"table\";")
	fmt.Fprintf(w, "  } else {")
	fmt.Fprintf(w, "    x.style.display = \"none\";")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "}")
	fmt.Fprintf(w, "</script>")
}

func main() {

	//Get the values from inside the YAP configuration file
	err := godotenv.Load("/etc/yap/yap.env")
	if err != nil {
		panic("Error loading yap.env file for database details")
	}

	//Get the database connection details
	dbUsername := os.Getenv("dbUsername")
	dbPassword := os.Getenv("dbPassword")
	dbTransport := os.Getenv("dbTransport")
	dbAddress := os.Getenv("dbAddress")
	dbPort := os.Getenv("dbPort")
	dbTls := os.Getenv("dbTls")

	//Values allowed for dbTransport Variable
	var allowedTransportValue = []string{"tcp", "udp"}
	validDbTransport := slices.Contains(allowedTransportValue, dbTransport)

	validateDbAddress := validator.New()
	validateDbAddressErr := validateDbAddress.Var(dbAddress, "required,ip_addr")

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic("DATABASE PORT MUST BE A NUMBER IN /etc/yap/yap.env")
	}

	//Values allowed for dbTls Variable
	var allowedTlsValue = []string{"false", "true"}
	validDbTls := slices.Contains(allowedTlsValue, dbTls)

	//Catch if any errors were made in yap.env and feed back where to correct error
	if dbUsername == "" {
		panic("DATABASE USERNAME CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if dbPassword == "" {
		panic("DATABASE PASSOWRD CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if dbTransport == "" {
		panic("DATABASE TRANSPORT OPTION CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if validDbTransport == false {
		panic("DATABASE TRANSPORT OPTION MUST BE udp OR tcp IN /etc/yap/yap.env")
	} else if validateDbAddressErr != nil && dbAddress != "localhost" {
		panic("DATABASE ADDRESS MUST BE A VALID INTERENT PROTOCOL (IP) ADDRESS OR localhost IN /etc/yap/yap.env")
	} else if dbPortInt <= 0 || dbPortInt >= 65536 {
		panic("DATABASE PORT MUST BE IN THE NUMBER RANGE 1-65535 IN /etc/yap/yap.env")
	} else if dbTls == "" {
		panic("DATABASE TLS OPTION CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if validDbTls == false {
		panic("DATABASE TRANSPORT OPTION MUST BE false OR true IN /etc/yap/yap.env")
	}

	startHTML := csvcell.FileData(dirHTML, fileStartHTML)
	endHTML := csvcell.FileData(dirHTML, fileEndHTML)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w)
		// Code to call the emailHeaderHTTP function
		//fmt.Fprintf(w, "<h2><p>"+emailHeaderHTTP(r)+"</p></h2>")
		//fmt.Fprintf(w, "<br>")
		userInfomation(w)
		footer(w)
		fmt.Fprintf(w, endHTML)

	})

	yapPort := os.Getenv("yapPort")
	yapPortInt, err := strconv.Atoi(yapPort)

	if err != nil {
		panic("YAP PORT MUST BE A NUMBER IN /etc/yap/yap.env")
	}

	yapAddress := os.Getenv("yapAddress")
	validateYapAddress := validator.New()
	validateYapAddressErr := validateYapAddress.Var(yapAddress, "required,ip_addr")

	if yapPortInt <= 1023 || yapPortInt >= 49152 {
		panic("YAP LISTENING PORT MUST BE IN THE NUMBER RANGE 1024-49151 IN /etc/yap/yap.env")
	} else if validateYapAddressErr != nil && yapAddress != "localhost" {
		panic("YAP ADDRESS MUST BE A VALID INTERENT PROTOCOL (IP) ADDRESS OR localhost IN /etc/yap/yap.env")
	} else if dbAddress == yapAddress && dbPort == yapPort {
		panic("YAP ADDRESS & PORT NUMBER CANNOT BE THE SAME AS DATABASE ADDRESS & PORT NUMBER IN /etc/yap/yap.env")
	} else {
		socket := yapAddress + ":" + yapPort
		fmt.Println("YAP is running on port " + socket)
		//Start server on port specified above
		log.Fatal(http.ListenAndServe(socket, nil))
	}
}

// Contributor(s):
// Elliot Michael Keavney
