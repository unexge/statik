statik
=====
simple static file server


example config
```yml
listen: "0.0.0.0:6029"
https:
  cert: "./cert.pem"
  key: "./key.pem" 
files:
  "/index.html":
    cache: true
    gzip: true
    push: 
      - "/assets/css/app.css"
      - "/assets/js/app.js"

```
