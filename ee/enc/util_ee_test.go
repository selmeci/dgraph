// +build !oss

/*
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Dgraph Community License (the "License"); you
 * may not use this file except in compliance with the License. You
 * may obtain a copy of the License at
 *
 *     https://github.com/dgraph-io/dgraph/blob/master/licenses/DCL.txt
 */

package enc

import (
	//"io"
	"crypto/cipher"
	//"net"
	"os"
	"testing"

	//"github.com/hashicorp/vault/api"
	// "github.com/hashicorp/vault/http"
	// "github.com/hashicorp/vault/vault"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func resetConfig(config *viper.Viper) {
	config.Set(encKeyFile, "")

	config.Set(vaultAddr, "http://localhost:8200")
	config.Set(vaultRoleIDFile, "")
	config.Set(vaultSecretIDFile, "")
	config.Set(vaultPath, "dgraph")
	config.Set(vaultField, "enc_key")
}

// TODO: The function below allows instantiating a real Vault server. But results in go.mod issues.
// func startVaultServer(t *testing.T, kvPath, kvField, kvEncKey string) (net.Listener, *api.Client) {
// 	core, _, rootToken := vault.TestCoreUnsealed(t)
// 	ln, addr := http.TestServer(t, core)
// 	t.Logf("addr = %v", addr)

// 	conf := api.DefaultConfig()
// 	conf.Address = addr
// 	client, err := api.NewClient(conf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	client.SetToken(rootToken)

// 	err = client.Sys().EnableAuthWithOptions("approle/", &api.EnableAuthOptions{
// 		Type: "approle",
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("approle enabled")

// 	return ln, client
// }

func TestNewKeyReader(t *testing.T) {
	config := viper.New()
	flags := &pflag.FlagSet{}
	RegisterFlags(flags)
	config.BindPFlags(flags)

	// Vault and Local combination tests
	// Both Local and Vault options is invalid.
	resetConfig(config)
	config.Set(encKeyFile, "blah")
	config.Set(vaultRoleIDFile, "aaa")
	config.Set(vaultSecretIDFile, "bbb")
	kR, err := NewKeyReader(config)
	require.Error(t, err)
	require.Nil(t, kR)

	// RoleID only is invalid.
	resetConfig(config)
	config.Set(vaultRoleIDFile, "aaa")
	kR, err = NewKeyReader(config)
	require.Error(t, err)
	require.Nil(t, kR)

	// SecretID only is invalid.
	resetConfig(config)
	config.Set(vaultSecretIDFile, "bbb")
	kR, err = NewKeyReader(config)
	require.Error(t, err)
	require.Nil(t, kR)

	// RoleID and SecretID given but RoleID file doesn't exist.
	resetConfig(config)
	config.Set(vaultRoleIDFile, "aaa")
	config.Set(vaultSecretIDFile, "bbb")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.NotNil(t, kR)
	require.IsType(t, &vaultKeyReader{}, kR)
	k, err := kR.ReadKey()
	require.Nil(t, k)
	require.Error(t, err)

	// RoleID and SecretID given but RoleID file exists. SecretID file doesn't exists.
	resetConfig(config)
	config.Set(vaultRoleIDFile, "./test-fixtures/dummy_role_id_file")
	config.Set(vaultSecretIDFile, "bbb")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.NotNil(t, kR)
	require.IsType(t, &vaultKeyReader{}, kR)
	k, err = kR.ReadKey()
	require.Nil(t, k)
	require.Error(t, err)

	// RoleID and SecretID given but RoleID file and SecretID file exists and is valid.
	resetConfig(config)
	//nl, _ := startVaultServer(t, "dgraph", "enc_key", "1234567890123456")

	config.Set(vaultRoleIDFile, "./test-fixtures/dummy_role_id_file")
	config.Set(vaultSecretIDFile, "./test-fixtures/dummy_secret_id_file")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.NotNil(t, kR)
	require.IsType(t, &vaultKeyReader{}, kR)
	k, err = kR.ReadKey()
	require.Nil(t, k) // still fails because we need to mock Vault server.
	require.Error(t, err)
	//nl.Close()

	// Bad Encryption Key File
	resetConfig(config)
	config.Set(encKeyFile, "blah")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.NotNil(t, kR)
	require.IsType(t, &localKeyReader{}, kR)
	k, err = kR.ReadKey()
	require.Nil(t, k)
	require.Error(t, err)

	// Nil Encryption Key File
	resetConfig(config)
	config.Set(encKeyFile, "")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.Nil(t, kR)

	// Bad Length Encryption Key File.
	resetConfig(config)
	config.Set(encKeyFile, "./test-fixtures/bad-length-enc-key")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.NotNil(t, kR)
	require.IsType(t, &localKeyReader{}, kR)
	k, err = kR.ReadKey()
	require.Nil(t, k)
	require.Error(t, err)

	// Good Encryption Key File.
	resetConfig(config)
	config.Set(encKeyFile, "./test-fixtures/enc-key")
	kR, err = NewKeyReader(config)
	require.NoError(t, err)
	require.NotNil(t, kR)
	require.IsType(t, &localKeyReader{}, kR)
	k, err = kR.ReadKey()
	require.NotNil(t, k)
	require.NoError(t, err)
}

func TestGetReaderWriter(t *testing.T) {
	// Test GetWriter()
	f, err := os.Create("/tmp/enc_test")
	require.NoError(t, err)
	defer os.Remove("/tmp/enc_test")

	// empty key
	neww, err := GetWriter(nil, f)
	require.NoError(t, err)
	require.Equal(t, f, neww)

	// valid key
	neww, err = GetWriter(ReadEncryptionKeyFile("./test-fixtures/enc-key"), f)
	require.NoError(t, err)
	require.NotEqual(t, f, neww)
	require.IsType(t, cipher.StreamWriter{}, neww)
	require.Equal(t, neww.(cipher.StreamWriter).W, f)
	// lets encrypt
	data := []byte("this is plaintext form")
	_, err = neww.Write(data)
	require.NoError(t, err)
	f.Close()

	// Test GetReader()
	f, err = os.Open("/tmp/enc_test")
	require.NoError(t, err)

	// empty key.
	newr, err := GetReader(nil, f)
	require.NoError(t, err)
	require.Equal(t, f, newr)

	// valid key
	newr, err = GetReader(ReadEncryptionKeyFile("./test-fixtures/enc-key"), f)
	require.NoError(t, err)
	require.NotEqual(t, f, newr)
	require.IsType(t, cipher.StreamReader{}, newr)
	require.Equal(t, newr.(cipher.StreamReader).R, f)

	// lets decrypt
	plain := make([]byte, len(data))
	_, err = newr.Read(plain)
	require.NoError(t, err)
	require.Equal(t, data, plain)
	f.Close()
}
