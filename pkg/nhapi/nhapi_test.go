package nhapi

import (
	"fmt"
	"testing"
	"time"
)

func TestClient_Get(t *testing.T) {
	client := New(Options{timeout: time.Minute, Url: NHentaiApiUrl})
	res, err := client.Get(312479)
	if err != nil {
		fmt.Printf("Test failed prematurely: %v", err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestClient_Search(t *testing.T) {
	client := New(Options{timeout: time.Minute, Url: NHentaiApiUrl})
	res, err := client.Search("jojo", "", 1)
	if err != nil {
		fmt.Printf("Test failed prematurely: %v", err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestResult_GetThumbnailUrl(t *testing.T) {
	testResult := Result{MediaID: "1234"}
	wantedUrl := "https://t.nhentai.net/galleries/" + testResult.MediaID + "/cover.jpg"
	gotUrl := testResult.GetThumbnailUrl()
	if wantedUrl != gotUrl {
		fmt.Printf("Wrong image URL\nGot: %v\nWanted:%v\n", gotUrl, wantedUrl)
		t.Fail()
	}
}