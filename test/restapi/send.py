
#-*- coding: utf8 -*-

import urllib
import urllib2
import json

def curl_keystone_failed():
    url = 'http://14.29.84.60:18086/cs/msg/send'
    values = '{"to_user":"zzz", "msg_type":1, "content":"aaaaaaaaa"}'
    #values = "aaaabbbcccdd"
    # params = urllib.urlencode(values)
    #params = json.dumps(values)
    params = values
    headers = {"Content-type":"application/json","Accept": "application/json", "Connection":"close", "Authorization":"Basic OWU1ZmQ0MjcyYzQ1N2NlZjFlNWMyNjA1OmYyMmU2MDI5ZjgxNWJhOGVkMzI3OGM3NTJkMjNjZmVm"}
    req = urllib2.Request(url, params, headers)
    response = urllib2.urlopen(req)
    print response.read()
     
if __name__ == "__main__":
    curl_keystone_failed()

