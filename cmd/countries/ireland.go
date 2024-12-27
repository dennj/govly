package countries

import (
	"fmt"

	"github.com/dennj/govly/cmd/lib"
)

// generateIRXML generates XML payload for Revenue Ireland
func GenerateIRXML(vatRequest lib.VATRequest) string {
	return fmt.Sprintf(`
		<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
		<VAT3 
			currency="E" 
			name="%s" 
			regnum="%s" 
			startdate="%s" 
			enddate="%s" 
			sales="%s" 
			purchs="%s" 
			goodsto="%s" 
			goodsfrom="%s" 
			servicesto="%s" 
			servicesfrom="%s" 
			type="0" 
			filefreq="0" 
			formversion="1" 
			language="E" 
			postponedAccounting="%s" 
			unusualExpenditure="%s" 
			unusualExpenditureAmt="%s" 
			unusualExpenditureDtl="%s"/>
	`,
		vatRequest.Name,
		vatRequest.RegNum,
		vatRequest.StartDate,
		vatRequest.EndDate,
		vatRequest.Sales,
		vatRequest.Purchases,
		vatRequest.GoodsTo,
		vatRequest.GoodsFrom,
		vatRequest.ServicesTo,
		vatRequest.ServicesFrom,
		vatRequest.PostponedAccounting,
		vatRequest.UnusualExpenditure,
		vatRequest.UnusualExpenditureAmount,
		vatRequest.UnusualExpenditureDetail,
	)
}
