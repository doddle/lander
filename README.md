# lander

Simple kubernetes landing page..

- Search all the Ingresses in the current context cluster
- Show any with the appropriate Annotations
- simple vuejs UI to query the golang server
- there is a lightweight cache system to ensure k8s API chatter is limited
- guess colour scheme based on hostname
- Unique "identicon" favicons generated against hostname (nicer+predictable bookmarks for multiple k8s clusters)
- Identicons themed with the colours you decide (see: `-hex` when deploying lander)

The available annotations for an Ingress (with `-annotationBase` startup flag of `lander.doddle.tech`)
 - `"lander.doddle.tech/show": "true"`
 Note string, not bool
 - `"lander.doddle.tech/name": "Kibana"`
 Service name, defaults to K8s Ingress service name
 - `"lander.doddle.tech/description": "A free text description"` Defaults to blank string
 - `"lander.doddle.tech/icon": "kibana.png"`
 A file in the ./assets directory. Defaults to matching a file with a lowercase("lander.doddle.tech/name" or service name) and fallback to `link.png` if none is found
 - `"lander.doddle.com/url": "https://doddle.com"` URL to link to. Defaults to K8s Ingress URL.


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
