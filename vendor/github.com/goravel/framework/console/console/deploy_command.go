package console

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/file"
)

/*
DeployCommand
===============

Overview
--------
This command implements a simple, opinionated deployment pipeline for Goravel applications.
It builds the application locally, performs a one-time remote server setup, uploads the
required artifacts to the server, restarts a systemd service, and supports rollback to the
previous binary. The goal is to provide a pragmatic, single-command deploy for small-to-medium
workloads.


Architecture assumptions
------------------------
Two primary deployment topologies are supported:
1) Reverse proxy in front of the app (recommended)
   - reverseProxyEnabled=true
   - App listens on 127.0.0.1:<app.deploy.reverse_proxy_port> (e.g. 9000)
   - Caddy proxies public HTTP(S) traffic to the app
   - If reverseProxyTLSEnabled=true and a valid domain is configured, Caddy terminates TLS
     and automatically provisions certificates; otherwise Caddy serves plain HTTP on :80

2) No reverse proxy
   - reverseProxyEnabled=false
   - App listens directly on :80 (APP_HOST=0.0.0.0, APP_PORT=80)

Artifacts & layout on server
----------------------------
Remote base directory: /var/www/<APP_NAME>
Files managed by this command on the remote host:
  - main        : current binary (running)
  - backups/    : timestamped zip archives of previous states (used for rollback)
  - .env        : environment file (uploaded from app.deploy.prod_env_file_path)
  - public/     : optional static assets
  - storage/    : optional storage directory (uploaded only during setup)
  - resources/  : optional resources directory

Idempotency & first-time setup
------------------------------
The initial server setup is performed exactly once per server (per app name). The command first
checks if /etc/systemd/system/<APP_NAME>.service exists over SSH. If it exists, setup is skipped.
Otherwise, the command:
  - Installs and configures Caddy (only when reverseProxyEnabled=true)
  - Creates the app directory and sets ownership
  - Writes the systemd unit for <APP_NAME>
  - Enables the service and configures the firewall (ufw)

Subsequent deploys skip the setup entirely for speed and safety (unless --force-setup is used).
Note: If you change proxy/TLS/domain settings later, pass --force-setup to re-apply provisioning
changes (e.g., regenerate Caddyfile, adjust firewall rules, rewrite the unit file).

Rollback model
--------------
Every deployment creates a timestamped zip archive of the current state under backups/.
Rollback restores from the latest archive and restarts the service.

Build & artifacts (local)
-------------------------
The command builds the binary (name: APP_NAME) using the configured target OS/ARCH and static
linking preference. See Goravel docs for compiling guidance, artifacts, and what to upload:
https://www.goravel.dev/getting-started/compile.html

Configuration (app.config)
--------------------------
This command reads from application configuration (see app/config/app.go), not directly from .env.
Required:
  - app.name                             : Application name (used in remote paths/service name)
  - app.deploy.ssh_ip                    : Target server IP
  - app.deploy.reverse_proxy_port        : Backend app port when reverse proxy is used (e.g. 9000)
  - app.deploy.ssh_port                  : SSH port (e.g. 22)
  - app.deploy.ssh_user                  : SSH username (user must have sudo privileges)
  - app.deploy.ssh_key_path              : Path to SSH private key (e.g. ~/.ssh/id_rsa)
  - app.build.os                         : Target OS for build (e.g. linux)
  - app.build.arch                       : Target arch for build (e.g. amd64)
  - app.deploy.prod_env_file_path        : Local path to production .env file to upload

Optional / boolean flags (default false if unset):
  - app.build.static                     : Build statically when true
  - app.deploy.reverse_proxy_enabled     : Use Caddy reverse proxy when true
  - app.deploy.reverse_proxy_tls_enabled : Enable TLS (requires domain) when true
  - app.deploy.domain                    : Domain name for TLS or HTTP vhost when using Caddy
                                           (required only if TLS is enabled)

CLI flags
---------
  - --only                                : Comma-separated subset to deploy: main,env,public,storage,resources
  - -r, --rollback                        : Rollback to previous binary
  - -f, --force-setup                     : Force re-run of provisioning even if already set up

Security & firewall
-------------------
The command uses SSH with StrictHostKeyChecking=no for convenience. For production, consider
manually trusting the host key to avoid MITM risks. Firewall rules are applied via ufw with
safe ordering: allow OpenSSH and required HTTP(S) ports first, then enable ufw to avoid losing
SSH connectivity.

Systemd service
---------------
The unit runs under app.deploy.ssh_user. Environment variables are provided via the unit for host/port,
and the working directory points to /var/www/<APP_NAME>. Service restarts are used (brief downtime).
For zero-downtime swaps, a more advanced process manager or socket activation would be required.

High-level deployment flow
--------------------------
1) Build: compile the binary for the specified target (OS/ARCH, static optional) with name APP_NAME
2) Determine artifacts to upload: main, .env, public, storage (setup only), resources (filter via --only)
3) Setup (first deploy only, or when --force-setup):
   - Create directories and permissions
   - Install/configure Caddy based on reverse proxy + TLS settings
   - Write systemd unit and enable service
   - Configure ufw rules (OpenSSH, 80, and 443 as needed)
4) Upload:
   - Binary: upload to main.new, atomically move main.new to main
   - .env:   upload to .env.new, atomically move to .env
   - public, storage, resources: recursively upload if they exist locally
5) Restart service: systemctl daemon-reload, then restart (or start) the service

Known limitations
-----------------
  - No migrations or database orchestration
  - Rollback covers only the binary; assets/env are not rolled back
  - StrictHostKeyChecking is disabled by default for convenience
  - Changing proxy/TLS/domain requires --force-setup to re-apply provisioning
  - Assumes Debian/Ubuntu with apt-get and ufw available

Usage examples
--------------

Usage example (1 - with reverse proxy):

Assuming you have the following .env file stored in the root of your project as .env.production:
```
APP_NAME=my-app
DEPLOY_SSH_IP=127.0.0.1
DEPLOY_REVERSE_PROXY_PORT=9000
DEPLOY_SSH_PORT=22
DEPLOY_SSH_USER=deploy
DEPLOY_SSH_KEY_PATH=~/.ssh/id_rsa
DEPLOY_OS=linux
DEPLOY_ARCH=amd64
DEPLOY_PROD_ENV_FILE_PATH=.env.production
DEPLOY_STATIC=true
DEPLOY_REVERSE_PROXY_ENABLED=true
DEPLOY_REVERSE_PROXY_TLS_ENABLED=true
DEPLOY_DOMAIN=my-app.com
```
You can then deploy your application to the server with the following command:
```
go run . artisan deploy
```
This will:
1. Build the application
2. On the remote server: install Caddy as a reverse proxy, support TLS, configure Caddy to proxy traffic to the application on port 9000, and only allow traffic from the domain my-app.com.
3. On the remote server: install ufw, and set up the firewall to allow traffic to the application.
4. On the remote server: create the systemd unit file and enable it
5. Upload the application binary, environment file, public directory, storage directory, and resources directory to the server
6. Restart the systemd service that manages the application


Usage example (2 - without reverse proxy):

You can also deploy without a reverse proxy by setting the DEPLOY_REVERSE_PROXY_ENABLED environment variable to false. For example,
assuming you have the following .env file stored in the root of your project as .env.production and you want to deploy your application to the server without a reverse proxy:
```
APP_NAME=my-app
DEPLOY_SSH_IP=127.0.0.1
DEPLOY_REVERSE_PROXY_PORT=80
DEPLOY_SSH_PORT=22
DEPLOY_SSH_USER=deploy
DEPLOY_SSH_KEY_PATH=~/.ssh/id_rsa
DEPLOY_OS=linux
DEPLOY_ARCH=amd64
DEPLOY_PROD_ENV_FILE_PATH=.env.production
DEPLOY_STATIC=true
DEPLOY_REVERSE_PROXY_ENABLED=false
DEPLOY_REVERSE_PROXY_TLS_ENABLED=false
DEPLOY_DOMAIN=
```

You can then deploy your application to the server with the following command:
```
go run . artisan deploy
```

This will:
1. Build the application
2. On the remote server: install ufw, and set up the firewall to allow traffic to the application that is listening on port 80 (http).
3. On the remote server: create the systemd unit file and enable it
4. Upload the application binary, environment file, public directory, storage directory, and resources directory to the server
5. Restart the systemd service that manages the application
```

Usage example (3 - rollback):

You can also rollback a deployment to the previous binary by running the following command:
```
go run . artisan deploy --rollback
```


Usage example (4 - force setup):

You can also force the setup of the server by running the following command:
```
go run . artisan deploy --force-setup
```


Usage example (5 - only deploy subset of files):

You can also deploy only a subset of the files (such as only the main binary and the environment file) by running the following command:
```
go run . artisan deploy --only main,env
```
*/

