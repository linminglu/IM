
#-*- coding: utf8 -*-

import urllib
import urllib2
import json

def curl_keystone_failed():
    url = 'http://14.29.84.60:18086/cs/account/enable'
    values = '{"app_key":"9e5fd4272c457cef1e5c2605", "account":"kefuadmin", "enable":1}'
    #values = "aaaabbbcccdd"
    # params = urllib.urlencode(values)
    #params = json.dumps(values)
    params = values
    headers = {"Content-type":"application/json","Accept": "application/json", "Connection":"close", "Authorization":"Basic OWU1ZmQ0MjcyYzQ1N2NlZjFlNWMyNjA1OmYyMmU2MDI5ZjgxNWJhOGVkMzI3OGM3NTJkMjNjZmVm"}
    #headers = {"Content-type":"application/x-www-form-urlencoded; charset=UTF-8","Accept": "application/json", "Connection":"close", "Authorization":"Basic OWU1ZmQ0MjcyYzQ1N2NlZjFlNWMyNjA1OmYyMmU2MDI5ZjgxNWJhOGVkMzI3OGM3NTJkMjNjZmVm"}
    req = urllib2.Request(url, params, headers)
    response = urllib2.urlopen(req)
    cookie = response.info().getheader('Set-Cookie')
    print cookie
    print response.read()
     
if __name__ == "__main__":
    curl_keystone_failed()

