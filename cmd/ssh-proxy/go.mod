module github.com/SUSE/eirini-ssh/cmd/ssh-proxy

go 1.12

require (
	code.cloudfoundry.org/bbs v0.0.0-20190724191824-01b96ad15a77
	code.cloudfoundry.org/cfhttp/v2 v2.0.0 // indirect
	code.cloudfoundry.org/clock v0.0.0-20180518195852-02e53af36e6c
	code.cloudfoundry.org/consuladapter v0.0.0-20190222031846-a0ec466a22b6
	code.cloudfoundry.org/debugserver v0.0.0-20180612203758-a3ba348dfede
	code.cloudfoundry.org/diego-logging-client v0.0.0-20190626151511-6278d4119f52
	code.cloudfoundry.org/diego-ssh v0.0.0-20190726165408-4e330a244ce1
	code.cloudfoundry.org/durationjson v0.0.0-20170223024715-718cecacb217
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible
	code.cloudfoundry.org/inigo v0.0.0-20190725181809-38aa015ae590
	code.cloudfoundry.org/lager v2.0.0+incompatible
	code.cloudfoundry.org/locket v0.0.0-20190524173003-285105ed8d9a
	code.cloudfoundry.org/tlsconfig v0.0.0-20190710180242-462f72de1106
	github.com/SUSE/eirini-ssh/proxy v0.0.0
	github.com/bmizerany/pat v0.0.0-20170815010413-6226ea591a40 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/hashicorp/consul/api v1.1.0
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/tedsuo/ifrit v0.0.0-20180802180643-bea94bb476cc
	github.com/tedsuo/rata v1.0.0 // indirect
	github.com/vito/go-sse v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
)

replace github.com/SUSE/eirini-ssh/proxy v0.0.0 => github.com/jimmykarily/eirini-ssh/proxy v0.0.0
