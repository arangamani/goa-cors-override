# Goa CORS Override

This is a simple example that demonstrates how to override Goa's built-in CORS handler.

There may be cases where it is necessary to have some methods that use HTTP OPTIONS verb and since CORS handler by default mounts on all methods, we need to route the request based on the presence of the Origin header.