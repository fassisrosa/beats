# Openstackbeat

Welcome to Openstackbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/fassisrosa/beats`

## Getting Started with Openstackbeat

### Init Project
To get running with Openstackbeat, run the following commands:

```
glide update --no-recursive
make update
```


To push Openstackbeat in the git repository, run the following commands:

```
git init
git add .
git commit
git remote set-url origin https://github.com/fassisrosa/beatsopenstackbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Openstackbeat run the command below. This will generate a binary
in the same directory with the name openstackbeat.

```
make
```


### Run

To run Openstackbeat with debugging output enabled, run:

```
./openstackbeat -c openstackbeat.yml -e -d "*"
```


### Test

To test Openstackbeat, run the following commands:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/openstackbeat.template.json and etc/openstackbeat.asciidoc

```
make update
```


### Cleanup

To clean  Openstackbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Openstackbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/fassisrosa/beats
cd ${GOPATH}/github.com/fassisrosa/beats
git clone https://github.com/fassisrosa/beats/openstackbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
