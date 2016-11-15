# k8stail

`tail -f` experience for Kubernetes Pods

As you know, `kubectl logs` can stream only ONE pod at the same time. `k8stail` enables you to watch __log streams of ALL pods__ in the specified namespace or labels in real time, like `tail -f`.

## Installation

### Using Homebrew (OS X only)

Preparing... :construction_worker:

### Precompiled binary

Precompiled binaries for Windows, OS X, Linux are available at [Releases](https://github.com/dtan4/k8stail/releases).

### From source

```bash
$ go get -d github.com/dtan4/k8stail
$ cd $GOPATH/src/github.com/dtan4/k8stail
$ make deps
$ make install
```

## Usage

Logs of all Pods in the specified namespace are streaming. When new Pod is added, logs of the Pod also appears.

```bash
$ bin/k8stail -namespace awesome-app
Namespace: awesome-app
Labels:
----------
Pod awesome-app-web-4212725599-67vd4 has detected
Pod awesome-app-web-4212725599-6pduy has detected
Pod awesome-app-web-4212725599-lbuny has detected
Pod awesome-app-web-4212725599-mh3g1 has detected
Pod awesome-app-web-4212725599-pvjsm has detected
[awesome-app-web-4212725599-mh3g1]  | creating base compositions...
[awesome-app-web-4212725599-zei9h]  |    (47.1ms)  CREATE TABLE "schema_migrations" ("version" character varying NOT NULL)
[awesome-app-web-4212725599-zei9h]  |    (45.1ms)  CREATE UNIQUE INDEX  "unique_schema_migrations" ON "schema_migrations"  ("version")
[awesome-app-web-4212725599-zei9h]  |   ActiveRecord::SchemaMigration Load (1.8ms)  SELECT "schema_migrations".* FROM "schema_migrations"
[awesome-app-web-4212725599-zei9h]  | Migrating to CreatePosts (20160218082522)
```

With `-timestamps` option, log timestamp is printed together.


```bash
$ bin/k8stail -namespace awesome-app -timestamps
Namespace: awesome-app
Labels:
----------
Pod awesome-app-web-4212725599-67vd4 has detected
Pod awesome-app-web-4212725599-6pduy has detected
Pod awesome-app-web-4212725599-lbuny has detected
Pod awesome-app-web-4212725599-mh3g1 has detected
Pod awesome-app-web-4212725599-pvjsm has detected
[awesome-app-web-4212725599-mh3g1] 2016-11-15T10:57:22.178667425Z  | creating base compositions...
[awesome-app-web-4212725599-zei9h] 2016-11-15T10:57:22.309011520Z  |    (47.1ms)  CREATE TABLE "schema_migrations" ("version" character varying NOT NULL)
[awesome-app-web-4212725599-zei9h] 2016-11-15T10:57:22.309053601Z  |    (45.1ms)  CREATE UNIQUE INDEX  "unique_schema_migrations" ON "schema_migrations"  ("version")
[awesome-app-web-4212725599-zei9h] 2016-11-15T10:57:22.463700110Z  |   ActiveRecord::SchemaMigration Load (1.8ms)  SELECT "schema_migrations".* FROM "schema_migrations"
[awesome-app-web-4212725599-zei9h] 2016-11-15T10:57:22.463743373Z  | Migrating to CreatePosts (20160218082522)
```

With `-labels` option, you can filter Pods to watch.

```bash
$ bin/k8stail -namespace awesome-app -labels name=awesome-app-web
Namespace: awesome-app
Labels:    name=awesome-app-web
----------
Pod awesome-app-web-4212725599-67vd4 has detected
Pod awesome-app-web-4212725599-6pduy has detected
Pod awesome-app-web-4212725599-lbuny has detected
Pod awesome-app-web-4212725599-mh3g1 has detected
Pod awesome-app-web-4212725599-pvjsm has detected
[awesome-app-web-4212725599-mh3g1]  | creating base compositions...
[awesome-app-web-4212725599-zei9h]  |    (47.1ms)  CREATE TABLE "schema_migrations" ("version" character varying NOT NULL)
[awesome-app-web-4212725599-zei9h]  |    (45.1ms)  CREATE UNIQUE INDEX  "unique_schema_migrations" ON "schema_migrations"  ("version")
[awesome-app-web-4212725599-zei9h]  |   ActiveRecord::SchemaMigration Load (1.8ms)  SELECT "schema_migrations".* FROM "schema_migrations"
[awesome-app-web-4212725599-zei9h]  | Migrating to CreatePosts (20160218082522)
```

### Options

|Option|Description|Required|Default|
|---------|-----------|-------|-------|
|`-kubeconfig=KUBECONFIG`|Path of kubeconfig||`~/.kube/config`|
|`-labels=LABELS`|Label filter query (e.g. `app=APP,role=ROLE`)|||
|`-namespace=NAMESPACE`|Kubernetes namespace||`default`|
|`-timestamps`|Include timestamps on each line||`false`|

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
