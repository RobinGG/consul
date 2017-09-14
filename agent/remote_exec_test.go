package agent

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
)

func generateUUID() (ret string) {
	var err error
	if ret, err = uuid.GenerateUUID(); err != nil {
		panic(fmt.Sprintf("Unable to generate a UUID, %v", err))
	}
	return ret
}

// todo(fs): func TestRexecWriter(t *testing.T) {
// todo(fs): 	// t.Parallel() // timing test. no parallel
// todo(fs): 	writer := &rexecWriter{
// todo(fs): 		BufCh:    make(chan []byte, 16),
// todo(fs): 		BufSize:  16,
// todo(fs): 		BufIdle:  100 * time.Millisecond,
// todo(fs): 		CancelCh: make(chan struct{}),
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Write short, wait for idle
// todo(fs): 	start := time.Now()
// todo(fs): 	n, err := writer.Write([]byte("test"))
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if n != 4 {
// todo(fs): 		t.Fatalf("bad: %v", n)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	select {
// todo(fs): 	case b := <-writer.BufCh:
// todo(fs): 		if len(b) != 4 {
// todo(fs): 			t.Fatalf("Bad: %v", b)
// todo(fs): 		}
// todo(fs): 		if time.Now().Sub(start) < writer.BufIdle {
// todo(fs): 			t.Fatalf("too early")
// todo(fs): 		}
// todo(fs): 	case <-time.After(2 * writer.BufIdle):
// todo(fs): 		t.Fatalf("timeout")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Write in succession to prevent the timeout
// todo(fs): 	writer.Write([]byte("test"))
// todo(fs): 	time.Sleep(writer.BufIdle / 2)
// todo(fs): 	writer.Write([]byte("test"))
// todo(fs): 	time.Sleep(writer.BufIdle / 2)
// todo(fs): 	start = time.Now()
// todo(fs): 	writer.Write([]byte("test"))
// todo(fs):
// todo(fs): 	select {
// todo(fs): 	case b := <-writer.BufCh:
// todo(fs): 		if len(b) != 12 {
// todo(fs): 			t.Fatalf("Bad: %v", b)
// todo(fs): 		}
// todo(fs): 		if time.Now().Sub(start) < writer.BufIdle {
// todo(fs): 			t.Fatalf("too early")
// todo(fs): 		}
// todo(fs): 	case <-time.After(2 * writer.BufIdle):
// todo(fs): 		t.Fatalf("timeout")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Write large values, multiple flushes required
// todo(fs): 	writer.Write([]byte("01234567890123456789012345678901"))
// todo(fs):
// todo(fs): 	select {
// todo(fs): 	case b := <-writer.BufCh:
// todo(fs): 		if string(b) != "0123456789012345" {
// todo(fs): 			t.Fatalf("bad: %s", b)
// todo(fs): 		}
// todo(fs): 	default:
// todo(fs): 		t.Fatalf("should have buf")
// todo(fs): 	}
// todo(fs): 	select {
// todo(fs): 	case b := <-writer.BufCh:
// todo(fs): 		if string(b) != "6789012345678901" {
// todo(fs): 			t.Fatalf("bad: %s", b)
// todo(fs): 		}
// todo(fs): 	default:
// todo(fs): 		t.Fatalf("should have buf")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecGetSpec(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	testRemoteExecGetSpec(t, nil, "", true)
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecGetSpec_ACLToken(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	testRemoteExecGetSpec(t, cfg, "root", true)
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecGetSpec_ACLAgentToken(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLAgentToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	testRemoteExecGetSpec(t, cfg, "root", true)
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecGetSpec_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	testRemoteExecGetSpec(t, cfg, "root", false)
// todo(fs): }
// todo(fs):
// todo(fs): func testRemoteExecGetSpec(t *testing.T, c *config.RuntimeConfig, token string, shouldSucceed bool) {
// todo(fs): 	a := NewTestAgent(t.Name(), c)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	event := &remoteExecEvent{
// todo(fs): 		Prefix:  "_rexec",
// todo(fs): 		Session: makeRexecSession(t, a.Agent, token),
// todo(fs): 	}
// todo(fs): 	defer destroySession(t, a.Agent, event.Session, token)
// todo(fs):
// todo(fs): 	spec := &remoteExecSpec{
// todo(fs): 		Command: "uptime",
// todo(fs): 		Script:  []byte("#!/bin/bash"),
// todo(fs): 		Wait:    time.Second,
// todo(fs): 	}
// todo(fs): 	buf, err := json.Marshal(spec)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	key := "_rexec/" + event.Session + "/job"
// todo(fs): 	setKV(t, a.Agent, key, buf, token)
// todo(fs):
// todo(fs): 	var out remoteExecSpec
// todo(fs): 	if shouldSucceed != a.remoteExecGetSpec(event, &out) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs): 	if shouldSucceed && !reflect.DeepEqual(spec, &out) {
// todo(fs): 		t.Fatalf("bad spec")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecWrites(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	testRemoteExecWrites(t, nil, "", true)
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecWrites_ACLToken(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	testRemoteExecWrites(t, cfg, "root", true)
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecWrites_ACLAgentToken(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLAgentToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	testRemoteExecWrites(t, cfg, "root", true)
// todo(fs): }
// todo(fs):
// todo(fs): func TestRemoteExecWrites_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	testRemoteExecWrites(t, cfg, "root", false)
// todo(fs): }
// todo(fs):
// todo(fs): func testRemoteExecWrites(t *testing.T, c *config.RuntimeConfig, token string, shouldSucceed bool) {
// todo(fs): 	a := NewTestAgent(t.Name(), c)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	event := &remoteExecEvent{
// todo(fs): 		Prefix:  "_rexec",
// todo(fs): 		Session: makeRexecSession(t, a.Agent, token),
// todo(fs): 	}
// todo(fs): 	defer destroySession(t, a.Agent, event.Session, token)
// todo(fs):
// todo(fs): 	if shouldSucceed != a.remoteExecWriteAck(event) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	output := []byte("testing")
// todo(fs): 	if shouldSucceed != a.remoteExecWriteOutput(event, 0, output) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs): 	if shouldSucceed != a.remoteExecWriteOutput(event, 10, output) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Bypass the remaining checks if the write was expected to fail.
// todo(fs): 	if !shouldSucceed {
// todo(fs): 		return
// todo(fs): 	}
// todo(fs):
// todo(fs): 	exitCode := 1
// todo(fs): 	if !a.remoteExecWriteExitCode(event, &exitCode) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	key := "_rexec/" + event.Session + "/" + a.Config.NodeName + "/ack"
// todo(fs): 	d := getKV(t, a.Agent, key, token)
// todo(fs): 	if d == nil || d.Session != event.Session {
// todo(fs): 		t.Fatalf("bad ack: %#v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	key = "_rexec/" + event.Session + "/" + a.Config.NodeName + "/out/00000"
// todo(fs): 	d = getKV(t, a.Agent, key, token)
// todo(fs): 	if d == nil || d.Session != event.Session || !bytes.Equal(d.Value, output) {
// todo(fs): 		t.Fatalf("bad output: %#v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	key = "_rexec/" + event.Session + "/" + a.Config.NodeName + "/out/0000a"
// todo(fs): 	d = getKV(t, a.Agent, key, token)
// todo(fs): 	if d == nil || d.Session != event.Session || !bytes.Equal(d.Value, output) {
// todo(fs): 		t.Fatalf("bad output: %#v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	key = "_rexec/" + event.Session + "/" + a.Config.NodeName + "/exit"
// todo(fs): 	d = getKV(t, a.Agent, key, token)
// todo(fs): 	if d == nil || d.Session != event.Session || string(d.Value) != "1" {
// todo(fs): 		t.Fatalf("bad output: %#v", d)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func testHandleRemoteExec(t *testing.T, command string, expectedSubstring string, expectedReturnCode string) {
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	event := &remoteExecEvent{
// todo(fs): 		Prefix:  "_rexec",
// todo(fs): 		Session: makeRexecSession(t, a.Agent, ""),
// todo(fs): 	}
// todo(fs): 	defer destroySession(t, a.Agent, event.Session, "")
// todo(fs):
// todo(fs): 	spec := &remoteExecSpec{
// todo(fs): 		Command: command,
// todo(fs): 		Wait:    time.Second,
// todo(fs): 	}
// todo(fs): 	buf, err := json.Marshal(spec)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	key := "_rexec/" + event.Session + "/job"
// todo(fs): 	setKV(t, a.Agent, key, buf, "")
// todo(fs):
// todo(fs): 	buf, err = json.Marshal(event)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	msg := &UserEvent{
// todo(fs): 		ID:      generateUUID(),
// todo(fs): 		Payload: buf,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Handle the event...
// todo(fs): 	a.handleRemoteExec(msg)
// todo(fs):
// todo(fs): 	// Verify we have an ack
// todo(fs): 	key = "_rexec/" + event.Session + "/" + a.Config.NodeName + "/ack"
// todo(fs): 	d := getKV(t, a.Agent, key, "")
// todo(fs): 	if d == nil || d.Session != event.Session {
// todo(fs): 		t.Fatalf("bad ack: %#v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify we have output
// todo(fs): 	key = "_rexec/" + event.Session + "/" + a.Config.NodeName + "/out/00000"
// todo(fs): 	d = getKV(t, a.Agent, key, "")
// todo(fs): 	if d == nil || d.Session != event.Session ||
// todo(fs): 		!bytes.Contains(d.Value, []byte(expectedSubstring)) {
// todo(fs): 		t.Fatalf("bad output: %#v", d)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify we have an exit code
// todo(fs): 	key = "_rexec/" + event.Session + "/" + a.Config.NodeName + "/exit"
// todo(fs): 	d = getKV(t, a.Agent, key, "")
// todo(fs): 	if d == nil || d.Session != event.Session || string(d.Value) != expectedReturnCode {
// todo(fs): 		t.Fatalf("bad output: %#v", d)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHandleRemoteExec(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	testHandleRemoteExec(t, "uptime", "load", "0")
// todo(fs): }
// todo(fs):
// todo(fs): func TestHandleRemoteExecFailed(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	testHandleRemoteExec(t, "echo failing;exit 2", "failing", "2")
// todo(fs): }
// todo(fs):
// todo(fs): func makeRexecSession(t *testing.T, a *Agent, token string) string {
// todo(fs): 	args := structs.SessionRequest{
// todo(fs): 		Datacenter: a.config.Datacenter,
// todo(fs): 		Op:         structs.SessionCreate,
// todo(fs): 		Session: structs.Session{
// todo(fs): 			Node:      a.config.NodeName,
// todo(fs): 			LockDelay: 15 * time.Second,
// todo(fs): 		},
// todo(fs): 		WriteRequest: structs.WriteRequest{
// todo(fs): 			Token: token,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out string
// todo(fs): 	if err := a.RPC("Session.Apply", &args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	return out
// todo(fs): }
// todo(fs):
// todo(fs): func destroySession(t *testing.T, a *Agent, session string, token string) {
// todo(fs): 	args := structs.SessionRequest{
// todo(fs): 		Datacenter: a.config.Datacenter,
// todo(fs): 		Op:         structs.SessionDestroy,
// todo(fs): 		Session: structs.Session{
// todo(fs): 			ID: session,
// todo(fs): 		},
// todo(fs): 		WriteRequest: structs.WriteRequest{
// todo(fs): 			Token: token,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out string
// todo(fs): 	if err := a.RPC("Session.Apply", &args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func setKV(t *testing.T, a *Agent, key string, val []byte, token string) {
// todo(fs): 	write := structs.KVSRequest{
// todo(fs): 		Datacenter: a.config.Datacenter,
// todo(fs): 		Op:         api.KVSet,
// todo(fs): 		DirEnt: structs.DirEntry{
// todo(fs): 			Key:   key,
// todo(fs): 			Value: val,
// todo(fs): 		},
// todo(fs): 		WriteRequest: structs.WriteRequest{
// todo(fs): 			Token: token,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var success bool
// todo(fs): 	if err := a.RPC("KVS.Apply", &write, &success); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func getKV(t *testing.T, a *Agent, key string, token string) *structs.DirEntry {
// todo(fs): 	req := structs.KeyRequest{
// todo(fs): 		Datacenter: a.config.Datacenter,
// todo(fs): 		Key:        key,
// todo(fs): 		QueryOptions: structs.QueryOptions{
// todo(fs): 			Token: token,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out structs.IndexedDirEntries
// todo(fs): 	if err := a.RPC("KVS.Get", &req, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if len(out.Entries) > 0 {
// todo(fs): 		return out.Entries[0]
// todo(fs): 	}
// todo(fs): 	return nil
// todo(fs): }