// deployOptions is a struct that contains all the options for the deploy command
type deployOptions struct {
	appName                string
	arch                   string
	deployBaseDir          string
	domain                 string
	envDecryptKey          string
	httpPort               string
	prodEnvFilePath        string
	reverseProxyEnabled    bool
	reverseProxyPort       string
	reverseProxyTLSEnabled bool
	remoteEnvDecrypt       bool
	sshIp                  string
	sshKeyPath             string
	sshPort                string
	sshUser                string
	staticEnv              bool
	targetOS               string
}

type uploadOptions struct {
	hasMain      bool
	hasProdEnv   bool
	hasPublic    bool
	hasStorage   bool
	hasResources bool
}

type DeployCommand struct {
	artisan console.Artisan
	config  config.Config
	process process.Process
}

func NewDeployCommand(artisan console.Artisan, config config.Config, process process.Process) *DeployCommand {
	return &DeployCommand{
		artisan: artisan,
		config:  config,
		process: process,
	}
}

// Signature The name and signature of the console command.
func (r *DeployCommand) Signature() string {
	return "deploy"
}

// Description The console command description.
func (r *DeployCommand) Description() string {
	return "Deploy the application"
}

// Extend The console command extend.
func (r *DeployCommand) Extend() command.Extend {
	return command.Extend{
		Flags: []command.Flag{
			&command.StringFlag{
				Name:  "only",
				Usage: "Comma-separated subset to deploy: main,public,storage,resources,env. For example, to only deploy the main binary and the environment file, you can use 'main,env'",
			},
			&command.BoolFlag{
				Name:               "rollback",
				Aliases:            []string{"r"},
				Value:              false,
				Usage:              "Rollback to previous deployment",
				DisableDefaultText: true,
			},
			&command.BoolFlag{
				Name:               "force-setup",
				Aliases:            []string{"force"},
				Value:              false,
				Usage:              "Force re-run server setup even if already configured",
				DisableDefaultText: true,
			},
		},
	}
}

