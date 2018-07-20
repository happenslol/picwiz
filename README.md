# picvoter.xyz application

## Setup
```bash
$ go get -u -v github.com/gobuffalo/buffalo/buffalo

$ go get -u -v github.com/happenslol/picwiz
$ cd $GOPATH/src/github.com/happenslol/picwiz

# database needs to have username/pw postgres/postgres
$ buffalo db create -a
$ buffalo db m up

# run dev server
$ buffalo dev
```

## Configuration
Create a `.env` file in the project root and set the following vars:
```
STORAGE_LOCATION=absolute path to where the images should be stored
SCAN_LOCATIONS=comma seperated paths where new folders with images can be found
```

The storage location needs to contain a `static` and a `imports` directory

## Dependencies
Only libvips >= 8.6.3 needs to be installed, it's the same library that the nodejs sharpen uses for resizing images
