# pypi

This repo was pretty much a clone of https://github.com/Akshaya-T/pypi.git from the tutorial found here:
https://python.plainenglish.io/private-pypi-server-on-kubernetes-7df169864972

I have written the resource definitions in Go though. After compiling the executable, it can be used
to generate the required yaml files, fully parameterized, at least as far as I felt like it.

## Starting minikube cluster

I think you need en1 when you are on a wired internet connection, but it's important to
use the actual local IP and not "localhost" or "127.0.0.0". I think, maybe restarting everything
with `minikube delete` was enough to get rid of a glitch I was experiencing.
And oh yeah, this `ipconfig` command works on Mac, not sure about other OSs.

```bash
minikube start --insecure-registry="$(ipconfig getifaddr en0):5000"
```

## Local docker registry

This creates a local docker registry where the Kubernetes pod can pull our Docker image from.
This can of course be a remote one as well (For example `jmeidam/pypiserver`).

```bash
docker run -d -p 5000:5000 --name registry registry:2
```

## Build and push the Docker image

If not using the local registry as above, replace "localhost:5000" with you registry

```bash
docker build -t localhost:5000/pypi-server .
```

```bash
docker push localhost:5000/pypi-server
```

## Build Go executable and generate yaml files

Ensure you have Go installed ofcourse.
Move into the `gok8s` folder and run `go build -o ./pypilocal main.go`.
This will create an executable `pypilocal` in that folder, which you can run from anywhere.

Generate the yaml files using the following command

```bash
./pypilocal -outputpath yamls -pypilogins '{"password": "123", "username": "pypi"}' -image 192.123.1.23:5000/pypi-server
```

Where indeed `192.123.1.23` should be replaced with your localhost IP and the username and password can also be changed to anything you like. Just make sure you save them somewhere to access the pypi service later on.


## Apply kubernetes files in generated yaml folder

Do this when you are sure you are in the right workspace and context. When you run `minikube start` it
appears to set the defaults to the minikube cluster:
*"Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default"*

```bash
# Make sure to change secret values before apply
kubectl apply -f gok8s/yamls
```

Where `gok8s/yamls` is the folder you specified when running `pypilocal`.


## Port Forward and see the PyPi UI

First find the pypi pod ID by running `kubectl get pods`. You'll want the name.

```bash
kubectl port-forward <<pod-name>> 8080:80 --address 0.0.0.0
```


## Upload Package pypiserver using Poetry

This sets the configuration for de server referred to as "local", you can give it any name.
```bash
poetry config repositories.local http://localhost:8080
poetry config http-basic.local pypi 123
```

Then run `poetry publish -r local`

Or run the following when working with Twine, after updating the `~/.pypirc` file as in the tutorial.

```bash
python3 -m twine upload -r local dist/
```

## Use pip client to install and verify the package

```bash
pip install -i http://pypi:pass@localhost:8080/simple yourpackage
```

Or indeed when working with Poetry add this to `pyproject.toml`

```toml
[[tool.poetry.source]]
name = "local"
url = "http://localhost:8080/simple"
secondary = true

# The usual dependencies section
[tool.poetry.dependencies]
yourpackage = { version = "^1.0.0", source = "local" }
```

The `poetry config http-basic.local` command above has already configured the username and password. No need to hardcode those into toml file.
