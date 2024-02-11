package xmlparser

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"GoTestTask/pkg/db"
)

func ParseXML(ctx context.Context, url string, entries chan<- db.SdnEntry) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Ошибка декодирования XML:", err)
			return err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "sdnEntry" {
				var entry db.SdnEntry
				decoder.DecodeElement(&entry, &se)
				if entry.SdnType == "Individual" {
					select {
					case entries <- entry:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
		}
	}
	close(entries)
	return nil
}
