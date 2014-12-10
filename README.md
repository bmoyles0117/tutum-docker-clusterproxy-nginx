tutum-docker-clusterproxy
=========================

Nginx image that balances between linked containers and, if launched in Tutum, 
reconfigures itself when a linked cluster member joins or leaves


Usage - THIS SECTION WILL NOT WORK UNTIL THIS IS A PUBLIC IMAGE
-----

Make sure your application container exposes port 80. Then, launch it:

	docker run -d --name web1 tutum/hello-world
	docker run -d --name web2 tutum/hello-world

Then, run tutum/nginxproxy linking it to the target containers:

	docker run -d -p 80:80 --link web1:web1 --link web2:web2 tutum/nginxproxy


Configuration
-------------

You can overwrite the following Nginx configuration options:

* `PORT` (default: `80`): Port Nginx will bind to, and the port that will forward requests to.
* `POLLING_INTERVAL` (default: `30`): The default polling interval that the reloader checks for container updates

Check [the Nginx configuration manual](http://nginx.org/en/docs/) for more information on the above.


Usage within Tutum - THIS SECTION WILL NOT WORK UNTIL THIS IS A PUBLIC IMAGE
------------------

Launch the service you want to load-balance using Tutum.

Then, launch the load balancer. To do this, select `Jumpstarts` > `Proxies` and select `tutum/nginxproxy`. During the 3rd step of the wizard, link to the service created earlier, and add "Full Access" API role (this will allow Nginx to be updated dynamically by querying Tutum's API). 

That's it - the proxy container will start querying Tutum's API for an updated list of containers in the service and reconfigure itself automatically.

How to use this container
-------------------------
1. My service container(hello) exposes port 8080, I want the Nginx listens to port 80

    Run this container with `--link hello:hello -e PORT=8080 -p 80:80`

2. My service container(hello) exposes port 80, I want the Nginx listens to port 8080

    Run this container with `--link hello:hello -p 8080:80`

TODO
----
* SSL Support
* VHosts
* Multiple linked services - Currently the first one gains dominance