### Negotiator 

Negotiator negotiates between an RHMAP core and a OpenShift instance. It understands both.
Negotiator plays the role of being the middleman between rhmap and OpenShift. It uses the OpenShift client and Kubernetes client to directly interact with the Kubernetes API and OpenShift API. For example when a new cloud app from RHMAP needs to be deployed or a new environment needs to be created, a request will be sent to negotiator with the details required such as a giturl, an auth token, env vars etc. Negotiator will take care of turning this into the required OpenShift Objects and sending them on to OpenShift / Kubernetes directly.

Try it out locally.

Create a new OpenShift Project ``` oc new-project mine ```

- ensure go 1.7 or greater installed
- install glide package manager ``` curl https://glide.sh/get | sh ``` 
- clone this repo into $GOPATH/src/github.com/feedhenry/negotiator
- ``` cd $GOPATH/src/github.com/feedhenry/negotiator ```
- ``` go install``` (will be slow the first time but faster afterwards)
- export some envars for config 
    - ``` export API_HOST="http://anopenshifthost.com:8443" ```
    - ``` export API_TOKEN="some-os-token" ```
    - ``` export DEPLOY_NAMESPACE="mine" ```



Start the server 

```
negotiator

```

Curl request 

```
curl -X POST  http://localhost:3000/deploy/cloudapp -d '{"guid":"testguid","namespace":"mine","domain":"testing"}'

```

Note currently it just sets up a service, a route and an imagestream.

## Developing
 - Layout

```bash
.
├── config
├── domain #domain specfic logic 
│   ├── openshift # our logic around openshift
│   └── rhmap  # our logic specific to rhmap
└── pkg
    └── openshift # pkg for making the openshift and kubernetes client more simple to work with. Our domain logic does not go here
##handlers go in the root dir and deal with http specific logic 
deployHandler.go 
sysHandler.go      

``` 

## Test 

```bash
make test-unit 
```

## build and publish

```
env GOOS=linux go build .

docker build -t rhmap/negotiator:0.0.1 . ##change build number

docker tag rhmap/negotiator:0.0.1 rhmap/negotiator:latest

docker push rhmap/negotiator:0.0.1

```

## Run in OpenShift

```
oc new-app -f os_template.json --param=API_TOKEN=<YOUR_TOKEN>,API_HOST=<OpenShift_HOST>,DEPLOY_NAMESPACE=<SOME_NS>

```