// Handle Execute the console command.
func (r *DeployCommand) Handle(ctx console.Context) error {
	// Rollback check first: allow rollback without validating local host tools
	// (tests can short-circuit Spinner; real runs will still use ssh remotely)
	if ctx.OptionBool("rollback") {
		opts, err := r.getDeployOptions(ctx)
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}
		if res := r.process.WithSpinner("Rolling back...").Run(rollbackCommand(opts)); res.Failed() {
			ctx.Error(res.Error().Error())
			return nil
		}

		ctx.Info("Rollback successful.")
		return nil
	}

	// check if the local host is valid, requires scp, ssh, and bash to be installed and in your path.
	if err := validLocalHost(); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	// get all options
	opts, err := r.getDeployOptions(ctx)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	// continue normal deploy flow

	// Step 1: build the application by invoking the build command via Artisan (no shell exec)
	buildCmd := fmt.Sprintf("build --os %s --arch %s --name %s", opts.targetOS, opts.arch, opts.appName)
	if opts.staticEnv {
		buildCmd += " --static"
	}
	if err = r.artisan.Call(buildCmd); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	// Step 2: verify which files to upload (main, env, public, storage, resources)
	upload := getUploadOptions(ctx, opts.appName, opts.prodEnvFilePath)

	// If the production env file is encrypted (per Goravel docs), decrypt it first (locally or remotely)
	envPathToUpload := opts.prodEnvFilePath
	remoteDecrypt := false
	remoteEncName := ""
	if upload.hasProdEnv {
		// Detect encrypted env by content (base64 + AES block structure with IV)
		if data, readErr := os.ReadFile(opts.prodEnvFilePath); readErr == nil {
			if isEncryptedEnvContent(data) {
				if opts.remoteEnvDecrypt {
					remoteDecrypt = true
					remoteEncName = filepath.Base(opts.prodEnvFilePath)
					envPathToUpload = opts.prodEnvFilePath
				} else {
					cmd := fmt.Sprintf("env:decrypt --name %q", opts.prodEnvFilePath)
					if strings.TrimSpace(opts.envDecryptKey) != "" {
						cmd += fmt.Sprintf(" --key %q", opts.envDecryptKey)
					}
					if err = r.artisan.Call(cmd); err != nil {
						ctx.Error(err.Error())
						return nil
					}
					// env:decrypt writes to .env in the working directory
					envPathToUpload = ".env"
				}
			}
		}
	}

	// Step 3: set up server on first run —- skip if already set up unless --force-setup is used
	forceSetup := ctx.OptionBool("force-setup")
	setupNeeded := forceSetup || !r.isServerAlreadySetup(opts)
	if setupNeeded {
		if forceSetup {
			if res := r.process.WithSpinner("Removing previous server configuration...").Run(teardownServerCommand(opts)); res.Failed() {
				ctx.Error(res.Error().Error())
				return nil
			}
		}

		if res := r.process.WithSpinner("Setting up server (first time only)...").Run(setupServerCommand(opts)); res.Failed() {
			ctx.Error(res.Error().Error())
			return nil
		}
	} else {
		ctx.Info("Server already set up. Skipping setup.")
	}

	// Enforce: storage can only be uploaded during setup stage
	if !setupNeeded {
		upload.hasStorage = false
	}

	// Step 4: upload files
	if res := r.process.WithSpinner("Uploading files...").Run(uploadFilesCommand(opts, upload, envPathToUpload, remoteDecrypt, remoteEncName)); res.Failed() {
		ctx.Error(res.Error().Error())
		return nil
	}

	// Optional: decrypt env remotely after upload
	if remoteDecrypt && upload.hasProdEnv {
		baseDir := opts.deployBaseDir
		if !strings.HasSuffix(baseDir, "/") {
			baseDir += "/"
		}
		appDir := fmt.Sprintf("%s%s", baseDir, opts.appName)
		decryptCmd := fmt.Sprintf("cd %s && ./main artisan env:decrypt --name %q", appDir, remoteEncName)
		if strings.TrimSpace(opts.envDecryptKey) != "" {
			decryptCmd += fmt.Sprintf(" --key %q", opts.envDecryptKey)
		}
		script := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s '%s'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, decryptCmd)

		if res := r.process.WithSpinner("Decrypting environment on remote...").Run(script); res.Failed() {
			ctx.Error(res.Error().Error())
			return nil
		}
	}

	// Step 5: restart service
	if res := r.process.WithSpinner("Restarting service...").Run(restartServiceCommand(opts)); res.Failed() {
		ctx.Error(res.Error().Error())
		return nil
	}

	ctx.Info("Deploy successful.")

	return nil
}

