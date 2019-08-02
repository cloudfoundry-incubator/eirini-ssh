# Eirini ssh

This is a component that enables `cf ssh` for [Eirini](https://github.com/cloudfoundry-incubator/eirini) CloudFoundry clusters.
It is the equivalent of [diego-ssh](https://github.com/cloudfoundry/diego-ssh) for Eirini.

## Test it out

- Compile:

```go build```

- Start it with a working config.json (check the example one under `cmd/ssh-proxy/config.json`) :

```./ssh-proxy -config config.json```

- Find your applications guid:

```
cf app myapp --guid
```

- Try to ssh to the proxy:

```
ssh -p 2222 cf:b654358e-edfd-4e9e-b646-7fa55d5f8eb7/0@127.0.0.1
```

- When promted for a password use the one from the command:

```
cf ssh-code
```


More information:

- https://github.com/cloudfoundry/diego-ssh#cloud-foundry-via-cloud-controller-and-uaa
- https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst#client-obtains-token-post-oauthtoken


## Notes

To create a correct config.json file you will need the ssh-proxy user's secret (password). You can find that with something like

```
kubectl get secrets -n scf secrets-2.17.1-1 -o yaml | less
```

(look for uaa-clients-diego-ssh-proxy-secret)

decode the value with `echo "the_value_here" | base64 -d`. This is the value you need for the config.json file.


