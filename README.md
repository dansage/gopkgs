# GoPkgs
GoPkgs is a simple HTTP server that ensures any Go packages I created using a custom domain name resolve correctly while
also ensuring any end users that visit the package URL are directed to the appropriate repository.

## Getting Started
There are two files you'll want to update before starting the server itself:

### `internal/http.go`
Replace the default value of the `to` variable in the `handle` method with the fallback URL. This URL is where clients
are redirected to when they visit a page that does not belong to a known package.

### `internal/resources/packages.json`
This file contains the packages relevant for my custom domain. You should replace the contains of this file with your
packages:

```json
{
  "gopkgs": "https://github.com/dansage/gopkgs"
}
```

In this example the `go.dsage.org/gopkgs` package will resolve to `https://github.com/dansage/gopkgs`.

### Build and start the server
With the required modifications complete, build and start the server
```shell
go build -o gopkgs
./gopkgs
```

If you want to run the server as a systemd service, an example service definition has been included:

```shell
# copy the required files
sudo cp gopkgs.service /usr/lib/systemd/system/gopkgs.service
sudo cp gopkgs /usr/bin/gopkgs

# reload the systemd definitions
sudo systemctl daemon-reload

# enable and start the server
sudo systemctl enable --now gopkgs
```

GoPkgs is intended to be placed behind a reverse proxy and does not support TLS.

## Issues & Support
Issues are welcome but there is no guaranteed support available at this time.

## Security Vulnerabilities
If you discover a security vulnerability within GoPkgs, please email me at [security@mail.dsage.org][1] with
details. Please understand that this project receives no funding of any kind and there are no bug bounties available.

## License
GoPkgs is open-source software released under the [MIT License][2].

[1]: mailto:security@mail.dsage.org
[2]: https://choosealicense.com/licenses/mit/