func (r *DeployCommand) getDeployOptions(ctx console.Context) (deployOptions, error) {
	opts := deployOptions{}
	opts.appName = r.config.GetString("app.name")
	opts.sshIp = r.config.GetString("app.deploy.ssh_ip")
	// Preferred: use HTTP server port (APP_PORT via http.port)
	opts.httpPort = r.config.GetString("http.port")
	// Back-compat: allow explicit reverse proxy backend port if provided, else fall back to http.port
	opts.reverseProxyPort = r.config.GetString("app.deploy.reverse_proxy_port")
	opts.sshPort = r.config.GetString("app.deploy.ssh_port")
	opts.sshUser = r.config.GetString("app.deploy.ssh_user")
	opts.sshKeyPath = r.config.GetString("app.deploy.ssh_key_path")
	opts.targetOS = r.config.GetString("app.build.os")
	opts.arch = r.config.GetString("app.build.arch")
	opts.domain = r.config.GetString("app.deploy.domain")
	opts.prodEnvFilePath = r.config.GetString("app.deploy.prod_env_file_path")
	opts.deployBaseDir = r.config.GetString("app.deploy.base_dir", "/var/www/")
	opts.envDecryptKey = r.config.GetString("app.deploy.env_decrypt_key")

	opts.staticEnv = r.config.GetBool("app.build.static")
	opts.reverseProxyEnabled = r.config.GetBool("app.deploy.reverse_proxy_enabled")
	opts.reverseProxyTLSEnabled = r.config.GetBool("app.deploy.reverse_proxy_tls_enabled")
	opts.remoteEnvDecrypt = r.config.GetBool("app.deploy.remote_env_decrypt")

	// Validate required options and report all missing at once
	var missing []string
	if opts.appName == "" {
		missing = append(missing, "APP_NAME")
	}
	if opts.sshIp == "" {
		missing = append(missing, "DEPLOY_SSH_IP")
	}
	if opts.reverseProxyPort == "" {
		missing = append(missing, "DEPLOY_REVERSE_PROXY_PORT")
	}
	if opts.sshPort == "" {
		missing = append(missing, "DEPLOY_SSH_PORT")
	}
	if opts.sshUser == "" {
		missing = append(missing, "DEPLOY_SSH_USER")
	}
	if opts.sshKeyPath == "" {
		missing = append(missing, "DEPLOY_SSH_KEY_PATH")
	}
	if opts.targetOS == "" {
		missing = append(missing, "DEPLOY_OS")
	}
	if opts.arch == "" {
		missing = append(missing, "DEPLOY_ARCH")
	}
	// domain is only required if reverse proxy TLS is enabled
	if opts.reverseProxyEnabled && opts.reverseProxyTLSEnabled && opts.domain == "" {
		missing = append(missing, "DEPLOY_DOMAIN")
	}
	if opts.prodEnvFilePath == "" {
		missing = append(missing, "DEPLOY_PROD_ENV_FILE_PATH")
	}
	if len(missing) > 0 {
		return deployOptions{}, fmt.Errorf("missing required environment variables: %s. Please set them in the .env file. Deployment cancelled", strings.Join(missing, ", "))
	}

	// expand ssh key ~ path if needed
	if after, ok := strings.CutPrefix(opts.sshKeyPath, "~"); ok {
		if home, herr := os.UserHomeDir(); herr == nil {
			opts.sshKeyPath = filepath.Join(home, after)
		}
	}

	return opts, nil
}

