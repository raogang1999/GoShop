package model

// 过滤条件都建立在es_goods表上
type EsGoods struct {
	ID         int32 `json:"id"`
	CategoryID int32 `json:"category_id"`
	BrandID    int32 `json:"brand_id"`
	OnSale     bool  `json:"on_sale"`
	ShipFree   bool  `json:"ship_free"`
	IsNew      bool  `json:"is_new"`
	IsHot      bool  `json:"is_hot"`

	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

// 需要对应的mapping
func (EsGoods) GetMapping() string {
	goodsMapping := `
{
  "mappings": {
    "properties": {
      "brand_id": {
        "type": "integer"
      },
      "category_id": {
        "type": "integer"
      },
      "click_num": {
        "type": "integer"
      },
      "fav_num": {
        "type": "integer"
      },
      "id": {
        "type": "integer"
      },
      "is_hot": {
        "type": "boolean"
      },
      "is_new": {
        "type": "boolean"
      },
      "on_sale": {
        "type": "boolean"
      },
      "ship_free": {
        "type": "boolean"
      },
      "shop_price": {
        "type": "float"
      },
      "sold_num": {
        "type": "long"
      },
      "market_price": {
        "type": "float"
      },
      "name": {
        "type": "text",
        "analyzer": "ik_max_word",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "goods_brief": {
        "type": "text",
        "analyzer": "ik_max_word",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      }
    }
  }
}
`
	return goodsMapping
}

// 获得index名称
func (EsGoods) GetIndexName() string {
	return "goods"
}
