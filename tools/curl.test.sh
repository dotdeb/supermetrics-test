#
#   Read settings from test.http and execute api call. This script must be runned from tools/
# 

SRC="example.http"

ADDR=$(sed -n -e '/^GET/p' $SRC)
AUTH=$(sed -n -e '/^Authorization/p' $SRC)

curl -v --request $ADDR --header "$AUTH"
