package agent

// todo(fs): func checkForKey(key string, keyring *memberlist.Keyring) error {
// todo(fs): 	rk, err := base64.StdEncoding.DecodeString(key)
// todo(fs): 	if err != nil {
// todo(fs): 		return err
// todo(fs): 	}
// todo(fs):
// todo(fs): 	pk := keyring.GetPrimaryKey()
// todo(fs): 	if !bytes.Equal(rk, pk) {
// todo(fs): 		return fmt.Errorf("got %q want %q", pk, rk)
// todo(fs): 	}
// todo(fs): 	return nil
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_LoadKeyrings(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key := "tbLJg26ZJyJ9pK3qhc9jig=="
// todo(fs):
// todo(fs): 	// Should be no configured keyring file by default
// todo(fs): 	t.Run("no keys", func(t *testing.T) {
// todo(fs): 		a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a1.Shutdown()
// todo(fs):
// todo(fs): 		c1 := a1.Config.ConsulConfig
// todo(fs): 		if c1.SerfLANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("bad: %#v", c1.SerfLANConfig.KeyringFile)
// todo(fs): 		}
// todo(fs): 		if c1.SerfLANConfig.MemberlistConfig.Keyring != nil {
// todo(fs): 			t.Fatalf("keyring should not be loaded")
// todo(fs): 		}
// todo(fs): 		if c1.SerfWANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("bad: %#v", c1.SerfLANConfig.KeyringFile)
// todo(fs): 		}
// todo(fs): 		if c1.SerfWANConfig.MemberlistConfig.Keyring != nil {
// todo(fs): 			t.Fatalf("keyring should not be loaded")
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Server should auto-load LAN and WAN keyring files
// todo(fs): 	t.Run("server with keys", func(t *testing.T) {
// todo(fs): 		a2 := &TestAgent{Name: t.Name(), Key: key}
// todo(fs): 		a2.Start()
// todo(fs): 		defer a2.Shutdown()
// todo(fs):
// todo(fs): 		c2 := a2.Config.ConsulConfig
// todo(fs): 		if c2.SerfLANConfig.KeyringFile == "" {
// todo(fs): 			t.Fatalf("should have keyring file")
// todo(fs): 		}
// todo(fs): 		if c2.SerfLANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c2.SerfLANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if c2.SerfWANConfig.KeyringFile == "" {
// todo(fs): 			t.Fatalf("should have keyring file")
// todo(fs): 		}
// todo(fs): 		if c2.SerfWANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c2.SerfWANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Client should auto-load only the LAN keyring file
// todo(fs): 	t.Run("client with keys", func(t *testing.T) {
// todo(fs): 		cfg3 := TestConfig()
// todo(fs): 		cfg3.Server = false
// todo(fs): 		a3 := &TestAgent{Name: t.Name(), Config: cfg3, Key: key}
// todo(fs): 		a3.Start()
// todo(fs): 		defer a3.Shutdown()
// todo(fs):
// todo(fs): 		c3 := a3.Config.ConsulConfig
// todo(fs): 		if c3.SerfLANConfig.KeyringFile == "" {
// todo(fs): 			t.Fatalf("should have keyring file")
// todo(fs): 		}
// todo(fs): 		if c3.SerfLANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c3.SerfLANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if c3.SerfWANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("bad: %#v", c3.SerfWANConfig.KeyringFile)
// todo(fs): 		}
// todo(fs): 		if c3.SerfWANConfig.MemberlistConfig.Keyring != nil {
// todo(fs): 			t.Fatalf("keyring should not be loaded")
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_InmemKeyrings(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key := "tbLJg26ZJyJ9pK3qhc9jig=="
// todo(fs):
// todo(fs): 	// Should be no configured keyring file by default
// todo(fs): 	t.Run("no keys", func(t *testing.T) {
// todo(fs): 		a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a1.Shutdown()
// todo(fs):
// todo(fs): 		c1 := a1.Config.ConsulConfig
// todo(fs): 		if c1.SerfLANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("bad: %#v", c1.SerfLANConfig.KeyringFile)
// todo(fs): 		}
// todo(fs): 		if c1.SerfLANConfig.MemberlistConfig.Keyring != nil {
// todo(fs): 			t.Fatalf("keyring should not be loaded")
// todo(fs): 		}
// todo(fs): 		if c1.SerfWANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("bad: %#v", c1.SerfLANConfig.KeyringFile)
// todo(fs): 		}
// todo(fs): 		if c1.SerfWANConfig.MemberlistConfig.Keyring != nil {
// todo(fs): 			t.Fatalf("keyring should not be loaded")
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Server should auto-load LAN and WAN keyring
// todo(fs): 	t.Run("server with keys", func(t *testing.T) {
// todo(fs): 		cfg2 := TestConfig()
// todo(fs): 		cfg2.EncryptKey = key
// todo(fs): 		cfg2.DisableKeyringFile = true
// todo(fs):
// todo(fs): 		a2 := &TestAgent{Name: t.Name(), Config: cfg2}
// todo(fs): 		a2.Start()
// todo(fs): 		defer a2.Shutdown()
// todo(fs):
// todo(fs): 		c2 := a2.Config.ConsulConfig
// todo(fs): 		if c2.SerfLANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("should not have keyring file")
// todo(fs): 		}
// todo(fs): 		if c2.SerfLANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c2.SerfLANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if c2.SerfWANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("should not have keyring file")
// todo(fs): 		}
// todo(fs): 		if c2.SerfWANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c2.SerfWANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Client should auto-load only the LAN keyring
// todo(fs): 	t.Run("client with keys", func(t *testing.T) {
// todo(fs): 		cfg3 := TestConfig()
// todo(fs): 		cfg3.EncryptKey = key
// todo(fs): 		cfg3.DisableKeyringFile = true
// todo(fs): 		cfg3.Server = false
// todo(fs): 		a3 := &TestAgent{Name: t.Name(), Config: cfg3}
// todo(fs): 		a3.Start()
// todo(fs): 		defer a3.Shutdown()
// todo(fs):
// todo(fs): 		c3 := a3.Config.ConsulConfig
// todo(fs): 		if c3.SerfLANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("should not have keyring file")
// todo(fs): 		}
// todo(fs): 		if c3.SerfLANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c3.SerfLANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if c3.SerfWANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("bad: %#v", c3.SerfWANConfig.KeyringFile)
// todo(fs): 		}
// todo(fs): 		if c3.SerfWANConfig.MemberlistConfig.Keyring != nil {
// todo(fs): 			t.Fatalf("keyring should not be loaded")
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Any keyring files should be ignored
// todo(fs): 	t.Run("ignore files", func(t *testing.T) {
// todo(fs): 		dir := testutil.TempDir(t, "consul")
// todo(fs): 		defer os.RemoveAll(dir)
// todo(fs):
// todo(fs): 		badKey := "unUzC2X3JgMKVJlZna5KVg=="
// todo(fs): 		if err := initKeyring(filepath.Join(dir, SerfLANKeyring), badKey); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if err := initKeyring(filepath.Join(dir, SerfWANKeyring), badKey); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		cfg4 := TestConfig()
// todo(fs): 		cfg4.EncryptKey = key
// todo(fs): 		cfg4.DisableKeyringFile = true
// todo(fs): 		cfg4.DataDir = dir
// todo(fs):
// todo(fs): 		a4 := &TestAgent{Name: t.Name(), Config: cfg4}
// todo(fs): 		a4.Start()
// todo(fs): 		defer a4.Shutdown()
// todo(fs):
// todo(fs): 		c4 := a4.Config.ConsulConfig
// todo(fs): 		if c4.SerfLANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("should not have keyring file")
// todo(fs): 		}
// todo(fs): 		if c4.SerfLANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c4.SerfLANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if c4.SerfWANConfig.KeyringFile != "" {
// todo(fs): 			t.Fatalf("should not have keyring file")
// todo(fs): 		}
// todo(fs): 		if c4.SerfWANConfig.MemberlistConfig.Keyring == nil {
// todo(fs): 			t.Fatalf("keyring should be loaded")
// todo(fs): 		}
// todo(fs): 		if err := checkForKey(key, c4.SerfWANConfig.MemberlistConfig.Keyring); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_InitKeyring(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key1 := "tbLJg26ZJyJ9pK3qhc9jig=="
// todo(fs): 	key2 := "4leC33rgtXKIVUr9Nr0snQ=="
// todo(fs): 	expected := fmt.Sprintf(`["%s"]`, key1)
// todo(fs):
// todo(fs): 	dir := testutil.TempDir(t, "consul")
// todo(fs): 	defer os.RemoveAll(dir)
// todo(fs):
// todo(fs): 	file := filepath.Join(dir, "keyring")
// todo(fs):
// todo(fs): 	// First initialize the keyring
// todo(fs): 	if err := initKeyring(file, key1); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	content, err := ioutil.ReadFile(file)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if string(content) != expected {
// todo(fs): 		t.Fatalf("bad: %s", content)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try initializing again with a different key
// todo(fs): 	if err := initKeyring(file, key2); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Content should still be the same
// todo(fs): 	content, err = ioutil.ReadFile(file)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if string(content) != expected {
// todo(fs): 		t.Fatalf("bad: %s", content)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentKeyring_ACL(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key1 := "tbLJg26ZJyJ9pK3qhc9jig=="
// todo(fs): 	key2 := "4leC33rgtXKIVUr9Nr0snQ=="
// todo(fs):
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	a := &TestAgent{Name: t.Name(), Config: cfg, Key: key1}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// List keys without access fails
// todo(fs): 	_, err := a.ListKeys("", 0)
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "denied") {
// todo(fs): 		t.Fatalf("expected denied error, got: %#v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// List keys with access works
// todo(fs): 	_, err = a.ListKeys("root", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Install without access fails
// todo(fs): 	_, err = a.InstallKey(key2, "", 0)
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "denied") {
// todo(fs): 		t.Fatalf("expected denied error, got: %#v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Install with access works
// todo(fs): 	_, err = a.InstallKey(key2, "root", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Use without access fails
// todo(fs): 	_, err = a.UseKey(key2, "", 0)
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "denied") {
// todo(fs): 		t.Fatalf("expected denied error, got: %#v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Use with access works
// todo(fs): 	_, err = a.UseKey(key2, "root", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove without access fails
// todo(fs): 	_, err = a.RemoveKey(key1, "", 0)
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "denied") {
// todo(fs): 		t.Fatalf("expected denied error, got: %#v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove with access works
// todo(fs): 	_, err = a.RemoveKey(key1, "root", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): }
