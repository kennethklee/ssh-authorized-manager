SSH Authorized Manager
======================

Manage SSH keys for your servers. For small teams with many servers.

> TODO screenshot here


Features
---------

- [Pocketbase admin dashboard](https://pocketbase.io/)
- Header authentication (for authentication with other services, SAML, Webauthn, Tailscale, etc.)
- Manage users and SSH access on servers.
- Verify server host key to prevent MITM attacks.
- Verify server hostname to prevent port redirection attacks.


Docker Usage
-------------

Basic usage:

    $ docker run -dp 8090:8090 --name=ssh-authorized-manager <image>

With header authentication:

    $ docker run -dp 8090:8090 --name=ssh-authorized-manager -e "HEADER_AUTH_EMAIL=X-Auth-Email" -e "HEADER_AUTH_NAME=X-Auth-Name" -e "AUTO_CREATE_USERS=1" <image>

With listen address on a different port:

    $ docker run -dp 8090:8000 --name=ssh-authorized-manager <image> ssham serve --http 0:8000

Once running, visit http://localhost:8090/_ to create your first admin.

Options
--------

Environment Variables

- `HEADER_AUTH_EMAIL`: (default blank) Header name to use for user's email address.
- `HEADER_AUTH_NAME`: (default blank) Header name to use for user's name.
- `AUTO_CREATE_USERS`: (default blank) If not empty, automatically create users in the system.


Alternatives
-------------

A larger team would use a more robust system. This is a simple system for a small team.

If you have a large team, try these:

- https://www.boundaryproject.io/
- https://goteleport.com/
- ca-certs like in https://engineering.fb.com/2016/09/12/security/scalable-and-secure-access-with-ssh/
- ansible - https://gitlab.com/consensus.enterprises/ansible-roles/ansible-role-admin-users


Contributing
-------------

Feature requests and pull requests are welcome.

### Requirements

- docker
- docker-compose
- [go-task](https://taskfile.dev/installation/)

### Running

There are two modes -- development `ENV=dev` and test `ENV=test`.

In `development`, two containers start up, app and web. The app container does live-reloading. The web container does hot-reloading. Any code changes anywhere in the project will be automatically reloaded.

Start the stack:

    $ task up

Stop the stack:

    $ task stop

Destroy the stack:

    $ task down

Start the stack in a production-like mode (single container):

    $ ENV=test task up

In `test` mode, the app runs in a production-like fashion, hence the term staging. Live reloading and hot reloading is disabled.
