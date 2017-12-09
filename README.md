# tcb
This is a tool for building toolchains based on gcc. The toolchains are built using bitbake which means we need a linux environment to run in. By running in a Docker container we are able to execute these builds independently of the user operating system, as long as there is a docker engine available. Finally, by using go as the client interface we can provide the same client tool for both windows and linux environments. 

## Configuration
We are able to build a range of compilers based on different versions of gcc, libc, kernel and target, host, build-triplets.
tcb uses `<name>.conf` files inside the `meta-crosstools/conf/toolchains/` directory to identify the kind of compilers it can build.

So, in order to build the native-mingw configuration, which is a cross-compiler from `x86_64-pc-linux-gnu` to `x86_64-w64-mingw32`, one would issue:
```bash
>tcb install native-mingw
``` 

By default tcb will clone the https://github.com/staffano/meta-crosstools repository to use as the backbone builder, but this can be overridden by the `builder.repo.url` flag to the tcb command. For all commands and flags avaialable to the tcb command, please use
```bash
>tcb --help
```

## License

See [LICENSE](LICENSE).


