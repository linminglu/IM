
#-*- coding: utf8 -*-

import urllib
import urllib2
import json

def curl_keystone_failed():
    url = 'http://122.13.81.202:18086/cs/account/register'
    values = '{"account":"yaosha001", "password":"AEAA86A4CAA581DDDB4AA0EA3CAEAFBF", "ext_info":{ "nick_name":"要啥APP客服001", "image_id":"", "email":"", "tel":""}}'
    # 这里千万不要仿照网上的方法进行加密，因为它本身就没有加密的一个过程！不然还是会返回400的！
    # params = urllib.urlencode(values)
    #params = json.dumps(values)
    params = values
    headers = {"Content-type":"application/json","Accept": "application/json", "Authorization":"Basic YzdmYzUwODNjYzk2MzJkNDk0NTYzODliOjExMTEx"}
    req = urllib2.Request(url, params, headers)
    response = urllib2.urlopen(req)
    print response.read()
     
if __name__ == "__main__":
    curl_keystone_failed()

