host="coap.me"
port="5683"
testPath="/test"
largePath="/large"
etagPath="/etag"
payload="Hello world from absmach/coap-cli"

#Get
echo "Sending GET request to $host:$port$testPath"
coap-cli get $testPath -H $host -p $port
sleep 1

#Get with blockwise transfer
echo "Sending GET request with blockwise transfer to $host:$port$testPath"
coap-cli get $largePath -H $host -p $port
sleep 1

#Post
echo "Sending POST request to $host:$port$testPath"
coap-cli post $testPath -H $host -p $port -d $payload 
sleep 1

#Post with content format
echo "Sending POST request with content format to $host:$port$testPath"
coap-cli post $testPath -H $host -p $port -d $payload -c 50
sleep 1

#Post with authentication
echo "Sending POST request with authentication to $host:$port$testPath"
coap-cli post $testPath -H $host -p $port -d $payload --auth "test"
sleep 1

#Put
echo "Sending PUT request to $host:$port$testPath"
coap-cli put $testPath -H $host -p $port -d $payload
sleep 1

#Delete
echo "Sending DELETE request to $host:$port$testPath"
coap-cli delete $testPath -H $host -p $port
