package worker

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/melbahja/goph"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/crypto/ssh"
)

const AUTHORIZED_KEYS_START_MARKER = "# === ssh-authorized-manager START marker ===>"
const AUTHORIZED_KEYS_END_MARKER = "# <=== ssh-authorized-manager END marker ==="

type AuthorizedKey struct {
	Options string
	Type    string
	Key     ssh.PublicKey
	Comment string
}

func (key AuthorizedKey) Serialize() string {
	// NOTE authorized_keys line format is "<options?> <type> <base64 encoded key> <comment?>". no spaces allowed in any section.
	// See sshd(8) for more information.
	parts := []string{key.Type, serializePublicKey(key.Key)}
	if key.Options != "" {
		// prepend options
		parts = append([]string{key.Options}, parts...)
	}
	if key.Comment != "" {
		parts = append(parts, key.Comment)
	}
	return strings.Join(parts, " ")
}

type SyncServerWork struct {
	Server *models.Record
}

func (work *SyncServerWork) Name() string {
	return "Sync Server - " + work.Server.GetString("name")
}

func (work *SyncServerWork) Execute() error {
	// Prepare auth
	var auth goph.Auth
	if work.Server.GetBool("usePassword") {
		auth = goph.Password(work.Server.GetString("#password"))
	} else {
		var err error
		auth, err = goph.Key(work.Server.GetString("#privateKey"), work.Server.GetString("#privateKeyPassphrase"))
		if err != nil {
			CreateServerLog(work.Server, "error", "Failed to load private key. Passphrase may be incorrect.", err.Error())
			return err
		}
	}

	// Prepare client and verify host key
	client, err := goph.NewConn(&goph.Config{
		User:     work.Server.GetString("username"),
		Addr:     work.Server.GetString("host"),
		Port:     uint(work.Server.GetInt("port")),
		Auth:     auth,
		Callback: verifyHostKey(work),
	})
	if err != nil {
		// check if ends with "host key is not known"
		if !strings.HasSuffix(err.Error(), "ssh: required host key is not known") {
			CreateServerLog(work.Server, "error", "Failed to connect to server.", err.Error())
		}
		return err
	}
	defer client.Close()
	fmt.Println("[WORKER] Connected to", work.Server.GetString("name"))

	// Verify hostname
	if work.Server.GetString("hostname") != "" {
		if _, err := verifyServerHostname(client, work.Server.GetString("hostname")); err != nil {
			CreateServerLog(work.Server, "error", "Failed to verify hostname. Possible port redirection.", err.Error())
			return err
		}
		fmt.Println("[WORKER] Verified hostname", work.Server.GetString("hostname"))
	} else {
		// hostname is unknown, ask to if hostname is trusted
		hostname, _ := verifyServerHostname(client, work.Server.GetString("hostname"))
		msg := fmt.Sprintf("Hostname does not match, do you trust this server?\nHostname: %s", hostname)
		CreateServerLog(work.Server, "hostName", msg, hostname)
		return fmt.Errorf("ssh: required hostname is not known")
	}

	// collect authorized keys
	authorizedKeys, err := collectAuthorizedKeys(app.Dao(), work.Server)
	if err != nil {
		CreateServerLog(work.Server, "error", "Failed to collect authorized keys.", err.Error())
		return err
	}

	// inject authorized keys into our section
	// We may be able to do this better within the SSH client instead of download/upload, but this is a good start.
	if err := client.Download(".ssh/authorized_keys", "/tmp/authorized_keys"); err != nil {
		CreateServerLog(work.Server, "error", "Failed to download authorized keys. Is the file `~/.ssh/authorized_keys` missing?", err.Error())
		return err
	}
	startLines, _, endLines, err := parseAuthorizedKeys("/tmp/authorized_keys")
	if err != nil {
		CreateServerLog(work.Server, "error", "Failed to parse authorized_keys file.", err.Error())
		return err
	}
	start := []byte(strings.Join(startLines, "\n"))
	end := []byte(strings.Join(endLines, "\n"))
	if err := assembleAuthorizedKeys("/tmp/authorized_keys", start, authorizedKeys, end); err != nil {
		CreateServerLog(work.Server, "error", "Failed to assemble authorized_keys file.", err.Error())
		return err
	}
	client.Upload("/tmp/authorized_keys", ".ssh/authorized_keys")
	fmt.Println("[WORKER] Synced authorized keys")
	CreateServerLog(work.Server, "info", "Successfully synced authorized keys", "")

	return nil
}

