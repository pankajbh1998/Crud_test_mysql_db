package main

import (
	"fmt"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpGET(t *testing.T){
	testCases:=[]struct{
		id string
		output string
		statusCode int
	}{
		{
			"100","{\"Id\":\"100\",\"Name\":\"RB8\",\"Age\":35,\"Gender\":\"M\",\"Role\":\"4\"}",200,
		},
		{
			"104","{\"Id\":\"104\",\"Name\":\"Pankaj\",\"Age\":22,\"Gender\":\"M\",\"Role\":\"1\"}",200,
		},
		{
			"101","",400,//does not exist
		},
	}
	for i,tc:=range testCases {
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("GET","/employee",nil)
		r=mux.SetURLVars(r,map[string]string{
			"id": tc.id,
		})
		ReadDataId(w,r)
		res:=w.Result()
		readbyte,err:=ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(readbyte))
		if tc.statusCode!=res.StatusCode {
			t.Fatalf("statuscodes are not matching")
		} else if tc.statusCode !=http.StatusBadRequest && string(readbyte) !=	tc.output {
			t.Fatalf("Failed at %v\nOutput : %v\nActual : %v",i+1,tc.output,string(readbyte))
		}
		t.Logf("Passed at %v",i+1)
	}
}
func TestHttpPost(t *testing.T){
	testCases:=[]struct{
		inp string
		output string
		statusCode int
	}{
		{
			"{\"name\":\"RB8\",\"age\":35,\"Gender\":\"M\",\"role\":\"4\"}","RB8 is Added Successfully to Database and Generated id is 133",200,
		},
	}
	for i,tc:=range testCases {
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("POST","/employee",strings.NewReader(tc.inp))
		CreateData(w,r)
		res:=w.Result()
		readbyte,err:=ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(readbyte))
		if tc.statusCode!=res.StatusCode {
			t.Fatalf("statuscodes are not matching")
		} else if tc.statusCode !=http.StatusBadRequest && string(readbyte) !=	tc.output {
			t.Fatalf("Failed at %v\nOutput : %v\nActual : %v",i+1,tc.output,string(readbyte))
		}
		t.Logf("Passed at %v",i+1)
	}
}


func TestHttpPut(t *testing.T){
	testCases:=[]struct{
		id string
		inp string
		output string
		statusCode int
	}{
		{
			"104","{\"Name\":\"Pankaj\",\"Age\":22,\"Gender\":\"M\",\"Role\":\"1\"}","Updated Successfully",200,
		},
	}
	for i,tc:=range testCases{
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("PUT","/employee",strings.NewReader(tc.inp))
		r=mux.SetURLVars(r,map[string]string{
			"id":tc.id,
		})
		UpdateData(w,r)
		res:=w.Result()
		readbyte,err:=ioutil.ReadAll(res.Body)

		if err !=nil {
			t.Fatal(err)
		}
		fmt.Println(string(readbyte))
		if tc.statusCode !=res.StatusCode {
			t.Fatalf("status codes are not matching")
		}
		if tc.output != string(readbyte) {
			t.Fatalf("Failed at %v\nExpected Output : %v\nActual Output: %v",i+1,tc.output,string(readbyte))
		}
		t.Logf("Passed at %v",i+1)
	}
}
