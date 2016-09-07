#!/usr/bin/env python

#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements. See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership. The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License. You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.
#

import sys, glob
sys.path.append('./gen_py')
sys.path.insert(0, glob.glob('./git.apache.org/thrift.git/lib/py/build/lib.*')[0])

from gen_py.SystemMsgRpc.SystemMsgRpcSvr import *
from gen_py.SystemMsgRpc.SystemMsgRpcSvr import *

from thrift import Thrift
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol
import datetime

try:

  # Make socket
  transport = TSocket.TSocket('203.195.162.110', 8100)
  #transport = TSocket.TSocket('10.143.76.201', 8100)

  # Buffering is critical. Raw sockets are very slow
  transport = TTransport.TBufferedTransport(transport)

  # Wrap in a protocol
  protocol = TBinaryProtocol.TBinaryProtocol(transport)

  # Create a client to use the protocol encoder
  client = Client(protocol)

  # Connect!
  transport.open()
  

  try:
      
    Vercode = ''
    Appkey = "00b6413a92d4c1c84ad99e0a"
    CidList = ['Eddie']
    Title= 'titlet-test'
    MsgContent = 'MsgContent-test'
    #resp = client.SendSysMsgToCids(Appkey , Vercode , CidList , Title , MsgContent )
    #print "resp:", resp
    
    resp = client.SendSysMsgToAppkey(Appkey , Vercode  , Title , MsgContent )
    print "resp:", resp
    
  except Exception , e:
    print  "111111" , e 

  transport.close()

except Thrift.TException, tx:
  print '===== %s' % (tx.message)