// isServerAlreadySetup checks if the systemd unit already exists on remote host
func (r *DeployCommand) isServerAlreadySetup(opts deployOptions) bool {
	checkCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'test -f /etc/systemd/system/%s.service'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, opts.appName)
	if res := r.process.Run(checkCmd); res.Failed() {
		return false
	}

	return true
}

// isEncryptedEnvContent determines whether the provided bytes likely represent an encrypted env
// according to Goravel's env:encrypt format (base64 of IV||ciphertext using AES with 16-byte IV).
// Heuristic:
//   - Base64-decodable
//   - Decoded length >= aes.BlockSize*2 (IV + at least one block)
//   - Decoded length % aes.BlockSize == 0
func isEncryptedEnvContent(raw []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return false
	}
	if len(decoded) < aes.BlockSize*2 {
		return false
	}
	if len(decoded)%aes.BlockSize != 0 {
		return false
	}
	return true
}

func getUploadOptions(ctx console.Context, appName, prodEnvFilePath string) uploadOptions {
	res := uploadOptions{}
	res.hasMain = file.Exists(appName)
	res.hasProdEnv = file.Exists(prodEnvFilePath)
	res.hasPublic = file.Exists("public")
	res.hasStorage = file.Exists("storage")
	res.hasResources = file.Exists("resources")

	// Allow subset selection via --only
	only := strings.TrimSpace(ctx.Option("only"))
	if only != "" {
		parts := strings.Split(only, ",")
		include := map[string]bool{}
		for _, p := range parts {
			include[strings.TrimSpace(strings.ToLower(p))] = true
		}
		if !include["main"] {
			res.hasMain = false
		}
		if !include["env"] {
			res.hasProdEnv = false
		}
		if !include["public"] {
			res.hasPublic = false
		}
		if !include["storage"] {
			res.hasStorage = false
		}
		if !include["resources"] {
			res.hasResources = false
		}
	}
	return res
}

