# Gohook Server

Server for the gohook client application.

## What is it

Gohook server is an api which accepts long lived connections for gohook client application. It identifies each gohook client
that's connected in order to send them the correct webhooks as it receives them. It also exposes an REST HTTP api for creating,
updating, and removing the clients webhooks.


The server doesn't care about the command or script that the client will run when the webhook is hit. It's only job is to expose
the public webhooks and tunnel them down to the client.
