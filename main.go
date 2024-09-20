package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"main.go/models"
)

const (
	clientSecret = "your client secret"
	baseURL      = "https://us.api.blizzard.com/data/wow/item/"
	locale       = "en_US"
	namespace    = "static-us"
)

type config struct {
	InputFilePath   string `json:"input_file_path"`
	OutputDirectory string `json:"output_directory"`
	ServerName      string `json:"server_name"`
	CharacterName   string `json:"character_name"`
	ClientID        string `json:"client_id"`
	ClientSecret    string `json:"client_secret"`
}

func main() {
	var config config

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
	processLuaFile(config)
}

func processLuaFile(config config) {
	token, err := fetchAccessToken(config.ClientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadFile(config.InputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(config.OutputDirectory, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}
	fileInfo, err := os.Stat(config.InputFilePath)
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}
	timestamp := fileInfo.ModTime()
	filename := timestamp.Format("2006-01-02T15-04-05") + ".csv"
	csvFilePath := filepath.Join(config.OutputDirectory, filename)
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(string(data)); err != nil {
		log.Fatal(err)
	}
	processLuaTable(L, config.ServerName, config.CharacterName, token, csvFilePath)
}

func processLuaTable(L *lua.LState, serverName string, characterName string, token string, csvFilePath string) {
	bortherBagsSaved := L.GetGlobal("BrotherBags").(*lua.LTable)
	scans := L.GetField(bortherBagsSaved, serverName).(*lua.LTable)

	file, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '|'
	defer writer.Flush()

	writer.Write([]string{"Item", "Qty"})

	processTable(scans, characterName, token, writer)
}

func processTable(table *lua.LTable, characterName string, token string, writer *csv.Writer) {
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			currentCharacterName := key.String()
			if currentCharacterName != characterName {
				return
			}

			processInnerTable(tbl, token, writer)
		}
	})
}

func processInnerTable(table *lua.LTable, token string, writer *csv.Writer) {
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			processDeepTable(tbl, token, writer)
		}
	})
}

func processDeepTable(table *lua.LTable, token string, writer *csv.Writer) {
	itemQuantities := make(map[string]int)

	table.ForEach(func(key, value lua.LValue) {
		itemQuantity, err := getItemQuantity(value.String())
		if err != nil {
			return
		}

		itemQuantities[itemQuantity.Id] += itemQuantity.Quantity
	})

	for itemID, quantity := range itemQuantities {
		itemName, err := getItemName(itemID, token)
		if err != nil {
			continue
		}

		if isIgnoredItem(itemName) {
			continue
		}

		itemLink := fmt.Sprintf("[%s](%s)", itemName, getWoWheadLink(itemID))
		writer.Write([]string{itemLink, strconv.Itoa(quantity)})
	}
}

func isIgnoredItem(itemName string) bool {
	ignoredItemNames := []string{"Bag", "Hearthstone", "Pack", "Backpack"}
	for _, name := range ignoredItemNames {
		if strings.Contains(itemName, name) {
			return true
		}
	}
	return false
}

// getItemQuantity , parses the item quantity, separate by ::::::::
// ex: 17012::::::::1:::::::::;20
// returns item=17012 and quantity=20
func getItemQuantity(input string) (item models.Item, err error) {
	log.Println("getItemQuantity", input)
	parts := strings.Split(input, ":")

	if len(parts) < 17 {
		return item, fmt.Errorf("invalid input: %s", input)
	}

	log.Println("parts0", parts[0])
	item.Id = parts[0]

	quantity := "1"
	quantityStr := strings.Split(parts[17], ";")

	if len(quantityStr) == 2 {
		quantity = quantityStr[1]
	}

	qty := 0

	if qty, err = strconv.Atoi(quantity); err != nil {
		return item, err
	}

	item.Quantity = qty

	return item, nil
}

//getWoWheadLink

func getWoWheadLink(itemId string) string {

	return fmt.Sprintf("https://classic.wowhead.com/item=%s", itemId)

}

//getWoWItemFromAPI

// getItemName fetches the name of an item using its ID and access token
func getItemName(id string, token string) (string, error) {

	var item models.Item
	url := fmt.Sprintf("%s%s?namespace=%s&locale=%s&access_token=%s", baseURL, id, namespace, locale, token)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//printout the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Println(string(body))

	wowItem := models.WoWItem{}
	err = json.Unmarshal(body, &wowItem)
	if err != nil {
		return "", err
	}

	item.Name = wowItem.Name

	return item.Name, nil
}

// fetchAccessToken
func fetchAccessToken(clientID string, clientSecret string) (string, error) {
	tokenURL := fmt.Sprintf("https://us.battle.net/oauth/token?grant_type=client_credentials&client_id=%s&client_secret=%s", clientID, clientSecret)
	request, err := http.NewRequest(http.MethodPost, tokenURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