// validLocalHost checks if the local host is valid, requires scp, ssh, and bash to be installed and in your path.
func validLocalHost() error {

	missingBins := []string{}
	if _, err := exec.LookPath("scp"); err != nil {
		missingBins = append(missingBins, "scp")
	}
	if _, err := exec.LookPath("ssh"); err != nil {
		missingBins = append(missingBins, "ssh")
	}
	// Shell requirements depend on OS
	if env.IsWindows() {
		if _, err := exec.LookPath("cmd"); err != nil {
			missingBins = append(missingBins, "cmd")
		}
	} else {
		if _, err := exec.LookPath("bash"); err != nil {
			missingBins = append(missingBins, "bash")
		}
	}

	if len(missingBins) > 0 {
		return fmt.Errorf("environment validation errors:\n - the following binaries were not found on your path: %s\n - Please install them, add them to your path, and try again", strings.Join(missingBins, ", "))
	}

	return nil
}

// setupServerCommand ensures Caddy and a systemd service are installed; no-op on subsequent runs
func setupServerCommand(opts deployOptions) string {
	// Directories and service
	baseDir := opts.deployBaseDir
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	appDir := fmt.Sprintf("%s%s", baseDir, opts.appName)
	binCurrent := fmt.Sprintf("%s/main", appDir)

	// Build systemd unit based on whether reverse proxy is used
	listenHost := "127.0.0.1"
	// If reverse proxy is enabled, app should listen on http.port (APP_PORT). If not set, fallback to reverseProxyPort for BC.
	appPort := opts.httpPort
	if strings.TrimSpace(appPort) == "" {
		appPort = opts.reverseProxyPort
	}
	if !opts.reverseProxyEnabled {
		// App listens on port 80 directly
		appPort = "80"
		listenHost = "0.0.0.0"
	}

	unit := fmt.Sprintf(`[Unit]
Description=Goravel App %s
After=network.target

[Service]
User=%s
WorkingDirectory=%s
ExecStart=%s
Environment=APP_HOST=%s
Environment=APP_PORT=%s
Restart=always
RestartSec=5
KillSignal=SIGINT
SyslogIdentifier=%s

[Install]
WantedBy=multi-user.target
`, opts.appName, opts.sshUser, appDir, binCurrent, listenHost, appPort, opts.appName)

	// Build Caddyfile if reverse proxy enabled
	caddyfile := ""
	if opts.reverseProxyEnabled {
		site := ":80"
		if strings.TrimSpace(opts.domain) != "" {
			site = opts.domain
		}
		upstream := fmt.Sprintf("127.0.0.1:%s", appPort)
		tlsLine := ""
		if !opts.reverseProxyTLSEnabled {
			tlsLine = "    tls off\n"
		}
		caddyfile = fmt.Sprintf(`%s {
    reverse_proxy %s {
        lb_try_duration 30s
        lb_try_interval 250ms
    }
    encode gzip
%s}
`, site, upstream, tlsLine)
	}

	unitB64 := base64.StdEncoding.EncodeToString([]byte(unit))
	var caddyB64 string
	if caddyfile != "" {
		caddyB64 = base64.StdEncoding.EncodeToString([]byte(caddyfile))
	}

	// Firewall commands based on configuration
	ufwCmds := []string{"sudo apt-get update -y && sudo apt-get install -y ufw", "sudo ufw --force enable"}
	if opts.reverseProxyEnabled {
		ufwCmds = append(ufwCmds, "sudo ufw allow 80")
		if opts.reverseProxyTLSEnabled {
			ufwCmds = append(ufwCmds, "sudo ufw allow 443")
		}
	} else {
		// App listens on 80 directly
		ufwCmds = append(ufwCmds, "sudo ufw allow 80")
	}

	// Remote setup script: create directories, install Caddy optionally, write configs
	script := fmt.Sprintf(`ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s '
set -e
if [ ! -d %s ]; then
  sudo mkdir -p %s
  sudo chown -R %s:%s %s
fi
%s
if [ ! -f /etc/systemd/system/%s.service ]; then
  echo %q | base64 -d | sudo tee /etc/systemd/system/%s.service >/dev/null
  sudo systemctl daemon-reload
  sudo systemctl enable %s
fi
%s
%s'
`, opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp,
		appDir, appDir, opts.sshUser, opts.sshUser, appDir,
		// caddy install and config
		func() string {
			if !opts.reverseProxyEnabled {
				return ""
			}
			install := "sudo apt-get update -y && sudo apt-get install -y caddy"
			writeCfg := fmt.Sprintf("echo %q | base64 -d | sudo tee /etc/caddy/Caddyfile >/dev/null && sudo systemctl enable --now caddy && sudo systemctl reload caddy || sudo systemctl restart caddy", caddyB64)
			return install + " && " + writeCfg
		}(),
		opts.appName, unitB64, opts.appName, opts.appName,
		// Firewall: open before enabling to avoid SSH lockout
		func() string {
			cmds := append([]string{"sudo ufw allow OpenSSH"}, ufwCmds...)
			return strings.Join(cmds, " && ")
		}(),
		"true",
	)

	return script
}

