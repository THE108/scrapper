package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/THE108/scrapper/config"
	"github.com/THE108/scrapper/info"
	"github.com/THE108/scrapper/parser"
	"github.com/THE108/scrapper/pdp"
	"github.com/THE108/scrapper/result"
)

func diffDeliveryInfo(expected, obtained map[string]info.DeliveryDetails) (messages []string) {
	for deliveryTypeExpected, feeExpected := range expected {
		feeObtained, found := obtained[deliveryTypeExpected]
		if !found {
			messages = append(messages, fmt.Sprintf("delivery type %q not found", deliveryTypeExpected))
			continue
		}

		if feeExpected.Fee != feeObtained.Fee {
			messages = append(messages, fmt.Sprintf("shipping fee is not equal. delivery type: %s expected: %.2f obtained: %.2f",
				deliveryTypeExpected, feeExpected.Fee, feeObtained.Fee))
		}

		if feeExpected.Promo != "" && feeExpected.Promo != feeObtained.Promo {
			messages = append(messages, fmt.Sprintf("promo message is not equal. delivery type: %s expected: %s obtained: %s",
				deliveryTypeExpected, feeExpected.Promo, feeObtained.Promo))
		}
	}

	return
}

func diff(configItemList []config.Item, resultMap map[string]info.PageInfo, verbose bool) {
	for _, item := range configItemList {
		pageInfoObtained, found := resultMap[item.Url]
		if !found {
			log.Printf("can't find a PDP key: %s\n<-----\n\n", item.Url)
			continue
		}

		messages := diffDeliveryInfo(item.DeliveryOptions, pageInfoObtained.DeliveryInfo)

		if len(messages) > 0 || verbose {
			log.Printf("[%s]----->\n", item.Url)

			if verbose {
				printExpectedAndObtained(item.DeliveryOptions, pageInfoObtained.DeliveryInfo)
			}

			for _, msg := range messages {
				log.Printf("ERROR: %s", msg)
			}
			log.Printf("<-----\n\n")
		}
	}
}

func printExpectedAndObtained(expected, obtained map[string]info.DeliveryDetails) {
	filteredObtained := make(map[string]info.DeliveryDetails, len(expected))
	for deliveryType := range expected {
		if details, found := obtained[deliveryType]; found {
			filteredObtained[deliveryType] = details
		}
	}

	log.Printf("expected: %+v obtained: %+v", expected, filteredObtained)
}

func validateUrlList(urlListFilename string, allowedDomains []string, verbose bool) {
	configItemList, err := config.ParseConfig(urlListFilename)
	if err != nil {
		log.Fatalf("error parse url list file: %s", err.Error())
	}

	urlList := make([]string, 0, len(configItemList))
	for _, item := range configItemList {
		urlList = append(urlList, item.Url)
	}

	resultHolder := result.NewResultHolder()
	parser.Parse(urlList, allowedDomains, pdp.NewParser(resultHolder, verbose))

	diff(configItemList, resultHolder.Get(), verbose)
}

func parseSingleUrl(url string, allowedDomains []string, verbose bool) {
	resultHolder := result.NewResultHolder()
	parser.Parse([]string{url}, allowedDomains, pdp.NewParser(resultHolder, verbose))

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
	var verbose bool
	flag.StringVar(&url, "url", "", "Url to parse. Use to get shipping info from PDP page.")
	flag.StringVar(&urlListFilename, "file", "", "Config file name (yaml). Use to validate shipping info according to config.")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode. Use to see more data in logs.")
	flag.Parse()

	allowedDomains := []string{"www.lazada.sg", "pdp.lazada.sg", "lazada.sg"}

	if urlListFilename != "" {
		validateUrlList(urlListFilename, allowedDomains, verbose)
		return
	} else if url != "" {
		parseSingleUrl(url, allowedDomains, verbose)
		return
	}

	log.Fatalln("error parsing arguments: either url or file must be provided")
}
