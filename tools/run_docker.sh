#
#   Run dockerfile.
#

DockerfilePath=.
DotenvPath=.

# Find dockerfile
if [ -e $DockerfilePath/Dockerfile ]; then
    echo "Dockerfile exists in current path"
else
    echo "Check if parent has dockerfile"
    if [ -e ../Dockerfile ]; then
        echo "Parent has dockerfile. Use that"
        DockerfilePath=..
    else
        echo "Cannot find dockerfile. Exit!"
        exit 1
    fi
fi

# Find .env-file
if [ -e $DotenvPath/.env ]; then
    echo ".env exists in current path"
else
    echo "Check if parent has .env"
    if [ -e ../.env ]; then
        echo "Parent has .env"
        DotenvPath=..
    else
        echo "Cannot find .env. Exit!"
        exit 2
    fi
fi

#
#   Read data from the .env file. This should be more automated.
#
SECRET=$(sed -n -e '/^API_SECRET=/p' $DotenvPath/.env)
ISSUER=$(sed -n -e '/^JWT_ISSUER=/p' $DotenvPath/.env)
AUDIENCE=$(sed -n -e '/^JWT_AUDIENCE=/p' $DotenvPath/.env)
PORT=$(sed -n -e '/^PORT=/p' $DotenvPath/.env)
TIMEOUT_SEC=$(sed -n -e '/^TIMEOUT_SEC=/p' $DotenvPath/.env)

docker build --progress plain --tag homework $DockerfilePath
docker run -p 8000:8000 \
    -e $PORT \
    -e $ISSUER \
    -e $AUDIENCE \
    -e $SECRET \
    -e $TIMEOUT_SEC \
    homework