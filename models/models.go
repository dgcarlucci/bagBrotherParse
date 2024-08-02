package models

// Item, holds the local id, name and quantity for an item
type Item struct {
	Id       string
	Name     string
	Quantity int
}

// WowItem struct
type WoWItem struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Quality struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"quality"`
	Level        int `json:"level"`
	RequiredLvel int `json:"required_level"`
	Media        struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		ID int `json:"id"`
	} `json:"media"`
	ItemClass struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"item_class"`
	ItemSubclass struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"item_subclass"`
	InventoryType struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"inventory_type"`
	PurchasePrice int  `json:"purchase_price"`
	SellPrice     int  `json:"sell_price"`
	MaxCount      int  `json:"max_count"`
	IsEquippable  bool `json:"is_equippable"`
	IsStackable   bool `json:"is_stackable"`
	PreviewItem   struct {
		Item struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			ID int `json:"id"`
		} `json:"item"`
		Quality struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"quality"`
		Name  string `json:"name"`
		Media struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			ID int `json:"id"`
		} `json:"media"`
		ItemClass struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"item_class"`
		ItemSubclass struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"item_subclass"`
		InventoryType struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"inventory_type"`
		SellPrice struct {
			Value          int `json:"value"`
			DisplayStrings struct {
				Header string `json:"header"`
				Gold   string `json:"gold"`
				Silver string `json:"silver"`
				Copper string `json:"copper"`
			} `json:"display_strings"`
		} `json:"sell_price"`
		ContainerSlots struct {
			Value         int    `json:"value"`
			DisplayString string `json:"display_string"`
		} `json:"container_slots"`
	} `json:"preview_item"`
	PurchaseQuantity int `json:"purchase_quantity"`
}
