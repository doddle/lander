# Lander

A simple kubernetes (cluster) landing page

- search all the Ingresses in the current context cluster
- show any with the appropriate Annotations
- simple vuejs UI to query the golang server
- there is a lightweight cache system to ensure k8s API chatter is limited
- Identicons themed with the colours you decide (see: `-hex` when deploying lander)
- Unique "identicon" favicons generated against hostname (nicer+predictable bookmarks for multiple k8s clusters)


# Runtime flags

```sh
  -annotationBase string
    	The base of the annotations used for lander. e.g. lander.doddle.tech for annotations like lander.doddle.tech/show (default "lander.doddle.tech")
  -clusterFQDN string
    	The cluster this lander is operating in. Used for display and identicon purposes only. (default "k8s.example.com")
  -clusters string
    	comma seperated list of clusters (default "cluster1.example.com,cluster2.example.com")
  -color string
    	Main color scheme (See: https://vuetifyjs.com/en/styles/colors/#material-colors) (default "light-blue lighten-2")
  -debug
    	debug
  -hex string
    	identicon color, hex string, eg #112233, #123, #bAC (default "#26c5e8")
  -labels string
    	comma seperated list of node labels you care about (default "kubernetes.io/role,node.kubernetes.io/instance-type,node.kubernetes.io/instancegroup,topology.kubernetes.io/zone")
```

---
# local development

To run or develop lander follow the next steps:

## running the VueJS front-end

The frontend listens on `:8080` and will proxy traffic -> `localhost:8000` in development mode.

Simply:
```
cd frontend
npm i --from-lockfile
npm run serve
```

Good.. now that the front-end is running.. launch the backend:


## running the backend

The backend listens on `:8000` and will read assets from a relative path of: `./frontend/dist`

The golang app supports a few flags..:
- `-annotationBase` the base of the annotations to use. Appended with `/show` to figure the annotations above
- `-clusterFQDN` The cluster this lander is operating in. Used for display and identicon purposes only.
- `-hex` the hex colour to be used for the identicon.. EG: `#26c5e8`, `#123`, `#b2C`
- `-color` the colorscheme (vuetify) to be handed to the frontend.. See: [material-colors](https://vuetifyjs.com/en/styles/colors/#material-colors)

If you run this locally on your computer it will auto-detect your `$KUBECONFIG` and use the current context.
```
go run .
```

curl it and see if u get any endpoints detected
```
curl localhost:8000/v1/endpoints
```

NOTE: opening `:8000` in your web browser won't render anything until after you've run `npm run build` in the frontend. For local-development of the UI we recommend using the VueJS steps above.

## Adding a link

You'll need to add an annotation to an ingress to see it in the "Links" section

For example.. In the following example we add:
1. the annotation `lander.doddle.tech/show="true"`
2. to an ingress called `prometheus`
3. in the namespace `monitoring`

```sh
export NAMESPACE=monitoring
export INGRESS=prometheus
kubectl annotate \
  -n $NAMESPACE \
  ingress/$INGRESS \
  lander.doddle.tech/show="true" --overwrite
```

Landers cache may not make the result show up immediately, but if you restart the golang api or wait a minute you should be able to see a link now

# TODOS:

Nice-to-haves/ideas for the future:

- allow adding your own icons
- more documentation
- more tests
- as this will live behind oauth2proxy like many other components, we can determine their username/email.. greet them
- maybe pull some basic info from RBAC if we're able to determine their github groups (EG `hello <persons name>.. you have X namespaces you can play with on this cluster`)
- some information about cronjobs (or some CRDs potentially)
