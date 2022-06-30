package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

var urls = map[string]string{
	"tc":   "http://www.tcmb.gov.tr/kurlar/today.xml",
	"kktc": "http://www.mb.gov.ct.tr/kur/gunluk.xml",
}

func main() {
	// Tc()
	Kktc()
}

func Tc() {
	var res = downloadString(urls["tc"])
	var tc Currency
	if err := xml.Unmarshal([]byte(res), &tc); err != nil {
		panic(err)
	}
	for index, el := range tc.Currency {
		str := fmt.Sprint(index, `) `, el.Name, el.ForexBuying, el.ForexSelling)
		fmt.Println(str)
	}
}

func Kktc() {
	var res = downloadString(urls["kktc"])
	var tc CurrencyKKTC
	if err := xml.Unmarshal([]byte(res), &tc); err != nil {
		panic(err)
	}
	for index, el := range tc.Currencies {
		str := fmt.Sprint(index, `) `, el.Name, el.Buying, el.Selling)
		fmt.Println(str)
	}
}

func downloadString(param string) string {
	res, err := http.Get(param)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return string(body)
}

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortform = "20060102"
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortform, v)
	if err != nil {
		return err
	}
	*c = customTime{parse}
	return nil
}

type Currency struct {
	XMLName  xml.Name `xml:"Tarih_Date"`
	Currency []CurrencyItem
}

type CurrencyItem struct {
	XMLName         xml.Name `xml:"Currency"`
	Unit            int      `xml:"Unit"`
	Name            string   `xml:"Isim"`
	Currencyname    string   `xml:"CurrencyName"`
	ForexBuying     float32  `xml:"ForexBuying"`
	ForexSelling    float32  `xml:"ForexSelling"`
	BanknoteBuying  float32  `xml:"BanknoteBuying"`
	BanknoteSelling float32  `xml:"BanknoteSelling"`
	CrossRateUSD    float32  `xml:"CrossRateUSD"`
	CrossRateOther  float32  `xml:"CrossRateOther"`
}

type CurrencyKKTC struct {
	XMLName         xml.Name            `xml:"KKTCMB_Doviz_Kurlari"`
	CurrencyDate    customTime          `xml:"Kur_Tarihi"`
	AnnouncementNo  customTime          `xml:"Duyuru_No"`
	ValidDate       customTime          `xml:"Gecerli_Tarih_Araligi"`
	ValidDateStr    customTime          `xml:"Gecerli_Tarih_Araligi_Str"`
	Currencies      []CurrencyItemKKTC  `xml:"Resmi_Kurlar"`
	CrossCurrencies []CrossCurrencyItem `xml:"Capraz_Kurlar"`
}

type CurrencyItemKKTC struct {
	XMLName          xml.Name `xml:"Resmi_Kurlar"`
	Unit             int      `xml:"Birim"`
	Symbol           int      `xml:"Sembol"`
	Name             int      `xml:"Isim"`
	Buying           int      `xml:"Doviz_Alis"`
	Selling          int      `xml:"Doviz_Satis"`
	EffectiveBuying  int      `xml:"Efektif_Alis"`
	EffectiveSelling int      `xml:"Efektif_Satis"`
}

type CrossCurrencyItem struct {
	XMLName            xml.Name `xml:"Capraz_Kurlar"`
	SymbolFrom         string   `xml:"Sembol_From"`
	CurrencyFrom       float32  `xml:"Doviz_From"`
	CrossCurrencyValue float32  `xml:"Capraz_Kur_Deger"`
	CurrencyNameTo     string   `xml:"Doviz_To"`
	SymbolTo           string   `xml:"Sembol_To"`
}
