

develop | master
---|---
[![Build Status](https://ci.digtux.com/api/badges/digtux/lander/status.svg?ref=refs/heads/develop)](https://ci.digtux.com/digtux/lander) | [![Build Status](https://ci.digtux.com/api/badges/digtux/lander/status.svg?ref=refs/heads/master)](https://ci.digtux.com/digtux/lander)

Simple little landing homepage..

- Indexes all the ingresses in the cluster
- returns json with some meta-data only for ingress hostnames matching your `-host` setting
- simple vuejs UI to query the golang server
- there is a lightweight cache system to ensure k8s API chatter is limited
- guess colour scheme based on hostname
- "identicon" favicons generated against hostname


# running it locally


## back

The golang app supports a few flags..:
- `-debug` for extra logging
- `-host` which hostname to filter for results on the main page
- `-config` will allow you to read a `config.yaml` (see `config-example.yaml`)  **this allows you to set custom colour schemes**

If you run this locally on your computer it will auto-detect your `$KUBECONFIG` and use the current context.

For example, my k8s cluster's context is call: `mycluster.example.com`

and in this cluster I have:

- an ingress for prometheus: `mycluster.example.com/prom`
- an ingress for alertmanager: `mycluster.example.com/alerts`

This means I can run the app with:

```
go run . -host "$(kubectl config current-context)" -debug
```


if your "current-context" is something else.. simply over-ride it.. eg

```
go run . -host "ingress.tosee.com"
```


curl it and see if u get any endpoints detected
```
curl localhost:8000/v1/endpoints
```



## front

Once the backend is running.. you can now launch the frontend

The frontend listens on `:8080` and will proxy traffic -> `localhost:8000` as you develop

Simply:
```
cd frontend
npm i --from-lockfile
npm run serve
```

# TODOS:

- throw in some icons for common services such as prometheus+grafana etc and present a very simple list of links (initial purpose of this)
- Generate an Identicon for the server based on the hostname header
- search the hostname and set the identicon colour based on some globs.. etc `dev` `staging` `prod` should allow different (but predicible) colours
- map the identicon so we get a favicon into users browsers

- lightweight react/vue frontend (SPA) would rock.. we can actually extend this to return simple and useful information. EG prometheus metrics, service counts, warnings about ingresses without "alive" services (no k8s endpoints)

- as this will live behind oauth2proxy like many other components, we can determine their username/email.. greet them
- maybe pull some basic info from RBAC if we're able to determine their github groups (EG `hello <persons name>.. you have X namespaces you can play with on this cluster`)
