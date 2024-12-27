package countries

import (
	"fmt"

	"github.com/dennj/govly/cmd/lib"
)

// generateITXML generates an XML string for an Italian electronic invoice (Fattura Elettronica)
// based on the provided VATRequest data.
//
// Parameters:
//   - vatRequest: A lib.VATRequest object containing the necessary data for the invoice.
//
// Returns:
//   - A string containing the formatted XML for the Italian electronic invoice.
//
// The generated XML includes the following sections:
//   - FatturaElettronicaHeader: Contains information about the seller (CedentePrestatore) and the buyer (CessionarioCommittente).
//   - FatturaElettronicaBody: Contains general data about the invoice, such as document type, currency, date, invoice number, and total amount.
//
// Note:
//   - The buyer's country code is hardcoded as "IT".
//   - The buyer's VAT number is assumed to be the same as the seller's VAT number for simplicity.
//   - The invoice number is hardcoded as "12345" for example purposes.
func GenerateITXML(vatRequest lib.VATRequest) string {
	return fmt.Sprintf(`
		<?xml version="1.0" encoding="UTF-8"?>
		<FatturaElettronica versione="1.2" xmlns="http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2">
			<FatturaElettronicaHeader>
				<CedentePrestatore>
					<IdFiscaleIVA>
						<IdPaese>IT</IdPaese>
						<IdCodice>%s</IdCodice>
					</IdFiscaleIVA>
					<Denominazione>%s</Denominazione>
				</CedentePrestatore>
				<CessionarioCommittente>
					<IdFiscaleIVA>
						<IdPaese>%s</IdPaese>
						<IdCodice>%s</IdCodice>
					</IdFiscaleIVA>
					<Denominazione>%s</Denominazione>
				</CessionarioCommittente>
			</FatturaElettronicaHeader>
			<FatturaElettronicaBody>
				<DatiGenerali>
					<DatiGeneraliDocumento>
						<TipoDocumento>TD01</TipoDocumento>
						<Divisa>EUR</Divisa>
						<Data>%s</Data>
						<Numero>%s</Numero>
						<ImportoTotaleDocumento>%s</ImportoTotaleDocumento>
					</DatiGeneraliDocumento>
				</DatiGenerali>
			</FatturaElettronicaBody>
		</FatturaElettronica>
	`,
		vatRequest.RegNum,    // Seller's VAT number
		vatRequest.Name,      // Seller's name
		"IT",                 // Buyer country code
		vatRequest.RegNum,    // Buyer's VAT number (assumed same for simplicity)
		vatRequest.Name,      // Buyer's name
		vatRequest.StartDate, // Invoice date
		"12345",              // Invoice number (example)
		vatRequest.Sales,     // Total invoice amount
	)
}
