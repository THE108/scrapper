package main

import (
	"flag"
	"log"

	"github.com/THE108/scrapper/config"
	"github.com/THE108/scrapper/info"
	"github.com/THE108/scrapper/parser"
	"github.com/THE108/scrapper/pdp"
	"github.com/THE108/scrapper/result"
)

func diffDeliveryInfo(expected, obtained map[string]info.DeliveryDetails) {
	log.Printf("expected: %v obtained: %v", expected, obtained)

	for deliveryTypeExpected, feeExpected := range expected {
		feeObtained, found := obtained[deliveryTypeExpected]
		if !found {
			log.Printf("delivery type %q not found", deliveryTypeExpected)
			continue
		}

		if feeExpected.Fee != feeObtained.Fee {
			log.Printf("shipping fee is not equal. delivery type: %s expected: %.2f obtained: %.2f",
				deliveryTypeExpected, feeExpected.Fee, feeObtained.Fee)
		}

		if feeExpected.Promo != "" && feeExpected.Promo != feeObtained.Promo {
			log.Printf("promo message is not equal. delivery type: %s expected: %s obtained: %s",
				deliveryTypeExpected, feeExpected.Promo, feeObtained.Promo)
		}
	}
}

func diff(configItemList []config.Item, resultMap map[string]info.PageInfo) {
	for _, item := range configItemList {
		log.Printf("[%s] ----->", item.Url)

		pageInfoObtained, found := resultMap[item.Url]
		if !found {
			log.Printf("can't find a PDP key: %s\n<-----\n\n", item.Url)
			continue
		}

		diffDeliveryInfo(item.DeliveryOptions, pageInfoObtained.DeliveryInfo)

		log.Printf("<-----\n\n")
	}
}

func validateUrlList(urlListFilename string, allowedDomains []string) {
	configItemList, err := config.ParseConfig(urlListFilename)
	if err != nil {
		log.Fatalf("error parse url list file: %s", err.Error())
	}

	urlList := make([]string, 0, len(configItemList))
	for _, item := range configItemList {
		urlList = append(urlList, item.Url)
	}

	resultHolder := result.NewResultHolder()
	parser.Parse(urlList, allowedDomains, pdp.NewParser(resultHolder))

	diff(configItemList, resultHolder.Get())
}

func parseSingleUrl(url string, allowedDomains []string) {
	resultHolder := result.NewResultHolder()
	parser.Parse([]string{url}, allowedDomains, pdp.NewParser(resultHolder))

	resultMap := resultHolder.Get()
	resultItem, found := resultMap[url]
	if !found {
		log.Printf("can't find %s in resultItem map", url)
		return
	}

	log.Printf("[%s] %+v", url, resultItem)
}

func main() {
	var url, urlListFilename string
	flag.StringVar(&url, "url", "", "Url to parse. Use to get shipping info from PDP page.")
	flag.StringVar(&urlListFilename, "file", "", "Config file name (yaml). Use to validate shipping info according to config.")
	flag.Parse()

	allowedDomains := []string{"www.lazada.sg", "pdp.lazada.sg", "lazada.sg"}

	if urlListFilename != "" {
		validateUrlList(urlListFilename, allowedDomains)
		return
	} else if url != "" {
		parseSingleUrl(url, allowedDomains)
		return
	}

	log.Fatalln("error parsing arguments: either url or file must be provided")
}
