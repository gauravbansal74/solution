git checkout master && git pull

IMAGENAME=gauravbansal74/solution
SHORTREV=`git rev-parse --short HEAD`
BRANCH=`git branch | grep "*" | cut -d\  -f2`
TAG=`git describe --tags --abbrev=0`

if [ $TAG != "undefined" ]; then
    IMAGEREV=`echo $TAG | sed "s/\^0$//"`
    echo "Using tag $IMAGEREV"
    echo "Resetting repo..."
    git checkout $TAG
else
    IMAGEREV="$BRANCH-$SHORTREV"
    echo "Using tag $IMAGEREV"
fi

CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

if [ ! -d data ]; then
  mkdir data
fi
TAGINFO=`echo $TAG | sed "s/\^0$//"`
echo "Creating Info.json file to get latest version of commit"
echo "{ \"version\":"\"$TAGINFO\"", \"short\":"\"$SHORTREV\"", \"status\":\"OK\" }" > "data/info.json"

IMAGETAG=$IMAGENAME:$IMAGEREV
echo "Building Docker image with tag $IMAGETAG..."
docker build -t $IMAGETAG .
docker push $IMAGETAG
