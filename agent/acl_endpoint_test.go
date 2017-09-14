package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTestACL(t *testing.T, srv *HTTPServer) string {
	body := bytes.NewBuffer(nil)
	enc := json.NewEncoder(body)
	raw := map[string]interface{}{
		"Name":  "User Token",
		"Type":  "client",
		"Rules": "",
	}
	enc.Encode(raw)

	req, _ := http.NewRequest("PUT", "/v1/acl/create?token=root", body)
	resp := httptest.NewRecorder()
	obj, err := srv.ACLCreate(resp, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	aclResp := obj.(aclCreateResponse)
	return aclResp.ID
}

// todo(fs): func TestACL_Bootstrap(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLMasterToken = ""
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	tests := []struct {
// todo(fs): 		name   string
// todo(fs): 		method string
// todo(fs): 		code   int
// todo(fs): 		token  bool
// todo(fs): 	}{
// todo(fs): 		{"bad method", "GET", http.StatusMethodNotAllowed, false},
// todo(fs): 		{"bootstrap", "PUT", http.StatusOK, true},
// todo(fs): 		{"not again", "PUT", http.StatusForbidden, false},
// todo(fs): 	}
// todo(fs): 	for _, tt := range tests {
// todo(fs): 		t.Run(tt.name, func(t *testing.T) {
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			req, _ := http.NewRequest(tt.method, "/v1/acl/bootstrap", nil)
// todo(fs): 			out, err := a.srv.ACLBootstrap(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if got, want := resp.Code, tt.code; got != want {
// todo(fs): 				t.Fatalf("got %d want %d", got, want)
// todo(fs): 			}
// todo(fs): 			if tt.token {
// todo(fs): 				wrap, ok := out.(aclCreateResponse)
// todo(fs): 				if !ok {
// todo(fs): 					t.Fatalf("bad: %T", out)
// todo(fs): 				}
// todo(fs): 				if len(wrap.ID) != len("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx") {
// todo(fs): 					t.Fatalf("bad: %v", wrap)
// todo(fs): 				}
// todo(fs): 			} else {
// todo(fs): 				if out != nil {
// todo(fs): 					t.Fatalf("bad: %T", out)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Update(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestACL(t, a.srv)
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"ID":    id,
// todo(fs): 		"Name":  "User Token 2",
// todo(fs): 		"Type":  "client",
// todo(fs): 		"Rules": "",
// todo(fs): 	}
// todo(fs): 	enc.Encode(raw)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/acl/update?token=root", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.ACLUpdate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	aclResp := obj.(aclCreateResponse)
// todo(fs): 	if aclResp.ID != id {
// todo(fs): 		t.Fatalf("bad: %v", aclResp)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_UpdateUpsert(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"ID":    "my-old-id",
// todo(fs): 		"Name":  "User Token 2",
// todo(fs): 		"Type":  "client",
// todo(fs): 		"Rules": "",
// todo(fs): 	}
// todo(fs): 	enc.Encode(raw)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/acl/update?token=root", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.ACLUpdate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	aclResp := obj.(aclCreateResponse)
// todo(fs): 	if aclResp.ID != "my-old-id" {
// todo(fs): 		t.Fatalf("bad: %v", aclResp)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Destroy(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestACL(t, a.srv)
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/acl/destroy/"+id+"?token=root", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.ACLDestroy(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp, ok := obj.(bool); !ok || !resp {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/acl/info/"+id, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.ACLGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok := obj.(structs.ACLs)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", respObj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Clone(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestACL(t, a.srv)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/acl/clone/"+id, nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	_, err := a.srv.ACLClone(resp, req)
// todo(fs): 	if !acl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("PUT", "/v1/acl/clone/"+id+"?token=root", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.ACLClone(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	aclResp, ok := obj.(aclCreateResponse)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work: %#v %#v", obj, resp)
// todo(fs): 	}
// todo(fs): 	if aclResp.ID == id {
// todo(fs): 		t.Fatalf("bad id")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/acl/info/"+aclResp.ID, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.ACLGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok := obj.(structs.ACLs)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", respObj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Get(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("wrong id", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/acl/info/nope", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.ACLGet(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.ACLs)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if respObj == nil || len(respObj) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("right id", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		id := makeTestACL(t, a.srv)
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/acl/info/"+id, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.ACLGet(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.ACLs)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if len(respObj) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_List(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	var ids []string
// todo(fs): 	for i := 0; i < 10; i++ {
// todo(fs): 		ids = append(ids, makeTestACL(t, a.srv))
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/acl/list?token=root", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.ACLList(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok := obj.(structs.ACLs)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// 10 + anonymous + master
// todo(fs): 	if len(respObj) != 12 {
// todo(fs): 		t.Fatalf("bad: %v", respObj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACLReplicationStatus(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/acl/replication", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.ACLReplicationStatus(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	_, ok := obj.(structs.ACLReplicationStatus)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): }