// teardownServerCommand removes prior Caddy and systemd service configuration to allow re-provisioning
func teardownServerCommand(opts deployOptions) string {
	baseDir := opts.deployBaseDir
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	appDir := fmt.Sprintf("%s%s", baseDir, opts.appName)
	// Remove Caddyfile and disable/stop Caddy if present; remove service unit and disable/stop service
	script := fmt.Sprintf(`ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s '
set -e
# Remove Caddy config if exists
if [ -f /etc/caddy/Caddyfile ]; then
  sudo rm -f /etc/caddy/Caddyfile || true
  sudo systemctl reload caddy || sudo systemctl restart caddy || true
fi
# Remove systemd unit if exists
if [ -f /etc/systemd/system/%s.service ]; then
  sudo systemctl stop %s || true
  sudo systemctl disable %s || true
  sudo rm -f /etc/systemd/system/%s.service || true
  sudo systemctl daemon-reload
fi
# Ensure app directory exists and permissions are consistent
sudo mkdir -p %s
sudo chown -R %s:%s %s
'`, opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, opts.appName, opts.appName, opts.appName, opts.appName, appDir, opts.sshUser, opts.sshUser, appDir)
	return script
}

// uploadFilesCommand uploads available artifacts to remote server
func uploadFilesCommand(opts deployOptions, up uploadOptions, envPathToUpload string, remoteDecrypt bool, remoteEncName string) string {
	baseDir := opts.deployBaseDir
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	appDir := fmt.Sprintf("%s%s", baseDir, opts.appName)
	remoteBase := fmt.Sprintf("%s@%s:%s", opts.sshUser, opts.sshIp, appDir)
	// ensure remote base exists and permissions
	cmds := []string{
		fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'sudo mkdir -p %s && sudo chown -R %s:%s %s'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, opts.sshUser, opts.sshUser, appDir),
	}

	// Create a timestamped backup zip of existing deploy artifacts before replacing any of them
	// Backup includes: main, .env, public, storage, resources (if present)
	backupCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'set -e; APP_DIR=%q; BACKUP_DIR=\"$APP_DIR/backups\"; TS=\"$(date +%%Y%%m%%d%%H%%M%%S)\"; sudo mkdir -p \"$BACKUP_DIR\"; if ! command -v zip >/dev/null 2>&1; then sudo apt-get update -y && sudo apt-get install -y zip; fi; cd \"$APP_DIR\"; INCLUDE=\"\"; [ -f main ] && INCLUDE=\"$INCLUDE main\"; [ -f .env ] && INCLUDE=\"$INCLUDE .env\"; [ -d public ] && INCLUDE=\"$INCLUDE public\"; [ -d storage ] && INCLUDE=\"$INCLUDE storage\"; [ -d resources ] && INCLUDE=\"$INCLUDE resources\"; if [ -n \"$INCLUDE\" ]; then zip -r \"$BACKUP_DIR/$TS.zip\" $INCLUDE >/dev/null; fi'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir)
	cmds = append(cmds, backupCmd)

	// main binary
	if up.hasMain {
		// upload to temp and atomically move
		cmds = append(cmds,
			fmt.Sprintf("scp -o StrictHostKeyChecking=no -i %q -P %s %q %s/main.new", opts.sshKeyPath, opts.sshPort, filepath.Clean(opts.appName), remoteBase),
			fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'sudo mv %s/main.new %s/main && sudo chmod +x %s/main'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, appDir, appDir),
		)
	}

	if up.hasProdEnv {
		// Upload env to a temp path, then atomically place as .env
		if remoteDecrypt {
			// upload encrypted file as provided path name
			destName := filepath.Base(envPathToUpload)
			cmds = append(cmds,
				fmt.Sprintf("scp -o StrictHostKeyChecking=no -i %q -P %s %q %s/%s", opts.sshKeyPath, opts.sshPort, filepath.Clean(envPathToUpload), remoteBase, destName),
			)
		} else {
			cmds = append(cmds,
				fmt.Sprintf("scp -o StrictHostKeyChecking=no -i %q -P %s %q %s/.env.new", opts.sshKeyPath, opts.sshPort, filepath.Clean(envPathToUpload), remoteBase),
				fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'sudo mv %s/.env.new %s/.env'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, appDir),
			)
		}
	}
	if up.hasPublic {
		cmds = append(cmds,
			fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'if [ -d %s/public ]; then sudo rm -rf %s/public; fi'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, appDir),
			fmt.Sprintf("scp -o StrictHostKeyChecking=no -i %q -P %s -r %q %s", opts.sshKeyPath, opts.sshPort, filepath.Clean("public"), remoteBase),
		)
	}
	if up.hasStorage {
		cmds = append(cmds,
			fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'if [ -d %s/storage ]; then sudo rm -rf %s/storage; fi'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, appDir),
			fmt.Sprintf("scp -o StrictHostKeyChecking=no -i %q -P %s -r %q %s", opts.sshKeyPath, opts.sshPort, filepath.Clean("storage"), remoteBase),
		)
	}
	if up.hasResources {
		cmds = append(cmds,
			fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'if [ -d %s/resources ]; then sudo rm -rf %s/resources; fi'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, appDir),
			fmt.Sprintf("scp -o StrictHostKeyChecking=no -i %q -P %s -r %q %s", opts.sshKeyPath, opts.sshPort, filepath.Clean("resources"), remoteBase),
		)
	}

	return strings.Join(cmds, " && ")
}

