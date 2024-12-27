package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dennj/govly/cmd/countries"
	"github.com/dennj/govly/cmd/lib"

	"github.com/gin-gonic/gin"
)

const (
	HMRCAPIURL        = "https://test-api.service.hmrc.gov.uk/organisations/vat/%s/returns"
	RevenueIrelandURL = "https://api.revenue.ie/vat3/declaration"
)

// VATRequestHandler handles incoming VAT requests
func VATRequestHandler(c *gin.Context) {
	var vatRequest lib.VATRequest
	if err := c.ShouldBindJSON(&vatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch vatRequest.Country {
	case "IR":
		IR_VATRequestHandler(c, vatRequest)
	case "UK":
		UK_VATRequestHandler(c, vatRequest)
	case "IT":
		IT_VATRequestHandler(c, vatRequest)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported country"})
	}
}

func IT_VATRequestHandler(c *gin.Context, vatRequest lib.VATRequest) {
	// Generate XML payload based on Italy's FatturaPA schema
	xmlData := countries.GenerateITXML(vatRequest)

	// Define the SOAP envelope for the SdI platform
	soapEnvelope := fmt.Sprintf(`
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Header/>
			<soapenv:Body>
				<submitInvoice xmlns="http://www.fatturapa.gov.it/sdi">
					<xmlPayload><![CDATA[%s]]></xmlPayload>
				</submitInvoice>
			</soapenv:Body>
		</soapenv:Envelope>
	`, xmlData)

	// SOAP endpoint for the Italian SdI platform
	const SDIEndpoint = "https://test.fatturapa.gov.it/sdi"

	// Send SOAP request to SdI
	response, err := lib.SendHTTPRequest("POST", SDIEndpoint, "text/xml", soapEnvelope, vatRequest.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to SdI", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sdi_status":   response.Status,
		"sdi_response": response.Body,
	})
}

// IR_VATRequestHandler processes VAT submissions to Revenue Ireland
func IR_VATRequestHandler(c *gin.Context, vatRequest lib.VATRequest) {
	xmlData := countries.GenerateIRXML(vatRequest)

	response, err := lib.SendHTTPRequest("POST", RevenueIrelandURL, "application/xml", xmlData, vatRequest.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to Revenue Ireland", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"revenue_status":   response.Status,
		"revenue_response": response.Body,
	})
}

// UK_VATRequestHandler processes VAT submissions to HMRC
func UK_VATRequestHandler(c *gin.Context, vatRequest lib.VATRequest) {
	jsonData, err := json.Marshal(vatRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode JSON"})
		return
	}

	vatRegNum := vatRequest.RegNum // Replace with your actual VAT registration logic
	url := fmt.Sprintf(HMRCAPIURL, vatRegNum)

	response, err := lib.SendHTTPRequest("POST", url, "application/json", string(jsonData), vatRequest.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to HMRC", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hmrc_status":   response.Status,
		"hmrc_response": response.Body,
	})
}

// @title Govly API
// Govly is a API wrapper got whoile world governament API. It offer a standard and consistent way to interact with different governament API.
func main() {
	router := gin.Default()
	router.POST("/vat", VATRequestHandler)

	port := "8080"
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
