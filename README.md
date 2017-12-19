# s3hash-go
Golang Compute Hash for S3

## Installing

```
go get -u github.com/tsujimic/s3hash-go
```

## Using the s3hash-go
required ~/.aws/config and ~/.aws/credentials

```
$ s3hash-go.exe -h
NAME:
   compute hash - usage compute hash

USAGE:
   s3hash-go.exe [global options] command [command options] [arguments...]

VERSION:
   0.1.1

COMMANDS:
     md5         compute hash md5
     sha1        compute hash sha1
     sha224      compute hash sha224
     sha256      compute hash sha256
     sha384      compute hash sha384
     sha512      compute hash sha512
     sha512_224  compute hash sha512_224
     sha512_256  compute hash sha512_256
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        enable debugging
   --help, -h     show help
   --version, -v  print the version

$ s3hash-go.exe md5 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha1 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha224 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha256 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha384 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha512 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha512_224 --input "/bucket/object" --output "hash.json"
$ s3hash-go.exe sha512_256 --input "/bucket/object" --output "hash.json"
```
