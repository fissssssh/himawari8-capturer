# himawari8-capturer

himawari8-capturer is a tool that can get picture of the earth taken by the himawari 8

## Install

* [Download Release.](https://github.com/fissssssh/himawari8-capturer/releases)
* Get the source code and build.

### How to get source code ?

You have 2 ways to get source code

1. Clone the repository
   ```shell
   $ git clone https://github.com/fissssssh/himawari8-capturer.git
   ```
2. [Download the source code](https://github.com/fissssssh/himawari8-capturer/archive/refs/heads/main.zip)

### How to build ?

```shell
$ cd himawari8-capturer
$ go build -o build/himawari8-capturer cmd/cli/main.go
```

## Usage

```shell
$ ./himawari-capturer [-q quality] [-t unix_millisecond_timestamp] [-l shorelines_color] [-p proxy_url]
```
| Option | Access values                                                                                                          | Remark              |
| ------ | ---------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `-q`   | `1`: `550*550`<br>`2`: `1100*1100`<br>`3`: `2200*2200`<br>`4`: `4400*4400` <br>`5`: `8800*8800` <br>`6`: `11000*11000` | Resolution of image |
| `-t`   | Unix millisecond timestamp                                                                                             | Time of image       |
| `-l`   | `0`: `Ignore`<br>`1`: `Red`<br>`2`: `Green`<br>`3`: `Yellow`<br>                                                       | Color of shorelines |
| `-p`   | HTTP proxy url                                                                                                         | HTTP Proxy          |

### Example

```shell
$ ./build/himawari8-capturer -t 1634359675264
2021/10/16 23:01:19 Get 1-1 tile image from https://himawari8.nict.go.jp/img/D531106/2d/550/2021/10/16/044000_1_1.png...
2021/10/16 23:01:19 Get 0-1 tile image from https://himawari8.nict.go.jp/img/D531106/2d/550/2021/10/16/044000_0_1.png...
2021/10/16 23:01:19 Get 0-0 tile image from https://himawari8.nict.go.jp/img/D531106/2d/550/2021/10/16/044000_0_0.png...
2021/10/16 23:01:19 Get 1-0 tile image from https://himawari8.nict.go.jp/img/D531106/2d/550/2021/10/16/044000_1_0.png...
2021/10/16 23:01:23 Get 0-1 tile image done!
2021/10/16 23:01:23 Get 1-0 tile image done!
2021/10/16 23:01:23 Get 1-1 tile image done!
2021/10/16 23:01:23 Get 0-0 tile image done!
2021/10/16 23:01:23 Tile images were composed!
2021/10/16 23:01:23 Saving image to himawari8_20211016T044000Z.png
2021/10/16 23:01:24 All done!
```

![himawari8_20211016T044000Z.png](himawari8_20211016T044000Z.png)
