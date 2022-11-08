## GobblerD - Blood Bowl 2 replay parser

*All commands - unless specified otherwise - should be ran from the project root*

### Initializing the local CockroachDB cluster

**This only needs to be done if the cluster hasn't been initialized yet or the volumes have been recreated**

Start the Cockroach nodes:

```
$ docker compose up roach1 roach2 roach3
```

Once the nodes are running they'll warn about being unable to communicate with other nodes in the cluster.
To fix this, we'll have to initialize the cluster on the first node. In a new tab/window run the following:

```
$ docker compose exec roach1 ./cockroach init --insecure
```

Once this finishes, the nodes should find each other and they each report their own info in the original docker compose output.

Now the nodes can communicate with each other we should create the schema:

First we connect to the server running on the roach1 node

```
$ docker compose exec roach1 ./cockroach sql --insecure
```

Then we create the database and the single table we'll have for now:

```
> CREATE DATABASE gobb_dev;
> USE gobb_dev;
> CREATE TABLE replays (
      id uuid NOT NULL,
      competition jsonb NOT NULL,
      home_team jsonb NOT NULL,
      away_team jsonb NOT NULL
  );
```

When this is done, quit the cockroach client (Ctrl+D) and shut the cluster down by pressing Ctrl+C in the docker compose tab/window

### Starting the daemon

Provided the initialization worked, this is quite simple:

Build the image:

```
$ docker build -t gobblerd:latest .
```

Make sure the `image` field of the `gobblerd` service in `docker-compose.yml` is `gobblerd:latest` (or whatever tag you built the image with) and then just bring the cluster up:
(If you're feeling confident/adventurous you can use the `-d` switch to daemonize the docker compose cluster as the logging is a bit enthusiastic right now)

```
$ docker compose up
```

### Useful stuff for debugging

The CockroachDB dashboard can be accessed at http://localhost:8080
The CockroachDB can be connected directly via the included client: `docker compose exec roach1 ./cockroach sql --insecure`
The Gobbler server is exposed on port 80 (http://localhost/upload, http://localhost/api/replays, http://localhost/api/replays/{id})

### Uploading replays

There's currently no UI so the most convenient way to upload replays is using [Postman](http://postman.com)

* Set the request method to POST and the URL to http://localhost/upload
* Go to the `Body` tab and select `form-data`
* Set the `Key` for the first row to `replay`
* There's a dropdown when hovering over a field in the `Key` column, set it to `File` and 
* Once set to `File` you can browse for the file you want to upload in the `Value` column

### Deploying UI code

All of the files are embedded in the binary which means you'll have to rebuild the Docker image when they're changed.
It's either this or fumbling around with Docker volumes and configurations which is just inconvenient.

The UI server is set up to cater to SPAs, no matter what URL you pass in, if it's not caught by the routes before it (api and asset paths) it will render the `ui/templates/index.html` file.

If you've made changes to the UI and built it, you can copy the asset files to `ui/assets` and update the `index.html` file in `ui/templates`.

For example creating a boilerplate project with Vite and Vue and building it resulted in the following `dist/` directory:
```
	dist/assets/index.46fee175.js
	dist/assets/index.75e12fc4.css
	dist/assets/vite.4a748afd.svg
	dist/assets/vue.5532db34.svg
```

To add these to GobblerD, you'll have to copy them to the `ui/assets` directory. The files in this directory are going to be served at `http://localhost/assets/<filename>`.

When that's done, update the `ui/templates/index.html` template as well with the built index (`dist/index.html`).

As I mentioned earlier, you'll have to rebuild the image, so let's do that:
```
$ docker build -t gobblerd:latest .
```

When you bring the cluster up and navigate to `http://localhost` you should see your changes.