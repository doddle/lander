# Lander

A kubernetes (cluster) landing page


## Features/progress

- [x] golang powered interactings with k8s 
- [x] Simple server-side caching to keep kube-api calls limited.
- [x] List nodes
- [x] Configurable list of node labels to show table
      (Similar to `kubectl get nodes -L beta.kubernetes.io/arch -L beta.kubernetes.io/os -L beta.kubernetes.io/instance-type` you'd run lander with `-l ,beta.kubernetes.io/instance-type`)
- [x] VueJS based PWA/SPA
- [x] Show links to cluster ingress endpoints based on if it has lander annotations, add the annotation `lander.doddle.tech/show="true"` to the ingress objects.
- [x] Unique "identicon" favicons generated against hostname (nicer+predictable bookmarks for multiple k8s clusters)
- [x] Identicons themed with the colours you decide (see: `-hex` when deploying lander)
- [x] search all the Ingresses in the current context cluster
- [x] "at-a-glance" piechart for **Deployments**
- [x] "at-a-glance" piechart for **StatefulSets**
- [x] "at-a-glance" piechart for **Nodes**
- [x] "at-a-glance" piechart for **DaemonSets**
- [x] Basic node information (age, labels, `Ready` and `Unscheduable`)
- [x] front-end settings (colorscheme etc) configurable via launch options. We use different colours to hint at our kubernetes "lifecycles" (`dev`, `staging`, `prod`, etc..)
- [x] Links will indicate a lock icon if `nginx-ingress` and `oauth2proxy` protection was detected on ingress.
- [ ] Allow users to search all Ingress rules and see what hostnames/paths are mapped to which services
- [ ] provide information about the api-groups and supported api's detected
- [ ] provide basic cluster info such as the k8s version etc.. probably worth add a column for this to the `nodes` tab
- [ ] Represent to the user any objects detected in the cluster that have weave/flux ignore annotations (EG: potential state-drift from gitops definitions.)

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

> NOTE: before launching the backend, either `mkdir ./frontend/dist` or run `cd frontend ; npm run build`
>       needless to say if you make an empty directory, don't expect lander render.
>

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