func serializePublicKey(k ssh.PublicKey) string {
	return base64.StdEncoding.EncodeToString(k.Marshal())
}

func parsePublicKey(hostKey string) (ssh.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(hostKey)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePublicKey(keyBytes)
}

func verifyHostKey(work *SyncServerWork) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		hostKeyString := work.Server.GetString("hostKey")
		if hostKeyString != "" {
			var err error
			hostKey, err := parsePublicKey(hostKeyString)
			if err != nil {
				CreateServerLog(work.Server, "error", "Invalid host key", err.Error())
				return err
			}

			err = ssh.FixedHostKey(hostKey)(hostname, remote, key)
			if err != nil {
				CreateServerLog(work.Server, "error", "Host key verification failed. Possible MITM.", err.Error())
				return err
			}

			return nil

		} else {
			// ask to verify host key
			msg := fmt.Sprintf("Host key is not known, do you want to trust it?\nFingerprint: %s\n", ssh.FingerprintSHA256(key))
			CreateServerLog(work.Server, "hostKey", msg, serializePublicKey(key))
			return fmt.Errorf("ssh: required host key is not known")
		}
	}
}

func executeSSHCommand(client *goph.Client, command string, args ...string) (string, error) {
	cmd, err := client.Command(command, args...)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}

func verifyServerHostname(client *goph.Client, hostname string) (string, error) {
	clientHostname, err := executeSSHCommand(client, "hostname", "-f")
	if err != nil {
		return clientHostname, fmt.Errorf("failed to execute `hostname -f`: %s", err.Error())
	}

	clientHostname = strings.TrimSpace(clientHostname)
	if clientHostname != hostname {
		return clientHostname, fmt.Errorf("hostname mismatch: %s != %s", clientHostname, hostname)
	}
	return clientHostname, nil
}

func collectAuthorizedKeys(dao *daos.Dao, server *models.Record) (authorizedKeys []AuthorizedKey, err error) {
	// collect direct users associated with this server
	userServersRecords, err := dao.FindRecordsByExpr("userServers", &dbx.HashExp{"serverId": server.Id})
	if err != nil {
		err = fmt.Errorf("failed to find userServers records: %s", err.Error())
		return
	}

	userIds := []interface{}{}
	userOptionsMap := map[string]string{}
	for _, usersServer := range userServersRecords {
		userIds = append(userIds, usersServer.GetString("userId"))
		userOptionsMap[usersServer.GetString("userId")] = usersServer.GetString("options")
	}

	// collect public keys associated with users
	publicKeys, err := dao.FindRecordsByExpr("publicKeys", dbx.In("userId", userIds...))
	if err != nil {
		err = fmt.Errorf("failed to find publicKeys records: %s", err.Error())
		return
	}

	authorizedKeys = []AuthorizedKey{}
	for _, publicKey := range publicKeys {
		key, err := parsePublicKey(publicKey.GetString("publicKey"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %s", err.Error())
		}

		authorizedKeys = append(authorizedKeys, AuthorizedKey{
			Options: userOptionsMap[publicKey.GetString("userId")],
			Type:    key.Type(),
			Key:     key,
			Comment: publicKey.GetString("comment"),
		})
	}

	return
}

func parseAuthorizedKeys(filename string) (startLines []string, betweenLines []string, endLines []string, err error) {
	startLines = []string{}
	betweenLines = []string{}
	endLines = []string{}

	// read fileContents
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	// read lines until start marker
	isWithinMarker := false
	lines := bytes.Split(fileContents, []byte("\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte(AUTHORIZED_KEYS_START_MARKER)) {
			isWithinMarker = true
			continue
		} else if bytes.Contains(line, []byte(AUTHORIZED_KEYS_END_MARKER)) {
			isWithinMarker = false
			continue
		}
		if isWithinMarker {
			betweenLines = append(betweenLines, string(line))
		} else if len(startLines) > 0 {
			endLines = append(endLines, string(line))
		} else {
			startLines = append(startLines, string(line))
		}
	}

	return
}

func assembleAuthorizedKeys(filename string, start []byte, between []AuthorizedKey, end []byte) error {
	// overwrite file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(start)
	file.WriteString("\n" + AUTHORIZED_KEYS_START_MARKER + "\n")
	file.WriteString("# WARNING: This section is auto-generated and is periodically replaced.\n")
	file.WriteString("# DO NOT MODIFY\n")
	for _, key := range between {
		file.WriteString(key.Serialize() + "\n")
	}
	file.WriteString(AUTHORIZED_KEYS_END_MARKER + "\n")
	file.Write(end)

	return nil
}
