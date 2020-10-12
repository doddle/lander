

Simple little landing homepage..

- Indexes all the ingresses in the cluster
- returns json with some meta-data only for ones with the same hostname
- later to have a spiffy UI served from static contents in `./public`
- there is a lightweight cache system to ensure k8s API chatter is limited


# running it locally

The app will auto-detect your `$KUBECONFIG` or if its running inside k8s, so simply run it:

```
go run .
```

Curl the API and set an example hostname (it will filter all ingress records looking for ones "containing" the `Host: ` header..

So setting that to `.` will match all records..

EG:


```
curl localhost:8000/v1/endpoints -H "Host: ."
```

# TODOS:

- throw in some icons for common services such as prometheus+grafana etc and present a very simple list of links (initial purpose of this)
- Generate an Identicon for the server based on the hostname header
- search the hostname and set the identicon colour based on some globs.. etc `dev` `staging` `prod` should allow different (but predicible) colours
- map the identicon so we get a favicon into users browsers

- lightweight react/vue frontend (SPA) would rock.. we can actually extend this to return simple and useful information. EG prometheus metrics, service counts, warnings about ingresses without "alive" services (no k8s endpoints)

- as this will live behind oauth2proxy like many other components, we can determine their username/email.. greet them
- maybe pull some basic info from RBAC if we're able to determine their github groups (EG `hello <persons name>.. you have X namespaces you can play with on this cluster`)
