package pdp

const parsedJs = `require.config({
		paths: {
			"react": "//g-assets.daily.taobao.net/lzdfe/pdp-platform/0.0.1/react.development",
			"react-dom": "//g-assets.daily.taobao.net/lzdfe/pdp-platform/0.0.1/react-dom.development",
			"@alife/universal-intl": "//unpkg.alibaba-inc.com/@alife/universal-intl@0.4.4/dist/universal-intl.min"
		}
  	});

	requirejs(['//g-assets.daily.taobao.net/lzdfe/pdp-platform/0.0.2/pc.js?_=1515407243982'], function(app) {
      app.run({"data": {"root": {"fields": {"primaryKey": {"skuId": 123, "itemId": 444, "sellerId":"10143"}, "deliveryOptions": {"123": [{"a": 1, "b": 2}]}}}}});
	})`
