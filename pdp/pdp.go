package pdp

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/THE108/scrapper/conv"
	"github.com/THE108/scrapper/info"

	"github.com/dop251/goja"
	"github.com/gocolly/colly"
)

const injectedJs = `var window = {};
	var require = {
		config: function(data) {}
	};

	function requirejs(arr, fn) {
		var result = {};
		function run(cfg) {
			result.skuId = cfg.data.root.fields.primaryKey.skuId;
			result.itemId = cfg.data.root.fields.primaryKey.itemId;
			result.sellerId = cfg.data.root.fields.primaryKey.sellerId;
			result.deliveryOptions = cfg.data.root.fields.deliveryOptions[result.skuId];
		}
		fn({run: run});
		return result;
	}`

func dumpJsonToFile(url string, jsonContent string) {
	filename := strings.TrimPrefix(url, "http://")
	filename = strings.Replace(filename, "/", "-", -1)
	if err := ioutil.WriteFile(filename, []byte(jsonContent), 0666); err != nil {
		log.Printf("[%s] error writing dump file err: %s", url, err.Error())
	}
}

func parseStringValue(data map[string]interface{}, key string) (strVal string, err error) {
	ifcVal, found := data[key]
	if !found {
		err = fmt.Errorf("key not found: %s", key)
		return
	}

	strVal, ok := ifcVal.(string)
	if !ok {
		err = fmt.Errorf("key not found: %s", key)
	}

	return
}

func parseIntValue(data map[string]interface{}, key string) (intVal int, err error) {
	var strVal string
	strVal, err = parseStringValue(data, key)
	if err != nil {
		return
	}

	return strconv.Atoi(strVal)
}

func parseDeliveryInfoList(v interface{}) (deliveryInfo map[string]float64, err error) {
	if v == nil {
		err = fmt.Errorf("parsed delivery info value is nil")
		return
	}

	deliveryInfoList, ok := v.([]interface{})
	if !ok {
		err = fmt.Errorf("fail assert type: %T to []interface{}", v)
		return
	}

	deliveryInfo = make(map[string]float64)
	for _, infDeliveryInfoVal := range deliveryInfoList {
		deliveryInfoMap, ok := infDeliveryInfoVal.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("fail assert type: %T to map[string]interface{}}", deliveryInfoMap)
			return
		}

		var dataType string
		dataType, err = parseStringValue(deliveryInfoMap, "dataType")
		if err != nil {
			return
		}

		switch dataType {
		case "delivery", "liveup":
		default:
			continue
		}

		var deliveryTypeRaw string
		deliveryTypeRaw, err = parseStringValue(deliveryInfoMap, "type")
		if err != nil {
			return
		}

		var feeRaw string
		feeRaw, err = parseStringValue(deliveryInfoMap, "fee")
		if err != nil {
			return
		}

		var deliveryType string
		deliveryType, err = conv.ToDeliveryType(deliveryTypeRaw)
		if err != nil {
			err = fmt.Errorf("error parse delivery type value value: %q error: %s", deliveryTypeRaw, err.Error())
			return
		}

		var fee float64
		fee, err = conv.ToMoney(feeRaw)
		if err != nil {
			err = fmt.Errorf("error parse shipping fee value delivery type: %q error: %s", deliveryType, err.Error())
			return
		}

		deliveryInfo[deliveryType] = fee
	}

	return
}

func parseInfo(v interface{}) (pageInfo info.PageInfo, err error) {
	if v == nil {
		err = fmt.Errorf("parsed page info value is nil")
		return
	}

	data, ok := v.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("fail assert type: %T to map[string]interface{}", v)
		return
	}

	pageInfo.ItemID, err = parseIntValue(data, "itemId")
	if err != nil {
		return
	}

	pageInfo.SkuID, err = parseIntValue(data, "skuId")
	if err != nil {
		return
	}

	pageInfo.SellerID, err = parseIntValue(data, "sellerId")
	if err != nil {
		return
	}

	pageInfo.DeliveryInfo, err = parseDeliveryInfoList(data["deliveryOptions"])

	return
}

type ResultHolder interface {
	Add(url string, pageInfo info.PageInfo)
}

func NewParser(resultHolder ResultHolder) func(*colly.Collector) {
	vm := goja.New()

	return func(c *colly.Collector) {
		c.OnHTML("script", func(e *colly.HTMLElement) {
			if !strings.Contains(e.Text, "app.run(") {
				return
			}

			url := e.Response.Ctx.Get("orig_url")

			v, err := vm.RunString(injectedJs + e.Text)
			if err != nil {
				log.Printf("[%s] error interpret js: %s", url, err)
				dumpJsonToFile(url, e.Text)
				return
			}

			if url != e.Request.URL.String() {
				log.Printf("redirect: %s -> %s", url, e.Request.URL.String())
			}

			pageInfo, err := parseInfo(v.Export())
			if err != nil {
				log.Printf("[%s] error interpret js: %s", url, err)
				dumpJsonToFile(url, e.Text)
				return
			}

			resultHolder.Add(url, pageInfo)
		})
	}
}
