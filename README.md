# Bullet

![](assets/bullet_128.png)

Bullet is a fast and flexible application deployment tool built by Furqan Software.

At [Furqan Software](https://furqansoftware.com/), Bullet helps us setup and deploy prototype applications with minimal effort.

## Getting Started

### Install from Source

``` sh
go get github.com/FurqanSoftware/bullet
```

### Copy an Example App

``` sh
cp -r $GOPATH/src/github.com/FurqanSoftware/examples/hello .
```

### Set up a Server

``` sh
bullet -H {host} setup
```

### Deploy App to Server

``` sh
make release
bullet -H {host} deploy hello.tar.gz
```

### Scale Programs on Server

``` sh
bullet -H {host} scale web=1
```

## Acknowledgements

- [Nikita Golubev](http://www.flaticon.com/authors/nikita-golubev) - For the bullet icon
