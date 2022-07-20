# xdry - Programming language agnostic clone detection

xdry detects duplicated behaviour in your application, even if the duplicated passages are not exact matches
to each other. Likely you should read a bit about the [theoretical background](http://xdry.leovie.de/blog/0.html)
for a better understanding.

## Run via Docker (recommended)

```bash
docker run -v {path_to_project}:/project leovie/xdry -h
```

## Run via binary

Download latest release and run via

```bash
./xdry {path_to_project} -h
```

## Documentation
see [here](http://xdry.leovie.de/documentation)

## Thanks

Special thank you belongs to [queo GmbH](https://www.queo.de) for sponsoring the development and maintenance of xdry.