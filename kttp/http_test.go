package kttp

import (
	"log"
	"mime"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	TransportOptions(WithConnectTimeout(1000 * time.Millisecond))
	TransportOptions(WithInsecureSkipVerify(true))
	resp, err := NewRequest("https://bing.img.run/rand.php", nil, nil).Get()
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()

	result, err := resp.AsString()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
}

func TestPost(t *testing.T) {
	resp, err := NewRequest("https://www.baidu.com/?search=golang", nil, nil).Post()
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()

	result, err := resp.AsString()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
}

func TestAsDocument(t *testing.T) {
	resp, err := NewRequest("https://www.baidu.com", nil, nil).Get()
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()

	doc, err := resp.AsDom()
	if err != nil {
		t.Fatal(err)
	}

	head := doc.Find("head").Find("title")
	log.Println(head.Text())
}

func TestAsFile(t *testing.T) {
	resp, err := NewRequest("https://www.google.com/images/branding/googlelogo/2x/googlelogo_light_color_92x30dp.png", nil, nil).Get()
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()

	err = resp.AsFile("", "google")
	// str, err := resp.AsString()
	if err != nil {
		panic(err)
	}
	// log.Printf("%s", str)
}

func TestMime(t *testing.T) {
	testExtensionsByType("application/json")
	testExtensionsByType("image/jpeg")
	testExtensionsByType("image/png")
	testExtensionsByType("image/bmp")
	testExtensionsByType("text/plain")
	testExtensionsByType("text/html")
}

func testExtensionsByType(contentType string) {
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("extensions of %s: %s", contentType, strings.Join(extensions, ","))
	}
}

func TestExtractFileName(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Disposition", "inline")
	header.Add("Content-Disposition", "attachment; filename=test.txt")
	log.Printf("file name: %s", ExtractFileName(header))
}
