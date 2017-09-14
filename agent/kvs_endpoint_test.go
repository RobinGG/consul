package agent

// todo(fs): func TestKVSEndpoint_PUT_GET_DELETE(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	keys := []string{
// todo(fs): 		"baz",
// todo(fs): 		"bar",
// todo(fs): 		"foo/sub1",
// todo(fs): 		"foo/sub2",
// todo(fs): 		"zip",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, key := range keys {
// todo(fs): 		buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/kv/"+key, buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, key := range keys {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/kv/"+key, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		res, ok := obj.(structs.DirEntries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(res) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", res)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res[0].Key != key {
// todo(fs): 			t.Fatalf("bad: %v", res)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, key := range keys {
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/kv/"+key, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.KVSEndpoint(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_Recurse(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	keys := []string{
// todo(fs): 		"bar",
// todo(fs): 		"baz",
// todo(fs): 		"foo/sub1",
// todo(fs): 		"foo/sub2",
// todo(fs): 		"zip",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, key := range keys {
// todo(fs): 		buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/kv/"+key, buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	{
// todo(fs): 		// Get all the keys
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/kv/?recurse", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		res, ok := obj.(structs.DirEntries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(res) != len(keys) {
// todo(fs): 			t.Fatalf("bad: %v", res)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		for idx, key := range keys {
// todo(fs): 			if res[idx].Key != key {
// todo(fs): 				t.Fatalf("bad: %v %v", res[idx].Key, key)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	{
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/kv/?recurse", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.KVSEndpoint(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	{
// todo(fs): 		// Get all the keys
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/kv/?recurse", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if obj != nil {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_DELETE_CAS(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/kv/test", buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/kv/test", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	d := obj.(structs.DirEntries)[0]
// todo(fs):
// todo(fs): 	// Create a CAS request, bad index
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte("zip"))
// todo(fs): 		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/kv/test?cas=%d", d.ModifyIndex-1), buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); res {
// todo(fs): 			t.Fatalf("should NOT work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a CAS request, good index
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte("zip"))
// todo(fs): 		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/kv/test?cas=%d", d.ModifyIndex), buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify the delete
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/kv/test", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, _ = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("should be destroyed")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_CAS(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/kv/test?flags=50", buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/kv/test", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	d := obj.(structs.DirEntries)[0]
// todo(fs):
// todo(fs): 	// Check the flags
// todo(fs): 	if d.Flags != 50 {
// todo(fs): 		t.Fatalf("bad: %v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a CAS request, bad index
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte("zip"))
// todo(fs): 		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/kv/test?flags=42&cas=%d", d.ModifyIndex-1), buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); res {
// todo(fs): 			t.Fatalf("should NOT work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a CAS request, good index
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte("zip"))
// todo(fs): 		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/kv/test?flags=42&cas=%d", d.ModifyIndex), buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify the update
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/kv/test", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, _ = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	d = obj.(structs.DirEntries)[0]
// todo(fs):
// todo(fs): 	if d.Flags != 42 {
// todo(fs): 		t.Fatalf("bad: %v", d)
// todo(fs): 	}
// todo(fs): 	if string(d.Value) != "zip" {
// todo(fs): 		t.Fatalf("bad: %v", d)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_ListKeys(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	keys := []string{
// todo(fs): 		"bar",
// todo(fs): 		"baz",
// todo(fs): 		"foo/sub1",
// todo(fs): 		"foo/sub2",
// todo(fs): 		"zip",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, key := range keys {
// todo(fs): 		buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/kv/"+key, buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	{
// todo(fs): 		// Get all the keys
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/kv/?keys&seperator=/", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		res, ok := obj.([]string)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		expect := []string{"bar", "baz", "foo/", "zip"}
// todo(fs): 		if !reflect.DeepEqual(res, expect) {
// todo(fs): 			t.Fatalf("bad: %v", res)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_AcquireRelease(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Acquire the lock
// todo(fs): 	id := makeTestSession(t, a.srv)
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/kv/test?acquire="+id, bytes.NewReader(nil))
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if res := obj.(bool); !res {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify we have the lock
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/kv/test", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	d := obj.(structs.DirEntries)[0]
// todo(fs):
// todo(fs): 	// Check the flags
// todo(fs): 	if d.Session != id {
// todo(fs): 		t.Fatalf("bad: %v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Release the lock
// todo(fs): 	req, _ = http.NewRequest("PUT", "/v1/kv/test?release="+id, bytes.NewReader(nil))
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if res := obj.(bool); !res {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify we do not have the lock
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/kv/test", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	d = obj.(structs.DirEntries)[0]
// todo(fs):
// todo(fs): 	// Check the flags
// todo(fs): 	if d.Session != "" {
// todo(fs): 		t.Fatalf("bad: %v", d)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_GET_Raw(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/kv/test", buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if res := obj.(bool); !res {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/kv/test?raw", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Check the body
// todo(fs): 	if !bytes.Equal(resp.Body.Bytes(), []byte("test")) {
// todo(fs): 		t.Fatalf("bad: %s", resp.Body.Bytes())
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_PUT_ConflictingFlags(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/kv/test?cas=0&acquire=xxx", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.KVSEndpoint(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): 	if !bytes.Contains(resp.Body.Bytes(), []byte("Conflicting")) {
// todo(fs): 		t.Fatalf("expected conflicting args error")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestKVSEndpoint_DELETE_ConflictingFlags(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("DELETE", "/v1/kv/test?recurse&cas=0", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.KVSEndpoint(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): 	if !bytes.Contains(resp.Body.Bytes(), []byte("Conflicting")) {
// todo(fs): 		t.Fatalf("expected conflicting args error")
// todo(fs): 	}
// todo(fs): }
