package agent

// todo(fs): func TestTxnEndpoint_Bad_JSON(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	buf := bytes.NewBuffer([]byte("{"))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.Txn(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): 	if !bytes.Contains(resp.Body.Bytes(), []byte("Failed to parse")) {
// todo(fs): 		t.Fatalf("expected conflicting args error")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestTxnEndpoint_Bad_Method(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	buf := bytes.NewBuffer([]byte("{}"))
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/txn", buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.Txn(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 405 {
// todo(fs): 		t.Fatalf("expected 405, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestTxnEndpoint_Bad_Size_Item(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	buf := bytes.NewBuffer([]byte(fmt.Sprintf(`
// todo(fs): [
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "set",
// todo(fs):             "Key": "key",
// todo(fs):             "Value": %q
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `, strings.Repeat("bad", 2*maxKVSize))))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.Txn(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 413 {
// todo(fs): 		t.Fatalf("expected 413, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestTxnEndpoint_Bad_Size_Net(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	value := strings.Repeat("X", maxKVSize/2)
// todo(fs): 	buf := bytes.NewBuffer([]byte(fmt.Sprintf(`
// todo(fs): [
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "set",
// todo(fs):             "Key": "key1",
// todo(fs):             "Value": %q
// todo(fs):         }
// todo(fs):     },
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "set",
// todo(fs):             "Key": "key1",
// todo(fs):             "Value": %q
// todo(fs):         }
// todo(fs):     },
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "set",
// todo(fs):             "Key": "key1",
// todo(fs):             "Value": %q
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `, value, value, value)))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.Txn(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 413 {
// todo(fs): 		t.Fatalf("expected 413, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestTxnEndpoint_Bad_Size_Ops(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	buf := bytes.NewBuffer([]byte(fmt.Sprintf(`
// todo(fs): [
// todo(fs):     %s
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "set",
// todo(fs):             "Key": "key",
// todo(fs):             "Value": ""
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `, strings.Repeat(`{ "KV": { "Verb": "get", "Key": "key" } },`, 2*maxTxnOps))))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.Txn(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 413 {
// todo(fs): 		t.Fatalf("expected 413, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestTxnEndpoint_KV_Actions(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		// Make sure all incoming fields get converted properly to the internal
// todo(fs): 		// RPC format.
// todo(fs): 		var index uint64
// todo(fs): 		id := makeTestSession(t, a.srv)
// todo(fs): 		{
// todo(fs): 			buf := bytes.NewBuffer([]byte(fmt.Sprintf(`
// todo(fs): [
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "lock",
// todo(fs):             "Key": "key",
// todo(fs):             "Value": "aGVsbG8gd29ybGQ=",
// todo(fs):             "Flags": 23,
// todo(fs):             "Session": %q
// todo(fs):         }
// todo(fs):     },
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "get",
// todo(fs):             "Key": "key"
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `, id)))
// todo(fs): 			req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.Txn(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if resp.Code != 200 {
// todo(fs): 				t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			txnResp, ok := obj.(structs.TxnResponse)
// todo(fs): 			if !ok {
// todo(fs): 				t.Fatalf("bad type: %T", obj)
// todo(fs): 			}
// todo(fs): 			if len(txnResp.Results) != 2 {
// todo(fs): 				t.Fatalf("bad: %v", txnResp)
// todo(fs): 			}
// todo(fs): 			index = txnResp.Results[0].KV.ModifyIndex
// todo(fs): 			expected := structs.TxnResponse{
// todo(fs): 				Results: structs.TxnResults{
// todo(fs): 					&structs.TxnResult{
// todo(fs): 						KV: &structs.DirEntry{
// todo(fs): 							Key:       "key",
// todo(fs): 							Value:     nil,
// todo(fs): 							Flags:     23,
// todo(fs): 							Session:   id,
// todo(fs): 							LockIndex: 1,
// todo(fs): 							RaftIndex: structs.RaftIndex{
// todo(fs): 								CreateIndex: index,
// todo(fs): 								ModifyIndex: index,
// todo(fs): 							},
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 					&structs.TxnResult{
// todo(fs): 						KV: &structs.DirEntry{
// todo(fs): 							Key:       "key",
// todo(fs): 							Value:     []byte("hello world"),
// todo(fs): 							Flags:     23,
// todo(fs): 							Session:   id,
// todo(fs): 							LockIndex: 1,
// todo(fs): 							RaftIndex: structs.RaftIndex{
// todo(fs): 								CreateIndex: index,
// todo(fs): 								ModifyIndex: index,
// todo(fs): 							},
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			if !reflect.DeepEqual(txnResp, expected) {
// todo(fs): 				t.Fatalf("bad: %v", txnResp)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Do a read-only transaction that should get routed to the
// todo(fs): 		// fast-path endpoint.
// todo(fs): 		{
// todo(fs): 			buf := bytes.NewBuffer([]byte(`
// todo(fs): [
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "get",
// todo(fs):             "Key": "key"
// todo(fs):         }
// todo(fs):     },
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "get-tree",
// todo(fs):             "Key": "key"
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `))
// todo(fs): 			req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.Txn(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if resp.Code != 200 {
// todo(fs): 				t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			header := resp.Header().Get("X-Consul-KnownLeader")
// todo(fs): 			if header != "true" {
// todo(fs): 				t.Fatalf("bad: %v", header)
// todo(fs): 			}
// todo(fs): 			header = resp.Header().Get("X-Consul-LastContact")
// todo(fs): 			if header != "0" {
// todo(fs): 				t.Fatalf("bad: %v", header)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			txnResp, ok := obj.(structs.TxnReadResponse)
// todo(fs): 			if !ok {
// todo(fs): 				t.Fatalf("bad type: %T", obj)
// todo(fs): 			}
// todo(fs): 			expected := structs.TxnReadResponse{
// todo(fs): 				TxnResponse: structs.TxnResponse{
// todo(fs): 					Results: structs.TxnResults{
// todo(fs): 						&structs.TxnResult{
// todo(fs): 							KV: &structs.DirEntry{
// todo(fs): 								Key:       "key",
// todo(fs): 								Value:     []byte("hello world"),
// todo(fs): 								Flags:     23,
// todo(fs): 								Session:   id,
// todo(fs): 								LockIndex: 1,
// todo(fs): 								RaftIndex: structs.RaftIndex{
// todo(fs): 									CreateIndex: index,
// todo(fs): 									ModifyIndex: index,
// todo(fs): 								},
// todo(fs): 							},
// todo(fs): 						},
// todo(fs): 						&structs.TxnResult{
// todo(fs): 							KV: &structs.DirEntry{
// todo(fs): 								Key:       "key",
// todo(fs): 								Value:     []byte("hello world"),
// todo(fs): 								Flags:     23,
// todo(fs): 								Session:   id,
// todo(fs): 								LockIndex: 1,
// todo(fs): 								RaftIndex: structs.RaftIndex{
// todo(fs): 									CreateIndex: index,
// todo(fs): 									ModifyIndex: index,
// todo(fs): 								},
// todo(fs): 							},
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 				QueryMeta: structs.QueryMeta{
// todo(fs): 					KnownLeader: true,
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			if !reflect.DeepEqual(txnResp, expected) {
// todo(fs): 				t.Fatalf("bad: %v", txnResp)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Now that we have an index we can do a CAS to make sure the
// todo(fs): 		// index field gets translated to the RPC format.
// todo(fs): 		{
// todo(fs): 			buf := bytes.NewBuffer([]byte(fmt.Sprintf(`
// todo(fs): [
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "cas",
// todo(fs):             "Key": "key",
// todo(fs):             "Value": "Z29vZGJ5ZSB3b3JsZA==",
// todo(fs):             "Index": %d
// todo(fs):         }
// todo(fs):     },
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "get",
// todo(fs):             "Key": "key"
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `, index)))
// todo(fs): 			req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.Txn(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if resp.Code != 200 {
// todo(fs): 				t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			txnResp, ok := obj.(structs.TxnResponse)
// todo(fs): 			if !ok {
// todo(fs): 				t.Fatalf("bad type: %T", obj)
// todo(fs): 			}
// todo(fs): 			if len(txnResp.Results) != 2 {
// todo(fs): 				t.Fatalf("bad: %v", txnResp)
// todo(fs): 			}
// todo(fs): 			modIndex := txnResp.Results[0].KV.ModifyIndex
// todo(fs): 			expected := structs.TxnResponse{
// todo(fs): 				Results: structs.TxnResults{
// todo(fs): 					&structs.TxnResult{
// todo(fs): 						KV: &structs.DirEntry{
// todo(fs): 							Key:     "key",
// todo(fs): 							Value:   nil,
// todo(fs): 							Session: id,
// todo(fs): 							RaftIndex: structs.RaftIndex{
// todo(fs): 								CreateIndex: index,
// todo(fs): 								ModifyIndex: modIndex,
// todo(fs): 							},
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 					&structs.TxnResult{
// todo(fs): 						KV: &structs.DirEntry{
// todo(fs): 							Key:     "key",
// todo(fs): 							Value:   []byte("goodbye world"),
// todo(fs): 							Session: id,
// todo(fs): 							RaftIndex: structs.RaftIndex{
// todo(fs): 								CreateIndex: index,
// todo(fs): 								ModifyIndex: modIndex,
// todo(fs): 							},
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			if !reflect.DeepEqual(txnResp, expected) {
// todo(fs): 				t.Fatalf("bad: %v", txnResp)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Verify an error inside a transaction.
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		buf := bytes.NewBuffer([]byte(`
// todo(fs): [
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "lock",
// todo(fs):             "Key": "key",
// todo(fs):             "Value": "aGVsbG8gd29ybGQ=",
// todo(fs):             "Session": "nope"
// todo(fs):         }
// todo(fs):     },
// todo(fs):     {
// todo(fs):         "KV": {
// todo(fs):             "Verb": "get",
// todo(fs):             "Key": "key"
// todo(fs):         }
// todo(fs):     }
// todo(fs): ]
// todo(fs): `))
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/txn", buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.Txn(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 409 {
// todo(fs): 			t.Fatalf("expected 409, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		if !bytes.Contains(resp.Body.Bytes(), []byte("failed session lookup")) {
// todo(fs): 			t.Fatalf("bad: %s", resp.Body.String())
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
