package nhapi

import (
	"fmt"
	"testing"
	"time"
)

func TestClient_Get(t *testing.T) {
	client := New(Options{Timeout: time.Minute, Url: NHentaiApiUrl})
	_, err := client.Get(323888)
	if err != nil {
		fmt.Printf("Test failed prematurely: %v", err)
		t.Fail()
	}
}

func TestClient_Search(t *testing.T) {
	client := New(Options{Timeout: time.Minute, Url: NHentaiApiUrl})
	_, err := client.Search("jojo", "", 1)
	if err != nil {
		fmt.Printf("Test failed prematurely: %v", err)
		t.Fail()
	}
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

func TestClient_Random(t *testing.T) {
	client := New(Options{Timeout: time.Minute, Url: NHentaiApiUrl})
	id, err := client.Random()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(id)
}
