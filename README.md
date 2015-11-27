`extract` is a command line tool for extracting images and links from websites

install:

    go get github.com/zshipko/extract

usage:

    extract [-tag] [?search] urls...

examples:

    $ extract -img google.com
    > http://google.com/images/srpr/logo9w.png

    # find images on pages linked to from nytimes.com that include "nytimes" in the url
    $ extract -a nytimes.com | xargs extract ?nytimes -img
    > long list of images...

