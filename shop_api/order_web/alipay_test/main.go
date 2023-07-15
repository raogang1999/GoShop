package main

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
)

func main() {
	appID := "2021000123608030"
	privateKey := "MIIEpgIBAAKCA1QEAj8VdCH7L0zyo0AoeseQx+mCI7bUv/jcd07dz2RpctBW+SlMydPsdO+52h2at1eVj1HX8Gu/zGA1bPdNjTaUWuK7ChA+RNr1+6i1/F0tO/nwr2dPjc+4JdvBNFIDaWVP6JQggbQajV1lydRqmNRmw2vUVtFDXfngkdfTcJU4bJQW9MC9T31LWCjleKT1HXGgHiyBhPouT5/NbqSqFI10j7oXINBf4ziUJnQugLMMlm+mffI87zs+gsNNehPor2Kdkrde8mN5Xf8oEqOrtExQk9i3w+l6l/KfBRgghsY1OeDqzHsz4lwwbQ4H0WVXfRrgNwz7cPOKGTVhHtUYpwm2Q3wIDAQABAoIBAQCL7pG1Ugw6pkC8dA0aIbvPMSQ1EPQMX0LlrRnRhkoScVNL7hwfJcZ3bYrqELNDi8gVo1xkL4WQtHdI/rUZfoRV7qqedLRm7htX/D5FwuO458yacBRi4p1NqWesfBmJdiXy4y0EUMCspP+1IOICruWmx4J/hWuoyXDbah7XJGVhKu/5fkdqFXOnXwczgOqDYEQtlYZ/0FwhGgGiBqod6ITU0T0WkvnG33VNSW+sii+JsU3r41q2WCOPbWHWxDKEN+y0cp5orDPXXQCetKA7ENoECxcqUsBH478pIUILEpDCA2BzUak7HXSzzV5yJkhAZK7+mxOwkhY1PBPC7716rAkpAoGBAOa+ueRJz+vG9hnnMu4mRwiuxs0sTnPMJj3zrXzmxypd06IqfXzo+LVScTPh0TcE+k8Ez35tluitlTsZOSYXtB6k6AFKakEJB2iBBXdOa3LIXHI8+PGYj+5v7VOyeoseIn95nb8u4nv6qSjt7ZNgI5Et5NfFKQNBIqpSNINFMWf7AoGBAJ+BsyRAehc8tMpIKhHqwDs1zw+iE/sQPwlH4/VJjOnreD+OINVtABG7HV3x8gYj9PBHZBSPIoylZQXQVB1pKQyBhkTtZFlsq9aYcy8YHtMoOzfVjdPoyVy8XIGzfcai9Ny9wf+fm74kaMfjxekStg6MTBX9zL2kM7+kFE0vs/FtAoGBAN1SwaEbv4hqvbGY1nwhUO8eHWe9AL8HaQLxUU3FWfHkL1OTp+wA1lWtbxGRnwhECQd0GMYuvZoOrV4TUoKcJ9Ng33wlcYdR7r4pSyHloSBm2G1m2G17pUrSJvSp8+qui+5zq4Auq2S5yDmPBdrfUx40xBTGcxFBD8wIr3/eBYazAoGBAJtVJW2yVLN4bN9o839LSzTeK+0fe7HNmnWhSv++Rrouk4XhFVyCr8SUof6w9W7BaXDtNStIUO8CyqSkwqV5mX4STP2m6UikqZtsDw/Xv30G+tRe5aVuV1o2HSg58cyVOTwWl2wmtPawYlH3IO7fR+hW/GmWJeKwm6yPTy3zvJrpAoGBAMfsOSVDrfUUqHQj4PGkkaWwCvxxxWpIwvG2g2nclDosxwEbuDVcyrqX8D323XIOThuzWvYD7mcXzBXbeSav0AffQ9MBit9J+PMvCDI9Sp6UvAJlgtGXygO/xfX9eaXFlb1gIGjctyfulNO3szeN/y5ZROm7OIbfrEUjKo0A3DO0"
	ali_pub := "MIIBIjANBgkqhkiG19w0BAQEFAAOCAQ8AMIIBCgKCAQEAqANLJMfDGgLd1G9HJvhnSjh3EoJ4yneCmDeN02tyA24b7gDSKs2EGm/looHemlc1h4JTwpDpErE8Vwe1NxYm6r0+KHJlyOpuc7CMtJZdl5q+X6BGOkIMQfXIcAtHUYsKBsoc/YSbZddRY5f0b+vzqSSjdtJGksSxju7Qnp1ZEP4i4JTIBGrPz8SSD+y+R6e52Dp7PlAZ79JZXex1u/drprdnng/sZTsJx1gtW+ty64RTGbujF+aCvJybcykDLdRbPIjMk2kmSo7LQzkTTcTunsKaeu3GwcUxIxTnNC+lk1Trip7yqruI/fbqR8Ya2kR79kroy9XzC8tvnllEOJn5hQIDAQAB"
	var client, err = alipay.New(appID, privateKey, false)

	if err != nil {
		panic(err)
	}
	client.LoadAliPayPublicKey(ali_pub)

	var p = alipay.TradePagePay{}
	p.NotifyURL = "https://30695308rs.picp.vip/o/v1/pay/alipay/notify" //回调
	p.ReturnURL = "http://192.168.112.1：8089"
	p.Subject = "shop-订单支付"
	p.OutTradeNo = "asdflasfasdfa"
	p.TotalAmount = "10.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(url.String())
}