func restartServiceCommand(opts deployOptions) string {
	return fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s 'sudo systemctl daemon-reload && sudo systemctl restart %s || sudo systemctl start %s'", opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, opts.appName, opts.appName)
}

// rollbackCommand swaps main and main.prev if available, then restarts the service
func rollbackCommand(opts deployOptions) string {
	baseDir := opts.deployBaseDir
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	appDir := fmt.Sprintf("%s%s", baseDir, opts.appName)
	script := fmt.Sprintf(`ssh -o StrictHostKeyChecking=no -i %q -p %s %s@%s '
set -e
APP_DIR=%q
SERVICE=%q
BACKUP_DIR="$APP_DIR/backups"

# Ensure we have at least one backup to roll back to
TARGET_ZIP="$(ls -1t "$BACKUP_DIR"/*.zip 2>/dev/null | head -n1)"
if [ -z "$TARGET_ZIP" ]; then
  echo "No previous deployment backup to rollback to." >&2
  exit 1
fi

# Backup current state before rollback (so we can roll forward if needed)
TS="$(date +%%Y%%m%%d%%H%%M%%S)"
mkdir -p "$BACKUP_DIR"
if ! command -v zip >/dev/null 2>&1; then sudo apt-get update -y && sudo apt-get install -y zip; fi
cd "$APP_DIR"
INCLUDE=""
[ -f main ] && INCLUDE="$INCLUDE main"
[ -f .env ] && INCLUDE="$INCLUDE .env"
[ -d public ] && INCLUDE="$INCLUDE public"
[ -d storage ] && INCLUDE="$INCLUDE storage"
[ -d resources ] && INCLUDE="$INCLUDE resources"
if [ -n "$INCLUDE" ]; then zip -r "$BACKUP_DIR/rollback-$TS.zip" $INCLUDE >/dev/null; fi

# Restore from latest backup
if ! command -v unzip >/dev/null 2>&1; then sudo apt-get update -y && sudo apt-get install -y unzip; fi
unzip -o "$TARGET_ZIP" -d "$APP_DIR" >/dev/null

# Cleanup any *.newcurrent artifacts from previous failed operations (move this to the end as suggested)
find "$APP_DIR" -maxdepth 1 -name "*.newcurrent" -type f -exec sudo rm -f {} + || true

sudo systemctl daemon-reload
sudo systemctl restart "$SERVICE" || sudo systemctl start "$SERVICE"
 '`, opts.sshKeyPath, opts.sshPort, opts.sshUser, opts.sshIp, appDir, opts.appName)
	return script
}
