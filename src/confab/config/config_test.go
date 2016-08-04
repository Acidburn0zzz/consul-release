package config_test

import (
	"github.com/cloudfoundry-incubator/consul-release/src/confab/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("Default", func() {
		It("returns a default configuration", func() {
			Expect(config.Default()).To(Equal(config.Config{
				Consul: config.ConfigConsul{
					Agent: config.ConfigConsulAgent{
						Servers: config.ConfigConsulAgentServers{
							LAN: []string{},
							WAN: []string{},
						},
						DnsConfig: config.ConfigConsulAgentDnsConfig{
							AllowStale: false,
							MaxStale:   "5s",
						},
					},
				},
				Path: config.ConfigPath{
					AgentPath:       "/var/vcap/packages/consul/bin/consul",
					ConsulConfigDir: "/var/vcap/jobs/consul_agent/config",
					PIDFile:         "/var/vcap/sys/run/consul_agent/consul_agent.pid",
					KeyringFile:     "/var/vcap/store/consul_agent/serf/local.keyring",
					DataDir:         "/var/vcap/store/consul_agent",
				},
				Confab: config.ConfigConfab{
					TimeoutInSeconds: 55,
				},
			}))
		})
	})

	Describe("ConfigFromJSON", func() {
		It("returns a config given JSON", func() {
			json := []byte(`{
				"node": {
					"name": "nodename",
					"index": 1234,
					"external_ip": "10.0.0.1"
				},
				"path": {
					"agent_path": "/path/to/agent",
					"consul_config_dir": "/consul/config/dir",
					"pid_file": "/path/to/pidfile",
					"keyring_file": "/path/to/keyring",
					"data_dir": "/path/to/data/dir"
				},
				"consul": {
					"agent": {
						"services": {
							"myservice": {
								"name" : "myservicename"	
							}
						},
						"mode": "server",
						"datacenter": "dc1",
						"log_level": "debug",
						"protocol_version": 1,
						"servers": {
							"lan": ["server1", "server2", "server3"],
							"wan": ["wan-server1", "wan-server2", "wan-server3"]
						},
						"dns_config": {
							"allow_stale": true,
							"max_stale": "15s"
						}
					},
					"encrypt_keys": ["key-1", "key-2"]
				},
				"confab": {
					"timeout_in_seconds": 30
				}
			}`)

			cfg, err := config.ConfigFromJSON(json)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).To(Equal(config.Config{
				Path: config.ConfigPath{
					AgentPath:       "/path/to/agent",
					ConsulConfigDir: "/consul/config/dir",
					PIDFile:         "/path/to/pidfile",
					KeyringFile:     "/path/to/keyring",
					DataDir:         "/path/to/data/dir",
				},
				Node: config.ConfigNode{
					Name:       "nodename",
					Index:      1234,
					ExternalIP: "10.0.0.1",
				},
				Consul: config.ConfigConsul{
					Agent: config.ConfigConsulAgent{
						Services: map[string]config.ServiceDefinition{
							"myservice": {
								Name: "myservicename",
							},
						},
						Mode:            "server",
						Datacenter:      "dc1",
						LogLevel:        "debug",
						ProtocolVersion: 1,
						Servers: config.ConfigConsulAgentServers{
							LAN: []string{"server1", "server2", "server3"},
							WAN: []string{"wan-server1", "wan-server2", "wan-server3"},
						},
						DnsConfig: config.ConfigConsulAgentDnsConfig{
							AllowStale: true,
							MaxStale:   "15s",
						},
					},
					EncryptKeys: []string{"key-1", "key-2"},
				},
				Confab: config.ConfigConfab{
					TimeoutInSeconds: 30,
				},
			}))
		})

		It("returns a config with default values", func() {
			json := []byte(`{}`)
			cfg, err := config.ConfigFromJSON(json)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).To(Equal(config.Config{
				Path: config.ConfigPath{
					AgentPath:       "/var/vcap/packages/consul/bin/consul",
					ConsulConfigDir: "/var/vcap/jobs/consul_agent/config",
					PIDFile:         "/var/vcap/sys/run/consul_agent/consul_agent.pid",
					KeyringFile:     "/var/vcap/store/consul_agent/serf/local.keyring",
					DataDir:         "/var/vcap/store/consul_agent",
				},
				Consul: config.ConfigConsul{
					Agent: config.ConfigConsulAgent{
						Servers: config.ConfigConsulAgentServers{
							LAN: []string{},
							WAN: []string{},
						},
						DnsConfig: config.ConfigConsulAgentDnsConfig{
							AllowStale: false,
							MaxStale:   "5s",
						},
					},
				},
				Confab: config.ConfigConfab{
					TimeoutInSeconds: 55,
				},
			}))
		})

		It("returns an error on invalid json", func() {
			json := []byte(`{%%%{{}{}{{}{}{{}}}}}}}`)
			_, err := config.ConfigFromJSON(json)
			Expect(err).To(MatchError(ContainSubstring("invalid character")))
		})
	})
})
