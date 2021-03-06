# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 4367
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server  certs=${false}
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Check running twice
    ${status}=  Get State Of Github Issue  5015
    Run Keyword If  '${status}' == 'closed'  Fail  Test Group0-Bugs/4367.robot needs to be updated now that Issue #5015 has been resolved
    #${name}=  Generate Random String  15
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -t --name ${name} busybox ls
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  proc
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start -a ${name}
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  proc
